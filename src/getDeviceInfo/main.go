// getDeviceInfo
package main

import (
	"encoding/json"
	"os"

	"data"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Global AWS session variables
var sess *session.Session
var errSess error

// Init function creates an globally accessible AWS session
func init() {
	// Initialize an aws session
	region := os.Getenv("AWS_REGION")
	sess, errSess = session.NewSession(&aws.Config{
		Region: &region},
	)
}

// Handler: responsible for taking GET requests from user that provide a id by
// path parameter and produce appropriate response.
// Input: id (provided by GET/path)
// Output: json containing device information
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Try to get id from request
	id := request.PathParameters["id"]
	// If the input id is empty, return error 404
	if id == "" {
		return events.APIGatewayProxyResponse{
			Body:       "No id provided!",
			StatusCode: 404,
		}, nil
	}

	// Add /devices/ path to the begining of input id
	id = "/devices/" + id

	// If somthing went wrong with session creation, return error 500
	if errSess != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 1\nSession error: " + errSess.Error(),
			StatusCode: 500,
		}, nil
	}
	// Create DynamoDB client
	db := dynamodb.New(sess)

	// Get table name from OS
	tableName := aws.String(os.Getenv("DEVICES_TABLE_NAME"))
	// Try to get the requested item from DynamoDB table
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	// If somthing went wrong with database, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 2\nDatabase error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}
	// If the requested item was not found, return error 404
	if len(result.Item) == 0 {
		return events.APIGatewayProxyResponse{
			Body:       "Not found!",
			StatusCode: 404,
		}, nil
	}

	// Create a device object from result
	device := data.Device{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &device)
	// If somthing went wrong with unmarshaling, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 3\nDatabase unmarshaling error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	jsonResponse, err := json.Marshal(device)
	// If somthing went wrong with json creation, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 4\nJson marshaling error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Finally, everything went smoothly! return requested device information.
	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponse),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
