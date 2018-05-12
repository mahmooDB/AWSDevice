# Simple Restful API on AWS


## Project description
This project implements a simple Restful API on AWS using the following tech stack:

 * Serverless Framework ( https://serverless.com ) 
 * Go language ( https://golang.org )
 * AWS API Gateway
 * AWS Lambda
 * AWS DynamoDB


The API accepts the following JSON requests and produces the corresponding HTTP responses:



### Request 1:
```
HTTP POST
URL: https://`API-GATEWAY-URL`/api/devices
Body (application/json):
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}
```



#### Response 1 - Success:
```
HTTP 201 Created
Body (application/json):
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}
```



#### Response 1 - Failure 1:
```
HTTP 400 Bad Request
```
If any of the payload fields are missing. Response body should
 have a descriptive error message for the client to be able to
 detect the problem.



#### Response 1 - Failure 2:
```
HTTP 500 Internal Server Error
```
If any exceptional situation occurs on the server side.



### Request 2:
```
HTTP GET
URL: https://`API-GATEWAY-URL`/api/devices/{id}
Example: GET https://api123.amazonaws.com/api/devices/id1
```



#### Response 2 - Success:
```
HTTP 200 OK
Body (application/json):
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}
```



#### Response 2 - Failure 1:
```
HTTP 404 Not Found
```
If the request id does not exist.



#### Response 2 - Failure 2:
```
HTTP 500 Internal Server Error
```
If any exceptional situation occurs on the server side.


---

## Project Structure
Two AWS Lambda functions are responsible for mentioned tasks:
 1. `postNewDevice` for handling POST request
 2. `getDeviceInfo` for handling GET request


**Golang** code for these functions are stored in `src/postNewDevice/main.go` and `src/getDeviceInfo/main.go` files. Data structure for `Device` is shared for both lambda functions. The code for `Device` struct is avaialble at `src/vendor/data` folder.


Both of these functions connect to a DynamoDB table named "devices". This table and all of the necessary resources and permissions are described in `serverless.yml` file. The serverless will automatically generate required resources. Finally, unit tests and test results are stored in `tests` folder.

---

## Deployment
Serverless deployment requires following installed:
 * aws cli
 * go
 * go dep
 * node.js
 * serverless


Golang written AWS lambda functions must be compiled first before deployment. Therefore when you are inside the project folder run:
```
make
```
After a successful build, executables will be created in `bin` folder.
In order to deploy, run the following command:
```
serverless deploy
```
After successful deployment, two URLs will be generated. This is your deployment api URL that we call `API-GATEWAY-URL`.

---

## Testing
For testing by URL, replace you `API-GATEWAY-URL` with following `curl` commands:
### Create a new device:
```
curl -i -H "Content-Type: application/json" -X POST -d '{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}' https://API-GATEWAY-URL/devices
```
If you try to create another device using an existing `id`, the old item will be replaced by new item.

### Get an existing device by providing an id:
```
curl -i https://API-GATEWAY-URL/devices/id1
```
Also **unit tests** for code coverage are available inside `tests` folder. In order to test the functions yourself use the following command:
```
go test -coverprofile=cover.out
```
This will create a file named `cover.out`. To get a HTML representation of code coverage, use the following command after generating `cover.out`:
```
go tool cover -html=cover.out -o cover.html
```
Open the `cover.html` using a browser. Covered areas are shown in green. Coverage results from `go tool cover -func=cover.out` are presented here for both lamda functions:

#### getDeviceInfo:
```
device-db/tests/getDeviceInfo/main.go:23:	init		100.0%
device-db/tests/getDeviceInfo/main.go:35:	Handler		81.8%
device-db/tests/getDeviceInfo/main.go:117:	main		0.0%
total:						(statements)	80.0%
```

#### postNewDevice:
```
device-db/tests/postNewDevice/main.go:24:	init		100.0%
device-db/tests/postNewDevice/main.go:36:	createNewDevice	100.0%
device-db/tests/postNewDevice/main.go:82:	Handler		82.6%
device-db/tests/postNewDevice/main.go:163:	main		0.0%
total:						(statements)	89.6%
```

The `cover.html` files demonstrate that error handling codes for JSON marshaling/unmarshaling and also AWS session creation were not covered. However, it is still possible to obtain 100% coverage in both files. In order to achieve 100% coverage, one must implement mock methods that generate intentional errors for JSON marshaling/unmarshaling and AWS session as well.
