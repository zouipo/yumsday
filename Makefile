OUT=bin/yumsday
COVERAGE_FILE=test/coverage.out
TEST_REPORT=test/test-report.json

.PHONY: all
all: build

.PHONY: swagger
swagger:
	@swag init

.PHONY: build
build: swagger
	@go build -ldflags="-s -w" -o $(OUT) main.go

.PHONY: image
image:
	@docker build --target runtime -t zouipo/yumsday:latest .

.PHONY: run
run: swagger
	@go run .

.PHONY: test
test: swagger
	@go test -cover -coverprofile=$(COVERAGE_FILE) ./...

.PHONY: test-cicd
test-cicd: swagger
	@go test -v -race -cover -coverprofile=$(COVERAGE_FILE) -json ./... > $(TEST_REPORT)

.PHONY: benchmark
benchmark:
	@go test -bench=. -benchmem -run =^a ./...

.PHONY: coverage
coverage: test
	@go tool cover -html=$(COVERAGE_FILE)

.PHONY: clean
clean:
	@rm -rf bin/
