BINARY_NAME=iollama
BUILD_DIR=bin

.PHONY: build test clean docker-build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/iollama

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

docker-build:
	docker build -t iollama:latest .
