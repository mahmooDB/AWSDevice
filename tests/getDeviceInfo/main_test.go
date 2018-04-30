// postNewSensor test
package main

import (
	//"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Test struct {
	in  events.APIGatewayProxyRequest
	out int
}

var tests = []Test{
	{
		events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				"id": "",
			},
		},
		404,
	},
	{
		events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				"id": "idnotfound",
			},
		},
		404,
	},
	{
		events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				"id": "id1",
			},
		},
		200,
	},
}

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (d *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	out := dynamodb.GetItemOutput{}
	id := input.Key["id"].S
	if *id == "/devices/id1" {
		out.SetItem(
			map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{S: id},
			},
		)
	}
	return &out, nil
}

func TestHandler(t *testing.T) {
	for i, test := range tests {
		response, _ := Handler(test.in)
		if response.StatusCode != test.out {
			t.Errorf("#%d: Expected: %d, Actual: %d, %s", i, test.out, response.StatusCode, response.Body)
		}
	}
}
