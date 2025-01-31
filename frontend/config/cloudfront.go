package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/joho/godotenv"
)

var info *SignedURLInfo

// CloudFrontClient initializes a CloudFront client using the default AWS config.
func InitCloudFront() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	// Get the CloudFront credentials from the environment variables
	keyPairID = os.Getenv("CLOUDFRONT_KEY_PAIR_ID")
	privateKeyPath = os.Getenv("CLOUDFRONT_PRIVATE_KEY_PATH")
	cloudFrontDomain = os.Getenv("CLOUDFRONT_DOMAIN")

	// Read and parse the private key
	privateKey, err = readPrivateKey(privateKeyPath)
	if err != nil {
		log.Printf("failed to read private key: %v", err)
	}

	// Initialize AWS SDK configuration
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS configuration, %v", err)
	}

	cloudFrontClient = cloudfront.NewFromConfig(config)
}

// GenerateSignedURL generates a signed URL for accessing a photo via CloudFront.
func GenerateSignedURL(photoPath string) (string, error) {

	expires := time.Now().Add(1 * time.Minute)

	// Ensure the full CloudFront URL is used
	theFullURL := "https://" + cloudFrontDomain + "/" + photoPath

	// Initialize the CloudFront signer
	signer := sign.NewURLSigner(keyPairID, privateKey)

	// Generate the signed URL
	signedURL, err := signer.Sign(theFullURL, expires)
	if err != nil {
		log.Printf("failed to sign URL: %v", err)
		return "", err
	}

	return signedURL, nil
}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	// Read the private key from the file path
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Decode the private key
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Parse the RSA private key
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to assert type to *rsa.PrivateKey")
	}

	return privateKey, nil
}
