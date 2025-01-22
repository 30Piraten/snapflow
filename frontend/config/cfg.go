package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func S3Client() (*s3.Client, error) {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	// Get region from .env
	s3Region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Region),
	)

	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
