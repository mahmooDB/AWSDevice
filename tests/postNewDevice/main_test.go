// postNewSensor test
package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Test struct that holds an input and its corresponding output
type Test struct {
	in  events.APIGatewayProxyRequest
	out int
}

// A mock struct that emulates DynamoDB
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

// Overriding the PutItem method for mock DynamoDB
func (d *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	// Just return an empty output (no error)
	return new(dynamodb.PutItemOutput), nil
}

// An array of tests (input output pairs)
// Input: An events.APIGatewayProxyRequest object
// output: Expected status code
var tests = []Test{
	{
		// Test: Bad input json
		events.APIGatewayProxyRequest{
			Body: "{{}",
		},
		500,
	},
	{
		// Test: Incomplete input json (missing other fields)
		events.APIGatewayProxyRequest{
			Body: "{\"id\":\"/devices/id1\"}",
		},
		400,
	},
	{
		// Test: Incomplete input json (missing other fields)
		events.APIGatewayProxyRequest{
			Body: "{\"deviceModel\":\"id1\"}",
		},
		400,
	},
	{
		// Test: Incomplete input json (missing other fields)
		events.APIGatewayProxyRequest{
			Body: "{\"name\":\"Sensor.\"}",
		},
		400,
	},
	{
		// Test: Incomplete input json (missing other fields)
		events.APIGatewayProxyRequest{
			Body: "{\"note\":\"Testing a sensor.\"}",
		},
		400,
	},
	{
		// Test: Incomplete input json (missing other fields)
		events.APIGatewayProxyRequest{
			Body: "{\"serial\":\"A020000102\"}",
		},
		400,
	},
	{
		// Test: valid input json
		events.APIGatewayProxyRequest{
			Body: "{\"id\":\"/devices/id1\", \"deviceModel\":\"id1\", \"name\":\"Sensor.\", \"note\":\"Testing a sensor.\", \"serial\":\"A020000102\"}",
		},
		201,
	},
}

// Actual test function
func TestHandler(t *testing.T) {
	// For every test (input output pair) do:
	for i, test := range tests {
		response, _ := Handler(test.in)
		// If expected output does not match the actual output, throw error
		if response.StatusCode != test.out {
			t.Errorf("#%d: Expected: %d, Actual: %d", i, test.out, response.StatusCode)
		}
	}
}
