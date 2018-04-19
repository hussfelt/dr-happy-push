package main

import (
	"errors"
	"log"
	"strconv"
	"os"
	"fmt"

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

	// ErrCWNoSuccess is thrown when we could not send to CW
	ErrCWNoSuccess = errors.New("could not pass data to CloudWatch")
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// What is the current location of the box
	currentLocation := fmt.Sprintf("%s", os.Getenv("LOCATION"))

	if currentLocation == "" {
		currentLocation = "Office"
	}

	// Define map of statuses
	StatusMap := make(map[string]int)
	StatusMap["Green"] = 3
	StatusMap["Yellow"] = 2;
	StatusMap["Red"] = 1;

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

	// Trigger with ButtonPush
	var CWButtonPush bool
	CWButtonPush = sendToCloudWatch("ButtonPush", float64(1), currentLocation)

	// Break if error
	if ! CWButtonPush {
		return events.APIGatewayProxyResponse{}, ErrCWNoSuccess
	}

	// Trigger with ButtonPushAvg
	var CWButtonPushAvg bool
	CWButtonPushAvg = sendToCloudWatch("ButtonPushAvg", float64(StatusMap[request.Body]), currentLocation)

	// Break if error
	if ! CWButtonPushAvg {
		return events.APIGatewayProxyResponse{}, ErrCWNoSuccess
	}

	// Trigger with ButtonPushColor
	var CWButtonPushValue bool
	CWButtonPushValue = sendToCloudWatch("ButtonPush" + request.Body, float64(1), currentLocation)

	// Break if error
	if ! CWButtonPushValue {
		return events.APIGatewayProxyResponse{}, ErrCWNoSuccess
	}

	// Log the passed value
	log.Printf("Passing to CloudWatch %s\n", strconv.Itoa(StatusMap[request.Body]))

	return events.APIGatewayProxyResponse{
		Body:       strconv.Itoa(StatusMap[request.Body]),
		StatusCode: 200,
	}, nil

}

func sendToCloudWatch(metricName string, metricValue float64, currentLocation string) (bool) {

	// Setup new session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	// Create CloudWatch client
	cw := cloudwatch.New(sess)

	// Push metric to cloudwatch
	result, err := cw.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(metricName),
				Unit:       aws.String(cloudwatch.StandardUnitCount),
				Value:      aws.Float64(metricValue),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("Location"),
						Value: aws.String(currentLocation),
					},
				},
			},
		},
		Namespace: aws.String("HappyButton"),
	})

	if err != nil {
		log.Printf("Error in CW request %s\n", err)
		return false
	}

	log.Printf("Success: %s\n", result)

	return true
}

func main() {
	lambda.Start(Handler)
}

// Klarrt å betartt