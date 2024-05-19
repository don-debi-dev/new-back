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
	Password string `json:"password"`
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rb, err := getQueryParams(req)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	login := queries.LoginDetails{
		UserName: rb.UserName,
		Password: rb.Password,
	}

	db, err := rds.GetDB()
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	newUser, err := queries.GetUserIdFromNameAndPassword(db, login)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	// newUser, err := queries.SetUserSession(db, rb)
	// if err != nil {
	// 	return queries.GetLambdaResponse(queries.ErrDefault)
	// }

	marshalled, err := json.Marshal(newUser)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	return queries.GetLambdaResponse(queries.SuccessLoginUser, string(marshalled))
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
