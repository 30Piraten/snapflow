package services

import (
	"sync"

	"go.uber.org/zap"
)

// ImageProcessor handles all image processing operations
type ImageProcessor struct {
	Logger *zap.Logger
	cache  *sync.Map // -> Allow cache for processed images
}

// NewImageProcessor creates a new instance of ImageProcessor with the provided Logger.
// The Logger is used for logging messages, and a new cache is initialized for caching processed images.
func NewImageProcessor(Logger *zap.Logger) *ImageProcessor {
	return &ImageProcessor{
		Logger: Logger,
	}
}
