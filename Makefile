default: release

SHELL=/bin/bash
TEST_COVERAGE_THRESHOLD=88.7

fmt:
	@echo "Running go fmt"
	@go fmt

lint:
	@if [ ! -d /tmp/golangci-lint ]; then \
		echo "Installing golangci-lint"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.45.2; \
		mkdir -p /tmp/golangci-lint/; \
		mv ./bin/golangci-lint /tmp/golangci-lint/golangci-lint; \
	fi; \
	/tmp/golangci-lint/golangci-lint run ./... --issues-exit-code=1 \

tidy:
	@echo "Running go mod tidy"
	@go mod tidy

test:
	@echo "Running go test"
	@go test -race  ./...

coverage:
	@TEST_COVERAGE=$$(go test  -coverpkg ./... | grep coverage | grep -Eo '[0-9]+\.[0-9]+') ;\
	if [ $$(bc <<< "$$TEST_COVERAGE < $(TEST_COVERAGE_THRESHOLD)") -eq 1 ]; then \
		echo "Current test coverage $$TEST_COVERAGE is below threshold of $(TEST_COVERAGE_THRESHOLD)." ;\
		exit 1 ;\
	fi

release: fmt lint tidy test coverage
	@echo "Build Successful."