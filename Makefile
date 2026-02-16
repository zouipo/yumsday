COVERAGE_FILE=coverage.out
OUT=bin/yumsday

.PHONY: all
all: build

.PHONY: build
build:
	@go build -ldflags="-s -w" -o $(OUT) main.go

.PHONY: image
image:
	@docker build --target runtime -t zouipo/yumsday:latest .

.PHONY: run
run: build
	@go run .

.PHONY: run-race
run-race: build-race
	@$(OUT) -v

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
