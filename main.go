package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var gDB *sql.DB

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := applyConfiguration(); err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf("{\"msg\": \"%v\"}", err),
			StatusCode: 200,
		}, err
	}

	slog.Debug("got request", "body", request.Body)

	slog.Debug("extracting the request's data")
	extractedValues, err := extractData([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf("{\"msg\": \"%v\"}", err),
			StatusCode: 200,
		}, nil
	}

	insertDataPoint(gDB, extractedValues)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"msg": "ok"}`,
		StatusCode: 200,
	}, nil
}

// It turns out reading pays off: https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtime-environment.html
// The bottom line to the lifecycle of a lambda is that they indeed have some persistence whatsoever!
// We'll just initialize the DB once in the 'static' portion of the lambda and simply rely on a global
// handle to it. Doing so is the recommended way of working with GORM BTW...
func main() {
	db, err := openDB()

	if err != nil {
		slog.Error("couldn't open the DB", "err", err)
	}

	gDB = db

	lambda.Start(handler)
}
