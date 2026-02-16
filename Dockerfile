FROM golang:1.25-alpine AS base
RUN apk add --no-cache gcc=15.2.0-r2 musl-dev=1.2.5-r21 make
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .


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
