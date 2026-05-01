#!/bin/bash
set -e

SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(cd "$SCRIPTS_DIR/.." && pwd)"
APP_DIR="$DEPLOY_DIR/apps/cad-bot"

echo "=== Stopping CAD Bot ==="

cd "$APP_DIR"
docker-compose down

echo "=== CAD Bot Stopped ==="
