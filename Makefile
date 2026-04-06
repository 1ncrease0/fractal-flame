COVERAGE_FILE ?= coverage.out

TARGET ?= fractalflame # CHANGE THIS TO YOUR BINARY NAME

.PHONY: build
build:
	@echo "Выполняется go build для таргета ${TARGET}"
	@mkdir -p .bin
	@go build -o ./bin/${TARGET} ./cmd/${TARGET}

## test: run all tests
.PHONY: test
test:
	go test -coverpkg='./...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	go tool cover -func='$(COVERAGE_FILE)'
