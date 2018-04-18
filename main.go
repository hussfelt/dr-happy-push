package main

import (
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")

	// ErrAuthenticationFailed is thrown when a name is not provided
	ErrAuthenticationFailed = errors.New("recieved an invalid authentication token")

	// ErrBodyNotAccepted is thrown when a body is not provided
	ErrBodyNotAccepted = errors.New("recieved an invalid authentication token")

	// ErrCouldNotConvert is thrown when we could not convert body to float
	ErrCouldNotConvert = errors.New("could not convert body to float")

	// ErrCWNoSuccess is thrown when we could not send to CW
	ErrCWNoSuccess = errors.New("could not pass data to CloudWatch")
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Setup new session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	// Create CloudWatch client
	cw := cloudwatch.New(sess)

	// Define map of statuses
	StatusMap := make(map[string]string)
	StatusMap["Green"] = "3";
	StatusMap["Yellow"] = "2";
	StatusMap["Red"] = "1";

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	if request.Body != "Green" && request.Body != "Yellow" && request.Body != "Red" {
		return events.APIGatewayProxyResponse{}, ErrBodyNotAccepted
	}

	// Log headers
	log.Printf("Passed headers %s\n", request.Headers)
	if request.Headers["Auth"] != "banankontakt" {
		return events.APIGatewayProxyResponse{}, ErrAuthenticationFailed
	}

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	// Convert value to float
    floatBody, err := strconv.ParseFloat(StatusMap[request.Body][:len(StatusMap[request.Body])-1], 64)
    if err == nil {
        return events.APIGatewayProxyResponse{}, ErrCouldNotConvert
    }

    // Log the passed value
    log.Printf("Passing to CloudWatch %s\n", floatBody)

	// Push metric to cloudwatch
	result, err := cw.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String("ButtonPush"),
				Unit:       aws.String(cloudwatch.StandardUnitCount),
				Value:      aws.Float64(floatBody),
			},
		},
		Namespace: aws.String("HappyButton"),
	})

    if err != nil {
    	log.Printf("Error in CW request %s\n", err)
        return events.APIGatewayProxyResponse{}, ErrCWNoSuccess
    }

    log.Printf("Success: %s\n", result)

	return events.APIGatewayProxyResponse{
		Body:       StatusMap[request.Body],
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}