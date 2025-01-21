package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"strings"
	"sync"

	"github.com/30Piraten/snapflow/routes"
	"github.com/30Piraten/snapflow/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Constants for file size
const (
	MaxFileSize             = 100 * 1024 * 1024 // 100MB
	ResizeThreshold         = 10 * 1024 * 1024  // 10MB
	CompressionFactor       = 0.5               // 50% size reduction
	MaxConcurrentProcessing = 5
)

// ResponseData defines the structure of the API response
type ResponseData struct {
	Message      string            `json:"message"`
	Order        *utils.PhotoOrder `json:"order"`
	PresignedURL string            `json:"presigned_url"`
	OrderID      string            `json:"order_id"`
}

// FileProcessingResult holds the result of processing a single file
type FileProcessingResult struct {
	Filename string
	Error    error
	Path     string
}

// Handler configured the routes for our photo upload service
func Handler(app *fiber.App) {
	// Save the upload form
	app.Get("/", ServeUploadForm)

	// Handle form submissions
	app.Post("/submit-order", HandleOrderSubmission)

	// Register the presigned URL route
	routes.Upload(app)
}

// ServeUploadForm handles rendering the upload form
func ServeUploadForm(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Photo Upload Service",
	})
}

// HandleOrderSubmission processes the order form submission
func HandleOrderSubmission(c *fiber.Ctx) error {
	utils.Logger.Info("Starting order submission processing")

	// Parse the order details
	order, err := parseOrderDetails(c)
	if err != nil {
		return handleError(c, fiber.StatusBadRequest, "Failed to parse order details", err)
	}

	// Generate presigned URL
	presignedResponse, err := utils.GeneratePresignedURL(order)
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Failed to generate presigned URL", err)
	}

	// Process uploaded photos
	err = processUploadedFiles(c, order)
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Failed to process files", err)
	}

	// Return a successful response
	return c.JSON(ResponseData{
		Message:      "Order receibed successfully",
		Order:        order,
		PresignedURL: presignedResponse.URL,
		OrderID:      presignedResponse.OrderID,
	})
}

// parseOrderDetails extracts and validates order information from the request
func parseOrderDetails(c *fiber.Ctx) (*utils.PhotoOrder, error) {
	order := new(utils.PhotoOrder)
	if err := c.BodyParser(order); err != nil {
		utils.Logger.Error("Form parsing failed", zap.Error(err))
		return nil, err
	}

	// Validate required fields
	if err := validateOrder(order); err != nil {
		return nil, err
	}

	utils.Logger.Info("Order details parsed successfully",
		zap.String("fullName", order.FullName),
		zap.String("email", order.Email))

	return order, nil
}

// processUploadedFiles handles the concurrent processing of uploaded files
func processUploadedFiles(c *fiber.Ctx, order *utils.PhotoOrder) error {
	form, err := c.MultipartForm()
	if err != nil {
		return fmt.Errorf("failed to process multipart form: %w", err)
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return fmt.Errorf("no files uploaded")
	}

	// Create buffered channels for results and limiting concurrency
	results := make(chan FileProcessingResult, len(files))
	semaphore := make(chan struct{}, MaxConcurrentProcessing)

	// Start a worker pool for file processing
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			// Get the semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := processFile(file)
			results <- result
		}(file)
	}

	// Close results channel after all processing is complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and process results
	var errors []string
	for result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", result.Filename, result.Error))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("file processing errors: %v", errors)
	}

	return nil
}

// processFile handles the processing of a single file
// using the image processing logic
func processFile(file *multipart.FileHeader) FileProcessingResult {
	result := FileProcessingResult{
		Filename: file.Filename,
	}

	// Validate file size
	if file.Size > MaxFileSize {
		result.Error = fmt.Errorf("file exceeds maximum szie of 100MB")
	}

	// Validate MIME type
	if !isValidMimeType(file.Header.Get("Content-Type")) {
		result.Error = fmt.Errorf("invalid file type: only JPEG and PNG are allowed")
		return result
	}

	// Process the image
	processedImage, format, err := processImage(file)
	if err != nil {
		result.Error = err
		return result
	}

	// Save the processed image
	result.Path = fmt.Sprintf("./uploads/%s", file.Filename)
	if err := saveProcessedImage(processedImage, format, result.Path); err != nil {
		result.Error = err
		return result
	}

	return result
}

// ValidateOrder checks if all required fields are present and valid
func validateOrder(order *utils.PhotoOrder) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	var missingFields []string

	if strings.TrimSpace(order.FullName) == "" {
		missingFields = append(missingFields, "full name")
	}
	if strings.TrimSpace(order.Email) == "" {
		missingFields = append(missingFields, "email")
	}
	if strings.TrimSpace(order.Location) == "" {
		missingFields = append(missingFields, "location")
	}
	if strings.TrimSpace(order.Size) == "" {
		missingFields = append(missingFields, "size")
	}
	if strings.TrimSpace(order.PaperType) == "" {
		missingFields = append(missingFields, "paper type")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	// Basic email validation
	if !strings.Contains(order.Email, "@") || !strings.Contains(order.Email, ".") {
		return errors.New("invalid email format")
	}

	return nil
}

// Helper functions

func isValidMimeType(mimeType string) bool {
	return mimeType == "image/jpeg" || mimeType == "image/png"
}

// processImage handles the image processing with proper format detection
func processImage(file *multipart.FileHeader) (image.Image, string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open the file: %w", err)
	}
	defer src.Close()

	// Decode image and detect format
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, "", fmt.Errorf("failed to upload image: %w", err)
	}

	// Reset the file pointer so it can be resued
	if _, err := src.Seek(0, 0); err != nil {
		return nil, "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Resize the photo / image
	if file.Size >= ResizeThreshold {
		// Calcultae new dimensions while maintaining aspect ratio
		bounds := img.Bounds()
		originalWidth := bounds.Dx()
		originalHeight := bounds.Dy()

		// We get the newWidth & newHeight here
		newWidth := int(float64(originalWidth) * CompressionFactor)
		newHeight := int(float64(originalHeight) * CompressionFactor)

		// Then we perform the resizing of each image
		resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		// Simple bilinear scaling (changed to sophisticated)
		for y := 0; y < newHeight; y++ {
			for x := 0; x < newWidth; x++ {
				sx := float64(x) * float64(originalWidth) / float64(newWidth)
				sy := float64(y) * float64(originalHeight) / float64(newHeight)
				resized.Set(x, y, img.At(int(sx), int(sy)))
			}
		}
		img = resized

	}

	return img, format, nil
}

func saveProcessedImage(img image.Image, format string, path string) error {
	var buf bytes.Buffer
	var err error

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(&buf, img)
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	// Write the buffer to file
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

func handleError(c *fiber.Ctx, status int, message string, err error) error {
	utils.Logger.Error(message,
		zap.Error(err),
		zap.Int("status", status))

	return c.Status(status).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}
