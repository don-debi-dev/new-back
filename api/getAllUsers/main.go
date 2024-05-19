package main

import (
	"encoding/json"
	"new-back/queries"
	"new-back/rds"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
)

// type RequestBody struct {
// }

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// rb, err := getQueryParams(req)
	// if err != nil {
	// 	return queries.GetLambdaResponse(queries.ErrDefault)
	// }

	db, err := rds.GetDB()
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	users, err := queries.GetAllUserProfiles(db)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	marshalled, err := json.Marshal(users)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	defer db.Close()

	return queries.GetLambdaResponse(queries.SuccessGetUser, string(marshalled))
}

// func getQueryParams(req events.APIGatewayProxyRequest) (RequestBody, error) {
// 	rb := RequestBody{}

// 	if err := json.Unmarshal([]byte(req.Body), &rb); err != nil {
// 		log.Printf("Error parsing request payload Payload: %s\n", req.Body)
// 		return RequestBody{}, err
// 	}

// 	return rb, nil
// }
