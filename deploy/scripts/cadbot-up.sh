#!/bin/bash
set -e

SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(cd "$SCRIPTS_DIR/.." && pwd)"
APP_DIR="$DEPLOY_DIR/apps/cad-bot"
DOCKER_BUNDLE_DIR="$DEPLOY_DIR/docker"
SOURCE_ROOT="$(cd "$DEPLOY_DIR/.." && pwd)"

echo "=== Starting CAD Bot ==="

mkdir -p "$APP_DIR" "$APP_DIR/data" "$APP_DIR/cad_calls"

if [ ! -f "$APP_DIR/docker-compose.yml" ]; then
    cp "$DOCKER_BUNDLE_DIR/docker-compose.yml" "$APP_DIR/docker-compose.yml"
fi

if [ ! -f "$APP_DIR/setup.sh" ]; then
    cp "$DOCKER_BUNDLE_DIR/setup.sh" "$APP_DIR/setup.sh"
fi

if [ ! -f "$APP_DIR/.env.example" ]; then
    cp "$DOCKER_BUNDLE_DIR/.env.example" "$APP_DIR/.env.example"
fi

if [ ! -f "$APP_DIR/.env" ]; then
    echo "Creating .env from template..."
    cp "$APP_DIR/.env.example" "$APP_DIR/.env"
fi

if [ ! -f "$APP_DIR/cadbot" ] && [ -f "$SOURCE_ROOT/telegram-cad-bot/cadbot" ]; then
    echo "Copying cadbot binary from source build output..."
    cp "$SOURCE_ROOT/telegram-cad-bot/cadbot" "$APP_DIR/cadbot"
fi

if [ ! -f "$APP_DIR/cad_calls/direct_api_post.py" ] && [ -f "$SOURCE_ROOT/direct_api_post.py" ]; then
    echo "Copying Python CAD client files..."
    cp "$SOURCE_ROOT/direct_api_post.py" "$APP_DIR/cad_calls/direct_api_post.py"
fi

if [ ! -f "$APP_DIR/cad_calls/config.py" ] && [ -f "$SOURCE_ROOT/config.py" ]; then
    cp "$SOURCE_ROOT/config.py" "$APP_DIR/cad_calls/config.py"
fi

cd "$APP_DIR"

source .env

if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "Error: TELEGRAM_BOT_TOKEN is not set"
    echo "Edit .env and add your Telegram bot token (get from @BotFather)"
    exit 1
fi

if [ -z "$CAD_BASE_URL" ] || [ "$CAD_BASE_URL" = "https://your-agency.policetocitizen.com" ]; then
    echo "Error: CAD_BASE_URL not configured"
    echo "Edit .env with your agency's Police-to-Citizen URL"
    exit 1
fi

if [ -z "$CAD_AGENCY_ID" ] || [ "$CAD_AGENCY_ID" = "0" ]; then
    echo "Error: CAD_AGENCY_ID not configured"
    echo "Edit .env with your agency's ID"
    exit 1
fi

if [ ! -d "cad_calls" ]; then
    echo "Error: cad_calls directory not found in $APP_DIR"
    exit 1
fi

if [ ! -f "cadbot" ]; then
    echo "Error: cadbot binary not found in $APP_DIR"
    echo "Build it first: cd telegram-cad-bot && ./build.sh"
    exit 1
fi

if [ ! -f "setup.sh" ]; then
    echo "Error: setup.sh not found in $APP_DIR"
    exit 1
fi

chmod +x cadbot setup.sh

if [ ! -d "data" ] && [ -d "cad_calls/telegram-cad-bot/data" ]; then
    echo "Importing database from existing installation..."
    mkdir -p data
    cp cad_calls/telegram-cad-bot/data/cadbot.db data/
    echo "Database imported"
fi

echo "Starting CAD Bot..."
docker-compose up -d

echo ""
echo "=== CAD Bot Started ==="
echo "Logs: docker-compose logs -f cadbot"
