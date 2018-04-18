package main_test

import (
	"testing"

	"github.com/hussfelt/dr-happy-push"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	Headers := make(map[string]string)
	Headers["Auth"] = "banankontakt"

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "Green", Headers: Headers},
			expect:  "3",
			err:     nil,
		},
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "Yellow", Headers: Headers},
			expect:  "2",
			err:     nil,
		},
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "Red", Headers: Headers},
			expect:  "1",
			err:     nil,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "", Headers: Headers},
			expect:  "",
			err:     main.ErrNameNotProvided,
		},
		{
			// Test that the handler responds ErrAuthenticationFailed
			// when no auth is posted
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  "",
			err:     main.ErrAuthenticationFailed,
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}