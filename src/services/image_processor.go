package services

import (
	"sync"

	"go.uber.org/zap"
)

// ImageProcessor handles all image processing operations
type ImageProcessor struct {
	Logger *zap.Logger
	// Allow cache for processed images
	cache *sync.Map
}

// NewImageProcessor creates a new instance of
// ImageProcessor with the provided Logger.
func NewImageProcessor(Logger *zap.Logger) *ImageProcessor {
	return &ImageProcessor{
		Logger: Logger,
	}
}
