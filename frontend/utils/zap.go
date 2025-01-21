package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() error {
	var err error
	Logger, err = zap.NewDevelopment() // use zap.NewProduction() for prod

	if err != nil {
		return err
	}

	return nil
}
