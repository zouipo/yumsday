FROM golang:1.25-alpine AS base
RUN apk add --no-cache gcc=15.2.0-r2 make=4.4.1-r3 musl-dev=1.2.5-r21 && \
    go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY --exclude=test/compose.yaml --exclude=Dockerfile --exclude=Jenkinsfile --exclude=README.md . .


FROM base AS build
RUN make


FROM alpine:3.22 AS runtime
ENV USER="yumsday"
RUN addgroup -g 1000 -S $USER && \
    adduser -D -H -u 1000 -g 1000 -S -G $USER $USER
WORKDIR /app
COPY --from=build /app/bin/yumsday .
USER $USER
WORKDIR /data
VOLUME /data
ENTRYPOINT ["/app/yumsday"]
