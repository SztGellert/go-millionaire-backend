package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SztGellert/go-millionaire-backend/load_quiz"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	load_quiz.ConnectMongo()

	response := events.APIGatewayV2HTTPResponse{}

	quizData, err := load_quiz.LoadQuizData(request.QueryStringParameters["topic"], request.QueryStringParameters["difficulty"])
	if err != nil {
		response = events.APIGatewayV2HTTPResponse{Body: "Database error!", StatusCode: 500}
		return response, nil
	}
	questionJson, err := json.Marshal(quizData)
	if err != nil {
		response = events.APIGatewayV2HTTPResponse{Body: "Service error!", StatusCode: 500}
		return response, nil

	}

	// Switch for identifying the HTTP request
	switch request.RequestContext.HTTP.Method {
	case "GET":
		// Obtain the QueryStringParameter

		topic := request.QueryStringParameters["topic"]
		fmt.Println("Its GET METHOD:" + topic)

		if topic != "" {
			response = events.APIGatewayV2HTTPResponse{Body: string(questionJson), StatusCode: 200}
		} else {
			response = events.APIGatewayV2HTTPResponse{Body: "Error: Query Parameter topic missing", StatusCode: 500}
		}

	}

	// Response
	return response, nil

}
