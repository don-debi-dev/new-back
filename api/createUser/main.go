package main

import (
	"encoding/json"
	"errors"
	"log"
	"new-back/queries"
	"new-back/rds"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
)

type RequestBody struct {
	UserName string `json:"username"`
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rb, err := getQueryParams(req)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	db, err := rds.GetDB()
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	tempUser, err := queries.GetUserProfileFromName(db, rb.UserName)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	if tempUser.UserId > 0 {
		return queries.GetLambdaResponse(queries.ErrCreateUser)
	}

	_, err = queries.CreateUserProfile(db, rb.UserName)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	newUser, err := queries.GetUserProfileFromName(db, rb.UserName)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	defer db.Close()

	marshalled, err := json.Marshal(newUser)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	return queries.GetLambdaResponse(queries.SuccessCreateUser, string(marshalled))
}

func getQueryParams(req events.APIGatewayProxyRequest) (RequestBody, error) {
	rb := RequestBody{}

	if err := json.Unmarshal([]byte(req.Body), &rb); err != nil {
		log.Printf("Error parsing request payload Error: %s\n", err.Error())
		return RequestBody{}, err
	}

	if rb.UserName == "" {
		err := errors.New("username is empty")
		log.Printf("Error parsing request payload Error: %s\n", err.Error())
		return RequestBody{}, err
	}

	return rb, nil
}
