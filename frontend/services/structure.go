package services

import (
	"sync"
	"time"

	"github.com/30Piraten/snapflow/utils"
	"go.uber.org/zap"
)

// image.go file definition
// Quality settings for image processing
const (
	HighQuality   = 85
	MediumQuality = 75
	LowQuality    = 65

	QualityStep = 5  // Step to reduce JPEG quality
	MinQuality  = 10 // Minimum JPEG quality

	ProcessedImageDir = "uploads"

	MaxConcurrentProcessing = 3                // Max concurrent uploads
	MaxFileCount            = 10               // Maximum 10 files per request
	TargetFileSize          = 1 * 1024 * 1024  // 1MB total for all sizes
	MaxFileSize             = 100 * 100 * 1024 // 100MB per file
	MaxTotalUploadSize      = 5 * 1024 * 1024  // 5MB total for all files
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
}

// ImageProcessor handles all image processing operations
type ImageProcessor struct {
	logger *zap.Logger
	cache  *sync.Map // -> Allow cache for processed images
}

// file.go Struct definition

type ResponseData struct {
	Message      string            `json:"message"`
	Order        *utils.PhotoOrder `json:"order"`
	PresignedURL string            `json:"presigned_url"`
	OrderID      string            `json:"order_id"`
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
	Type    string                 `json:"type"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
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
