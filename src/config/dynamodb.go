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
)

var dynamoClient *dynamodb.Client

func InitDynamoDB() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("failed to load configuration, ", err)
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

// InsertMetadata inserts metadata for a new photo upload.
func InsertMetadata(customerFullName, customerEmail, paperType, paperSize, photoID string, timestamp int64) error {
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	_, err := dynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"customer_fullname": &types.AttributeValueMemberS{
				Value: customerFullName,
			},
			"customer_email": &types.AttributeValueMemberS{
				Value: customerEmail,
			},
			"paper_type": &types.AttributeValueMemberS{
				Value: paperType,
			},
			"paper_size": &types.AttributeValueMemberS{
				Value: paperSize,
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
