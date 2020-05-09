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

build_ipxe:
	apt install build-essential git liblzma-dev unzip wget coreutils isolinux
	cd ./pkg/ipxe && go run generator.go

.APP_NAME=pxecore
.BUILD_EXTENSION=$(if $(findstring windows, $(GOOS)),.exe,)
package: ## Packages aplication. Extra Vars: GOOS,GOARCH
	go build -o ./build/$(.APP_NAME)$(.BUILD_EXTENSION)

.GOOS=$(if $(GOOS),$(GOOS),linux)
.GOARCH=$(if $(GOARCH),$(GOARCH),amd64)
.RELEASE_NAME=$(if $(GITHUB_TAG_NAME),$(GITHUB_TAG_NAME),v.0.0.0)
.FLAVOUR_EXTENSION=$(if $(findstring windows, $(.GOOS)),.exe,)
.FLAVOUR_FILENAME=$(.APP_NAME)_$(.RELEASE_NAME)_$(.GOOS)_$(.GOARCH)$(.FLAVOUR_EXTENSION)
package_flavour: ## Packages aplication. Extra Vars: GOOS,GOARCH
	@echo Packaging Application...
	GOOS=$(.GOOS) GOARCH=$(.GOARCH) go build -o ./build/pxecore
	@echo Packaging Compleate!

.SUBST:={?name,label}
.RELEASE_UPLOAD_URL=$(subst $(.SUBST),,$(GITHUB_UPLOAD_URL))
.GITHUB_TOKEN=$(GITHUB_TOKEN)
github_release: package_flavour
	@echo Uploading Package...
	@tar cvfz "./build/$(.FLAVOUR_FILENAME).tar.gz" "pxecore"
	@curl \
      -X POST \
      --data-binary @./build/$(.FLAVOUR_FILENAME).tar.gz \
      -H 'Content-Type: application/gzip' \
      -H "Authorization: Bearer $(.GITHUB_TOKEN)" \
      "$(.RELEASE_UPLOAD_URL)?name=$(.FLAVOUR_FILENAME).tar.gz"