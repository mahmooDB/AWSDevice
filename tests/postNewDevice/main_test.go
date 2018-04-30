// postNewSensor test
package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Test struct {
	in  events.APIGatewayProxyRequest
	out int
}

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (d *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return new(dynamodb.PutItemOutput), nil
}

var tests = []Test{
	{
		events.APIGatewayProxyRequest{
			Body: "{{}",
		},
		500,
	},
	{
		events.APIGatewayProxyRequest{
			Body: "{\"id\":\"/devices/id1\"}",
		},
		400,
	},
	{
		events.APIGatewayProxyRequest{
			Body: "{\"deviceModel\":\"id1\"}",
		},
		400,
	},
	{
		events.APIGatewayProxyRequest{
			Body: "{\"name\":\"Sensor.\"}",
		},
		400,
	},
	{
		events.APIGatewayProxyRequest{
			Body: "{\"note\":\"Testing a sensor.\"}",
		},
		400,
	},
	{
		events.APIGatewayProxyRequest{
			Body: "{\"serial\":\"A020000102\"}",
		},
		400,
	},
	{
		events.APIGatewayProxyRequest{
			Body: "{\"id\":\"/devices/id1\", \"deviceModel\":\"id1\", \"name\":\"Sensor.\", \"note\":\"Testing a sensor.\", \"serial\":\"A020000102\"}",
		},
		201,
	},
}

func TestHandler(t *testing.T) {
	for i, test := range tests {
		response, _ := Handler(test.in)
		if response.StatusCode != test.out {
			t.Errorf("#%d: Expected: %d, Actual: %d", i, test.out, response.StatusCode)
		}
	}
}
