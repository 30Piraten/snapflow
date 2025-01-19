package utils

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func GeneratePresignedURL(client *s3.PresignClient, bucket, key string, metadata map[string]string, expiry time.Duration) (string, error) {

	presignedPut, err := client.PresignPutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			Metadata: metadata,
		},
		s3.WithPresignExpires(expiry),
	)

	if err != nil {
		return "", err
	}

	return presignedPut.URL, nil
}
