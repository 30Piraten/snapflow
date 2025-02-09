package models

import (
	"mime/multipart"
	"sync"
	"time"

	"go.uber.org/zap"
)

// PrintJob represents a print request
type PrintJob struct {
	CustomerEmail       string `json:"customer_email"`
	PhotoID             string `json:"photo_id"`
	ProcessedS3Location string `json:"processed_s3_location"`
}

// SignedURLInfo struct
type SignedURLInfo struct {
	CloudFrontDomain string `json:"cloudfront_domain"`
	ObjectKey        string `json:"object_key"`
	Expires          int64  `json:"expires"`
	KeyPairID        string `json:"key_pair_id"`
	Policy           string `json:"policy"`
	CustomerName     string `json:"customer_name"`
	OrderID          string `json:"order_id"`
}

type PhotoOrder struct {
	FullName  string                  `json:"fullName"`
	Location  string                  `json:"location"`
	Size      string                  `json:"size"`
	PaperType string                  `json:"paperType"`
	Email     string                  `json:"email"`
	Photos    []*multipart.FileHeader `json:"photos"`
}

// Quality settings for image processing
const (
	HighQuality   int = 85
	MediumQuality int = 75
	LowQuality    int = 65

	QualityStep int = 5  // Step to reduce JPEG quality
	MinQuality  int = 10 // Minimum JPEG quality

	ProcessedImageDir = "uploads"

	MaxConcurrentProcessing int   = 3               // Max concurrent uploads
	MaxFileCount            int   = 10              // Maximum 10 files per request
	TargetFileSize          int64 = 1 * 1024 * 1024 // 1MB total for all sizes
	MaxFileSize             int64 = 50 * 100 * 1024 // 50MB per file
	MaxTotalUploadSize      int   = 5 * 1024 * 1024 // 5MB total for all files
)

// ProcessingOptions defines configuration for image processing
type ProcessingOptions struct {
	MaxWidth         int
	MaxHeight        int
	Quality          int
	Sharpen          bool
	Format           string
	PreserveMetadata bool
	OptimiseSizeOnly bool
	TargetSizeBytes  int64 // -> New field for target file size
	MaxDimensions    Dimensions
}

type Dimensions struct {
	Width  int
	Height int
}

// ImageProcessor handles all image processing operations
type ImageProcessor struct {
	Logger *zap.Logger
	cache  *sync.Map // -> Allow cache for processed images
}

type ResponseData struct {
	Message      string      `json:"message"`
	Order        *PhotoOrder `json:"order"`
	PresignedURL []string    `json:"presigned_url"`
	OrderID      string      `json:"order_id"`
}

// FileProcessingResult holds the result of processing a single file
type FileProcessingResult struct {
	Filename      string           `json:"filename"`
	OriginalSize  int64            `json:"original_size"`
	ProcessedSize int64            `json:"processed_size"`
	Path          string           `json:"path"`
	Size          int64            `json:"size"`
	Error         *ProcessingError `json:"error,omitempty"`
	Duration      time.Duration    `json:"duration"`
	Quality       int              `json:"quality"`
}

// ProcessingError represents a structured processing error
type ProcessingError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	// Error   ProcessedError         `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type ProcessedError struct {
	Message string `json:"message"`
}

func (e *ProcessedError) Error() string {
	return e.Message
}

const (
	ErrCodeFileSave         = "CANNOT_SAVE_FILE"
	ErrCodeFileRead         = "CANNOT_READ_FILE"
	ErrCodeFileOpen         = "CANNOT_OPEN_FILE"
	ErrCodeInvalidRequest   = "INVALID_REQUEST"
	ErrCodeFileTooLarge     = "FILE_TOO_LARGE"
	ErrCodeTooManyFiles     = "TOO_MANY_FILES"
	ErrCodeProcessingFailed = "PROCESSING_FAILED"
	ErrCodeInvalidFormat    = "INVALID_FORMAT"
	ErrCodeStorageFailed    = "STORAGE_FAILED"
	ErrCodeFailedFileUpload = "NO_FILES_UPLOADED"
)
