#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== Building Telegram CAD Bot ==="

echo "Downloading dependencies..."
go mod download

# brew install filosottile/musl-cross/musl-cross is needed to run this
echo "Building binary..."
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o cadbot ./cmd/bot

echo "Binary built successfully: cadbot"
echo ""
echo "Build target: GOOS=${GOOS:-linux} GOARCH=${GOARCH:-amd64} CGO_ENABLED=${CGO_ENABLED:-1}"
echo "To run locally (uses Python config file):"
echo "  ./cadbot -token YOUR_TELEGRAM_TOKEN -python-script ../direct_api_post.py -db ./cadbot.db -interval 5"
