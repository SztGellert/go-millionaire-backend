package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	response := events.APIGatewayProxyResponse{}

	// Switch for identifying the HTTP request
	switch request.HTTPMethod {
	case "GET":
		// Obtain the QueryStringParameter
		topic := request.QueryStringParameters["topic"]
		if topic != "" {
			response = events.APIGatewayProxyResponse{Body: "You requested " + topic + " question! ", StatusCode: 200}
		} else {
			response = events.APIGatewayProxyResponse{Body: "Error: Query Parameter topic missing", StatusCode: 500}
		}
	}
	// Response
	return response, nil

}
