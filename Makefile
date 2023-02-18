GOCLEAN=go clean
GOTEST=go test
GOLINT=golangci-lint
GOFMT=gofmt

PACKAGE_PATH=./mailtrap

TEST_FLAGS=-race -v
TEST_COVERAGE_FLAGS=-race -coverprofile=coverage.out -covermode=atomic

all: test lint

test:
	$(GOTEST) $(TEST_FLAGS) $(PACKAGE_PATH)

cover:
	$(GOTEST) $(TEST_COVERAGE_FLAGS) -coverprofile coverage.txt -covermode atomic $(PACKAGE_PATH)

cover-html:
	$(GOTEST) $(TEST_COVERAGE_FLAGS) $(PACKAGE_PATH)
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

lint:
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOLINT) run $(PACKAGE_PATH)

fmt:
	$(GOFMT) -s -w -l .

clean:
	$(GOCLEAN)
	rm -f coverage.out coverage.html

.PHONY: clean test lint cover cover-html fmt
