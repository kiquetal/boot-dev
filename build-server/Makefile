.PHONY: build test clean

build:
	@echo "Building..."
	go build -o server main.go

test:
	@echo "Testing..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	go clean
	rm -f server

report:
	@echo "Generating test coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

