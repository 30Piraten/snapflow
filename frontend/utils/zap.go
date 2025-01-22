package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() error {
	var err error
	Logger, err = zap.NewProduction() //zap.NewDevelopment()   for prod

	if err != nil {
		Logger.Fatal("Failed to initialise logger: %v", zap.Error(err))
	}

	return nil
}
