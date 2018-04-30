// postNewDevice
package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// struct for storing a device object
type Device struct {
	Id          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

// createNewDevice: takes a string mapped json and checks for requested fields
// if all of required fields are provided returns a new Device object
// Input: a string mapped json
// Output: a new Device object
func createNewDevice(jsonMap map[string]interface{}) (Device, error) {

	// Check if input fields are missing
	errFlag := false
	errMsg := "Bad request\nMissing following field(s):"
	if jsonMap["id"] == nil {
		errMsg = errMsg + "\n id"
		errFlag = true
	}
	if jsonMap["deviceModel"] == nil {
		errMsg = errMsg + "\n deviceModel"
		errFlag = true
	}
	if jsonMap["name"] == nil {
		errMsg = errMsg + "\n name"
		errFlag = true
	}
	if jsonMap["note"] == nil {
		errMsg = errMsg + "\n note"
		errFlag = true
	}
	if jsonMap["serial"] == nil {
		errMsg = errMsg + "\n serial"
		errFlag = true
	}
	// If any fields are missing, return error
	errMsg = errMsg + "\n"
	if errFlag == true {
		return Device{}, errors.New(errMsg)
	}

	// Create new device object otherwise
	newDevice := Device{Id: jsonMap["id"].(string),
		DeviceModel: jsonMap["deviceModel"].(string),
		Name:        jsonMap["name"].(string),
		Note:        jsonMap["note"].(string),
		Serial:      jsonMap["serial"].(string),
	}

	return newDevice, nil
}

// Handler: responsible for taking POST requests from user that provide a json
// containing new Device object and store it in DynamoDB table
// Input: json containing new device information
// Output: json echo of newly inserted device
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Try to unmarshal input json to a map
	var mappedJson map[string]interface{}
	err := json.Unmarshal([]byte(request.Body), &mappedJson)
	// If the input json is not a valid json, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 1\nJson unmarshaling error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Try to create a new Device struct form input
	newDevice, err := createNewDevice(mappedJson)
	// If the input was not compliant with predefined struct, return error 400
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Initialize an aws session
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: &region},
	)
	// If somthing went wrong with session creation, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 2\nSession error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}
	// Create DynamoDB client
	db := dynamodb.New(sess)

	// Try to prepare a DynamoDB item structure from Device
	dbDevice, err := dynamodbattribute.MarshalMap(newDevice)
	// If somthing went wrong during data prepration, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 3\nDatabase unmarshaling error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Get table name from OS
	tableName := aws.String(os.Getenv("DEVICES_TABLE_NAME"))
	// Prepare a structure for putting the new item
	dbInput := &dynamodb.PutItemInput{
		Item:      dbDevice,
		TableName: tableName,
	}
	// Try to put new item into DynamoDB table
	_, err = db.PutItem(dbInput)
	// If somthing went wrong with database, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 4\nDatabase error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Create a json response from newly inserted device
	jsonResponse, err := json.Marshal(newDevice)
	// If somthing went wrong with json creation, return error 500
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error 5\nJson marshaling error: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Finally, everything went smoothly! Echo new device information.
	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponse),
		StatusCode: 201,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
