package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const keySize = 32 // 32 bytes for AES-256

// generateRandomKey generates a random key of the specified size.
func generateRandomKey(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// handler is the Lambda function handler.
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key, err := generateRandomKey(keySize)
	if err != nil {
		log.Println("Error generating key:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Internal Server Error: %v", err),
			}, nil
	}
		
		// Encode the key in base64 for transmission
	encodedKey := base64.StdEncoding.EncodeToString(key)
	response := map[string]string{
		"key": encodedKey,
	}
	
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: encodeResponse(response),
	}, nil
}
		
// encodeResponse encodes a response map to a JSON string.
func encodeResponse(response map[string]string) string {
	json, err := json.Marshal(response)
	if err != nil {
		log.Println("Error encoding response:", err)
		return "{}"
	}
	return string(json)
}
	
func main() {
	lambda.Start(handler)
}
