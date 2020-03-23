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
	@go build -o ./bin/server ./cmd/main.go
