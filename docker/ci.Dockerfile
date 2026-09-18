FROM golangci/golangci-lint:latest-alpine AS ci

WORKDIR /app

COPY docker/install-build-dependencies.sh .
RUN ./install-build-dependencies.sh

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN make swagger && make front

ENTRYPOINT ["golangci-lint", "run"]
