package main

import (
	"context"
	"encoding/json"
	"github.com/SztGellert/go-millionaire-backend/load_quiz"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	response := events.APIGatewayV2HTTPResponse{}

	switch request.RequestContext.HTTP.Method {
	case "POST":
		quizData, err := load_quiz.LoadQuizData(request.QueryStringParameters["topic"], request.QueryStringParameters["difficulty"], request.Body)
		if err != nil {
			response = events.APIGatewayV2HTTPResponse{Body: "Database error!", StatusCode: 500}
			return response, nil
		}
		questionJson, err := json.Marshal(quizData)
		if err != nil {
			response = events.APIGatewayV2HTTPResponse{Body: "Service error!", StatusCode: 500}
			return response, nil
		}
		response = events.APIGatewayV2HTTPResponse{Body: string(questionJson), StatusCode: 200}
	}

	return response, nil
}
