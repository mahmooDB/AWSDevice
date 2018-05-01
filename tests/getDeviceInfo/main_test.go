// getDeviceInfo test
package main

import (
	//"errors"
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


// An array of tests (input output pairs)
// Input: An events.APIGatewayProxyRequest object
// output: Expected status code
var tests = []Test{
	{
                // Test: id field is empty
		events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				"id": "",
			},
		},
		404,
	},
	{
                // Test: requested id does not exist
		events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				"id": "idnotfound",
			},
		},
		404,
	},
	{
                // Test: requested id exists
		events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				"id": "id1",
			},
		},
		200,
	},
}

// A mock struct that emulates DynamoDB
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

// Overriding the GetItem method for mock DynamoDB
func (d *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
        // Construct an empty output object
	out := dynamodb.GetItemOutput{}
	id := input.Key["id"].S
        // If the requested id exists in table, fill the output
	if *id == "/devices/id1" {
		out.SetItem(
			map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{S: id},
			},
		)
	}
        // Otherwise, just return the empty outout
	return &out, nil
}


// Actual test function
func TestHandler(t *testing.T) {
        // For every test (input output pair) do:
	for i, test := range tests {
		response, _ := Handler(test.in)
                // If expected output does not match the actual output, throw error
		if response.StatusCode != test.out {
			t.Errorf("#%d: Expected: %d, Actual: %d, %s", i, test.out, response.StatusCode, response.Body)
		}
	}
}
