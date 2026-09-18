#!/bin/sh

set -e

apk add --no-cache gcc make musl-dev npm
go install github.com/swaggo/swag/cmd/swag@latest
