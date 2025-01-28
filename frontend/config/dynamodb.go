package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
)

var (
	dynamoClient *dynamodb.Client
)

func Init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("failed to load configuration, ", err)
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

// InsertMetadata inserts metadata for a new photo upload.
func InsertMetadata(customerEmail, photoID string, timestamp int64) error {
	ok := godotenv.Load()
	if ok != nil {
		log.Fatal("failed to load .env file")
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	_, err := dynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"customer_email": &types.AttributeValueMemberS{
				Value: customerEmail,
			},
			"photo_id": &types.AttributeValueMemberS{
				Value: photoID,
			},
			"upload_timestamp": &types.AttributeValueMemberN{
				Value: strconv.FormatInt(timestamp, 10),
			},
			"photo_status": &types.AttributeValueMemberS{
				Value: "uploaded",
			},
		},
	})

	if err != nil {
		log.Printf("unable to insert metadata for photo %s: %v", photoID, err)
	}

	return nil
}

// UpdateMetadata updates the metadata for a photo upload once processed
func UpdateMetadata(customerEmail, photoID, processedLocation string) error {
	_, err := dynamoClient.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME")),
		Key: map[string]types.AttributeValue{
			"customer_email": &types.AttributeValueMemberS{
				Value: customerEmail,
			},
			"photo_id": &types.AttributeValueMemberS{
				Value: photoID,
			},
		},
		UpdateExpression: aws.String("SET photo_status = :status, processed_location = :location"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":location": &types.AttributeValueMemberS{
				Value: processedLocation,
			},
			":status": &types.AttributeValueMemberS{
				Value: "processed",
			},
		},
	})

	if err != nil {
		log.Printf("unable to update metadata for photo %s: %v", photoID, err)
	}

	return nil
}
