build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/postNewDevice postNewDevice/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getDeviceInfo getDeviceInfo/main.go
