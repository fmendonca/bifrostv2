APP_NAME = go-libvirt-api

build:
	go build -o $(APP_NAME) cmd/main.go

run:
	go run cmd/main.go

lint:
	golangci-lint run ./...

clean:
	rm -f $(APP_NAME)
