NAME=yumsday
MAIN=main.go
OUT=bin/$(NAME)
COVERAGE_REPORT=test/coverage.out
COVERAGE_REPORT_HTML=test/coverage.html

.PHONY: all
all: build

.PHONY: front
front:
	@cd front && \
		npm install && \
		npm run build

.PHONY: swagger
swagger:
	@swag init

.PHONY: build
build: swagger front
	@go build -trimpath -ldflags="-s -w -extldflags='-static'" -tags "sqlite_omit_load_extension" -o $(OUT) $(MAIN)

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
run: swagger
	@go run -tags dev $(MAIN)

.PHONY: test
test: swagger
	@mkdir -p test
	@go test -tags dev -cover -coverprofile=$(COVERAGE_REPORT) ./...

.PHONY: test-ci
test-ci: swagger
	@mkdir -p test

	@# lint fails if there is compile error
	@# and there is a compile error if front/dist does not exist or is empty
	@# because it is embedded with //go:embed
	@mkdir -p front/dist
	@touch front/dist/placeholder

	@CGO_ENABLED=1 go test -tags dev -race -cover -coverprofile=$(COVERAGE_REPORT) ./...

.PHONY: benchmark
benchmark:
	@go test -tags dev -bench=. -benchmem -run =^a ./...

.PHONY: coverage
coverage: test
	@go tool cover -html=$(COVERAGE_REPORT) -o=$(COVERAGE_REPORT_HTML)
	@xdg-open $(COVERAGE_REPORT_HTML)

.PHONY: lint
lint: swagger
	@golangci-lint run

.PHONY: clean
clean:
	@rm -r bin

.PHONY: gitclean
gitclean:
	@git clean -xdf
