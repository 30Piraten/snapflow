package config

import (
	"crypto/rsa"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	s3Client     *s3.Client
)

// CloudFront private key path
var privateKeyPath string

var (
	cloudFrontClient *cloudfront.Client
	privateKey       *rsa.PrivateKey
	keyPairID        string
	cloudFrontDomain string
	policy           string
)

// SignedURLInfo struct
type SignedURLInfo struct {
	CloudFrontDomain string `json:"cloudfront_domain"`
	ObjectKey        string `json:"object_key"`
	Expires          int64  `json:"expires"`
	KeyPairID        string `json:"key_pair_id"`
	Policy           string `json:"policy"`
	CustomerName     string `json:"customer_name"`
	OrderID          string `json:"order_id"`
}
