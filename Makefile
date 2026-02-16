COVERAGE_FILE=coverage.out
OUT=bin/yumsday

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
test:
	@go test -cover -coverprofile=$(COVERAGE_FILE) ./...

.PHONY: test-cicd
test-cicd:
	@go test -race ./...

.PHONY: benchmark
benchmark:
	@go test -bench=. -benchmem -run =^a ./...

.PHONY: coverage
coverage: test
	@go tool cover -html=$(COVERAGE_FILE)

.PHONY: clean
clean:
	@rm -rf bin/
