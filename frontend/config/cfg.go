package config

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

// S3Client initializes an S3 client using the default AWS config.
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

// S3 bucket config
type S3Bucket struct {
	S3Client *s3.Client
}

// UploadToS3 uploads a file to an S3 bucket.
// The function takes an S3 client, the name of the bucket to upload to, the key
// to use for the uploaded file, the file data, and the region the bucket is in.
// It returns an error if the upload fails.
// The function logs a successful upload to the console.
func UploadToS3(s3Client *s3.Client, bucketName, key string, fileData []byte, region string) error { // changes here

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileData),
	}

	// Here we upload the file
	_, err := s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error uploading to S3: %w", err)
	}

	log.Printf("File successfully uploaded to S3: s3://%s/%s\n", bucketName, key)

	return nil
}
