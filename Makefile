# Makefile for car-pooling-challenge
# Licence MIT

dep: ## Get the dependencies
	@go get -v -d ./...
	@go get -u golang.org/x/lint/golint

lint: ## Lint the files
	@golint -set_exit_status ./...

test: ## Run unit-tests
	@go test -short ./...

coverage: ## Generate global code coverage report
	@go test -covermode=count -coverprofile ./coverage.cov ./...
	@go tool cover -func=./coverage.cov
	@rm ./coverage.cov

all_tests: dep lint test coverage ## All Tests

build: dep ## Build the binary file
	GOOS=darwin GOARCH=amd64 go build -o ./build/pxecore.server.darwin.amd64 ./cmd/server/main.go
	GOOS=linux GOARCH=arm64 go build -o ./build/pxecore.server.linux.amd64 ./cmd/server/main.go
	GOOS=linux GOARCH=arm go build -o ./build/pxecore.server.linux.arm ./cmd/server/main.go
