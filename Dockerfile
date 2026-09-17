FROM golang:1.27.1-alpine AS base

WORKDIR /app
RUN apk add --no-cache gcc make musl-dev npm && \
    go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

FROM base AS build
RUN make

FROM gcr.io/distroless/static-debian13:nonroot AS runtime

COPY --from=build /app/bin/yumsday /yumsday
WORKDIR /data
VOLUME ["/data"]

ENTRYPOINT ["/yumsday"]
