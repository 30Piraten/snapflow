package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"sync"

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
		return handleError(c, fiber.StatusInternalServerError, "Failed to generate presigned URL")
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

	// Close aresults channel after all processing is complete
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
func processFile(file *multipart.FileHeader) FileProcessingResult {
	result := FileProcessingResult{
		Filename: file.Filename,
	}

	// Validate MIME type
	if !isValidMimeType(file.Header.Get("Content-Type")) {
		result.Error = fmt.Errorf("invalid file type: only JPEG and PNG are allowed")
		return result
	}

	// Process the image
	processedImage, err := processImage(file)
	if err != nil {
		result.Error = err
		return result
	}

	// Save the processed image
	result.Path = fmt.Sprintf("./uploads/%s", file.Filename)
	if err := saveProcessedImage(processedImage, result.Path); err != nil {
		result.Error = err
		return result
	}

	return result
}

// Helper functions

func isValidMimeType(mimeType string) bool {
	return mimeType == "image/jpeg" || mimeType == "image/png"
}

func processImage(file *multipart.FileHeader) (image.Image, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// format -> _
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	// Resize if needed
	if file.Size >= ResizeThreshold {
		img = imaging.Resize(img, int(float64(
			img.Bounds().Dx())*CompressionFactor), 0, imaging.Lanczos)
	}

	return img, nil
}

func saveProcessedImage(img image.Image, path string) error {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
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

// // Handler configures the routes for our photo upload service
// func Handler(app *fiber.App) {

// 	// Render the upload form
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.Render("index", fiber.Map{
// 			"Title": "Photo Upload Service",
// 		})
// 	})

// 	// Handle form submission
// 	app.Post("/submit-order", func(c *fiber.Ctx) error {

// 		// Log the start of the request
// 		utils.Logger.Info("Processing order submission")

// 		// Parse the multipart form
// 		order := new(utils.PhotoOrder)
// 		if err := c.BodyParser(order); err != nil {
// 			utils.Logger.Error("Failed to parse form", zap.Error(err))

// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"error": "Failed to parse form",
// 			})
// 		}

// 		utils.Logger.Info("Form parsed successfully", zap.String("fullName", order.FullName))

// 		// Use the shared presigned URL generation function
// 		presignedResponse, err := utils.GeneratePresignedURL(order)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}

// 		form, err := c.MultipartForm()
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"error": "Failed to process uploaded files",
// 			})
// 		}

// 		// Get the files from form
// 		files := form.File["photos"]

// 		// Channel for errors
// 		errorChan := make(chan error, len(files))

// 		// WaitGroup for synchronisation
// 		var wg sync.WaitGroup

// 		// Process files concurrently
// 		for _, file := range files {
// 			wg.Add(1)
// 			go func(file *multipart.FileHeader) {
// 				defer wg.Done()

// 				// Add logic for resize photo uploads by 50%
// 				// Enforce size limit
// 				if file.Size > MaxFileSize {
// 					utils.Logger.Warn("File exceeds the maximum size", zap.String("filename", file.Filename))
// 					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 						"error": fmt.Sprintf("File %s exceeds the maximum size of 100MB", file.Filename),
// 					})
// 				}

// 				// Open file for processing
// 				src, err := file.Open()
// 				if err != nil {
// 					utils.Logger.Error("Failed to open file", zap.String("filename", file.Filename), zap.Error(err))
// 					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 						"error": fmt.Sprintf("Failed to open file %s", file.Filename),
// 					})
// 				}
// 				defer src.Close()

// 				// Decode image
// 				img, format, err := image.Decode(src)
// 				if err != nil {
// 					utils.Logger.Error("Failed to decode image", zap.String("filename", file.Filename))
// 					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 						"error": fmt.Sprintf("File %s is not a valid image", file.Filename),
// 					})
// 				}

// 				// Reduce size for large files
// 				if file.Szie >= ResizeThreshold {
// 					utils.Logger.Info("Resizing image", zap.String("filename", file.Filename))
// 					img = imaging.Resize(img, int(float64(img.Bounds().Dx())*CompressionFactor), 0, imaging.Lanczos)
// 				}

// 				// Save compressed image to a buffer
// 				var buf bytes.Buffer
// 				if format == "jpeg" {
// 					err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
// 				} else if format == "png" {
// 					err = png.Encode(&buf, img)
// 				} else {
// 					utils.Logger.Warn("Unsupported format for compression", zap.String("format", format))
// 					return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
// 						"error": fmt.Sprintf("File %s has unsupported format %s", file.Filename, format),
// 					})
// 				}

// 				if err != nil {
// 					utils.Logger.Error("Failed to encode image", zap.String("filename", file.Filename), zap.Error(err))
// 					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 						"error": fmt.Sprintf("Failed to process image %s", file.Filename),
// 					})
// 				}

// 				// Validate MIME type
// 				if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
// 					err := fmt.Errorf("Invalid file type: %s. Only JPEG and PNG are allowed.", file.Header.Get("Content-Type"))
// 					utils.Logger.Warn("Invalid file type", zap.String("filename", file.Filename))
// 					errorChan <- err
// 					return
// 				}

// 				// Save file locally or to S3
// 				filePath := fmt.Sprintf("./uploads/%s", file.Filename)
// 				if err := c.SaveFile(file, filePath); err != nil {
// 					utils.Logger.Error("Failed to save file", zap.String("filename", file.Filename), zap.Error(err))
// 					errorChan <- err
// 					return
// 				}

// 				utils.Logger.Info("Photos processed successfully", zap.String("filename", file.Filename))

// 			}(file)

// 		}

// 		// Wait for all goroutines to finish
// 		wg.Wait()
// 		close(errorChan)

// 		// Check for errors
// 		if len(errorChan) > 0 {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "One or more files failed to upload",
// 			})
// 		}

// 		return c.JSON(fiber.Map{
// 			"message":       "Order received successfully",
// 			"order":         order,
// 			"presigned_url": presignedResponse.URL,
// 			"order_id":      presignedResponse.OrderID,
// 		})
// 	})

// 	// Register the presigned URL route
// 	routes.Upload(app)
// }
