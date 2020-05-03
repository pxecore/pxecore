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

all_tests: lint test coverage ## All Tests

build: ## Build the binary file
	GOOS=darwin GOARCH=amd64 go build -o ./build/pxecore.darwin.amd64.server ./cmd/server/main.go
	GOOS=linux GOARCH=arm64 go build -o ./build/pxecore.linux.amd64.server ./cmd/server/main.go
	GOOS=linux GOARCH=arm go build -o ./build/pxecore.linux.arm.server ./cmd/server/main.go

github_release:
	echo '$(EVENT_DATA)'
	pwd