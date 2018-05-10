build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/postNewDevice src/postNewDevice/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getDeviceInfo src/getDeviceInfo/main.go
