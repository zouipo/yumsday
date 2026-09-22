NAME=yumsday
MAIN=main.go
COVERAGE_REPORT=test/coverage.out
COVERAGE_REPORT_HTML=test/coverage.html

BACKEND_SOURCES=$(shell find -type f -name "*.go" -not -path "./docs/*" -not -path "./front/node_modules/*")
FRONT_SOURCES=$(shell find -type f -path "./front/*" -not -path "./front/node_modules/*" -not -path "./front/dist/*" -not -name "*.go")
SWAGGER_SOURCES=$(shell find -type f -name "*.go" -path "./backend/internal/handler/*" -not -name "*_test.go")

BACKEND_OUT=bin/$(NAME)
FRONT_OUT=front/dist/index.html
SWAGGER_OUT=docs/docs.go

.PHONY: all
all: $(BACKEND_OUT)

$(BACKEND_OUT): $(BACKEND_SOURCES) $(FRONT_OUT) $(SWAGGER_OUT)
	@go build -trimpath -ldflags="-s -w -extldflags='-static'" -tags "sqlite_omit_load_extension" -o $(BACKEND_OUT) $(MAIN)

$(FRONT_OUT): $(FRONT_SOURCES)
	@cd front && \
		npm install && \
		npm run build

$(SWAGGER_OUT): $(SWAGGER_SOURCES)
	@swag init

.PHONY: image
image:
	@docker build --file docker/Dockerfile --target runtime --tag $(NAME):latest .

.PHONY: image-ci
image-ci:
	@docker build --file docker/ci.Dockerfile --tag $(NAME)-ci:latest .

.PHONY: compose-up
compose-up:
	@mkdir -p test/data
	@docker compose -f test/compose.yaml up
	@make compose-down

.PHONY: compose-down
compose-down:
	@docker compose -f test/compose.yaml down

.PHONY: run
run: $(SWAGGER_OUT)
	@go run $(MAIN)

.PHONY: dev
dev: $(SWAGGER_OUT)
	@go run -tags dev $(MAIN)

.PHONY: test
test: $(SWAGGER_OUT)
	@mkdir -p test
	@go test -tags dev -cover -coverprofile=$(COVERAGE_REPORT) ./...

.PHONY: test-ci
test-ci: $(SWAGGER_OUT)
	@mkdir -p test
	@CGO_ENABLED=1 go test -tags dev -race -cover -coverprofile=$(COVERAGE_REPORT) ./...

.PHONY: benchmark
benchmark:
	@go test -tags dev -bench=. -benchmem -run =^a ./...

.PHONY: coverage
coverage: test
	@go tool cover -html=$(COVERAGE_REPORT) -o=$(COVERAGE_REPORT_HTML)
	@xdg-open $(COVERAGE_REPORT_HTML)

.PHONY: lint
lint: $(SWAGGER_OUT)
	@# lint fails if there is compile error
	@# and there is a compile error if front/dist does not exist or is empty
	@# because it is embedded with //go:embed
	@mkdir -p front/dist
	@touch front/dist/placeholder

	@golangci-lint run

.PHONY: fmt
fmt:
	@golangci-lint fmt

.PHONY: clean
clean:
	@rm -r bin

.PHONY: gitclean
gitclean:
	@git clean -xdf
