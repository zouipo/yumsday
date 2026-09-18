FROM golang:1.27.1-alpine AS base

WORKDIR /app

COPY docker/install-build-dependencies.sh .
RUN ./install-build-dependencies.sh

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN make

FROM gcr.io/distroless/static-debian13:nonroot AS runtime

COPY --from=base /app/bin/yumsday /yumsday
WORKDIR /data
VOLUME ["/data"]

ENTRYPOINT ["/yumsday"]
