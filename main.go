package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// Request structure for the JSON body
type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Global SES client to reuse connections across warm starts
var sesClient *ses.Client

func init() {
	// Initialize the SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	sesClient = ses.NewFromConfig(cfg)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. SECURITY: Check the Static API Key
	apiKey := os.Getenv("STATIC_API_KEY")
	requestApiKey := request.Headers["x-api-key"]

	// Handle case insensitivity for headers if necessary, but standard is usually exact match or lower-case
	if requestApiKey == "" {
		// Try lowercase if standard casing failed (API Gateway sometimes lowercases headers)
		requestApiKey = request.Headers["x-api-key"]
	}

	if apiKey == "" || requestApiKey != apiKey {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"message": "Unauthorized"}`,
		}, nil
	}

	// 2. PARSE: Decode the JSON body
	var emailReq EmailRequest
	err := json.Unmarshal([]byte(request.Body), &emailReq)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Invalid JSON body"}`,
		}, nil
	}

	// 3. EXECUTE: Send Email via SES
	senderEmail := os.Getenv("SENDER_EMAIL") // Must be verified in SES

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{emailReq.To},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(emailReq.Body),
				},
			},
			Subject: &types.Content{
				Data: aws.String(emailReq.Subject),
			},
		},
		Source: aws.String(senderEmail),
	}

	_, err = sesClient.SendEmail(ctx, input)
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message": "Failed to send email: %v"}`, err),
		}, nil
	}

	// 4. SUCCESS
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "Email sent successfully"}`,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
