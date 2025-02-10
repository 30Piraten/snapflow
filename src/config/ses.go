package config

// / func SendEmail(ctx context.Context, recipient, body string) error {
// 	// Load AWS SDK configuration
// 	cfg, err := config.LoadDefaultConfig(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to load AWS SDK config: %v", err)
// 	}

// 	subject := "📸 Your Photo is Ready for Pickup!"
// 	senderEmail := os.Getenv("SENDER_EMAIL")

// 	// Validate sender email
// 	if senderEmail == "" {
// 		log.Fatal("SENDER_EMAIL environment variable is not set")
// 	}

// 	// Initialize SES client
// 	sesClient := ses.NewFromConfig(cfg)

// 	// Create the email request
// 	input := &ses.SendEmailInput{
// 		Destination: &sesTypes.Destination{
// 			ToAddresses: []string{recipient},
// 		},
// 		Message: &sesTypes.Message{
// 			Subject: &sesTypes.Content{
// 				Data: aws.String(subject),
// 			},
// 			Body: &sesTypes.Body{
// 				Html: &sesTypes.Content{
// 					Charset: aws.String(CharSet),
// 					Data:    aws.String(body),
// 				},
// 				Text: &sesTypes.Content{
// 					Charset: aws.String(CharSet),
// 					Data:    aws.String(body),
// 				},
// 			},
// 		},
// 		Source: aws.String(senderEmail),
// 	}

// 	// Send the email
// 	_, err = sesClient.SendEmail(ctx, input)
// 	if err != nil {
// 		log.Printf("Failed to send email: %v", err)
// 		return err
// 	}/
