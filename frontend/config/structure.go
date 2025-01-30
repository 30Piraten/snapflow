package config

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// PrintJob represents a print request
type PrintJob struct {
	CustomerEmail       string `json:"customer_email"`
	PhotoID             string `json:"photo_id"`
	ProcessedS3Location string `json:"processed_s3_location"`
}

// PrintJob represents a print request message
var printJob *PrintJob

// AWS Clients
var (
	dynamoClient *dynamodb.Client
	snsClient    *sns.Client
)
