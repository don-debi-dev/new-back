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
	UserId int64 `json:"userid"`
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

	user, err := queries.GetUserProfileFromId(db, rb.UserId)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}
	if user.UserId == 0 {
		return queries.GetLambdaResponse(queries.ErrGetUser)
	}

	_, err = queries.DeleteUserProfileById(db, rb.UserId)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	marshalled, err := json.Marshal(user)
	if err != nil {
		return queries.GetLambdaResponse(queries.ErrDefault)
	}

	defer db.Close()

	return queries.GetLambdaResponse(queries.SuccessDeleteUser, string(marshalled))
}

func getQueryParams(req events.APIGatewayProxyRequest) (RequestBody, error) {
	rb := RequestBody{}

	if err := json.Unmarshal([]byte(req.Body), &rb); err != nil {
		log.Printf("Error parsing request payload Payload: %s\n", req.Body)
		return RequestBody{}, err
	}

	if rb.UserId == 0 {
		err := errors.New("userid is 0")
		log.Printf("Error in payload Error: %s\n", err.Error())
		return RequestBody{}, err
	}

	return rb, nil
}
