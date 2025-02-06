package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger initialises the logger for the
// application. It returns an error if it fails.
func InitLogger() error {
	var err error
	Logger, err = zap.NewProduction()

	if err != nil {
		Logger.Fatal("Failed to initialise logger: %v", zap.Error(err))
	}

	return nil
}
