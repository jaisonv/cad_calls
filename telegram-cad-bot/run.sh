#!/bin/bash
# Simple script to run the Telegram CAD bot

# Check if TELEGRAM_BOT_TOKEN is set
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "❌ Error: TELEGRAM_BOT_TOKEN environment variable is not set"
    echo ""
    echo "Set it with:"
    echo "  export TELEGRAM_BOT_TOKEN='your_token_here'"
    echo ""
    echo "Or run with:"
    echo "  TELEGRAM_BOT_TOKEN='your_token' ./run.sh"
    exit 1
fi

# Check if Python script exists
PYTHON_SCRIPT="../direct_api_post.py"
if [ ! -f "$PYTHON_SCRIPT" ]; then
    echo "❌ Error: Python script not found at $PYTHON_SCRIPT"
    echo ""
    echo "Make sure you're running this from the telegram-cad-bot directory"
    echo "and that direct_api_post.py exists in the parent directory"
    exit 1
fi

# Check if config.py is configured
if grep -q "REPLACE_WITH_YOUR_BASE_URL" ../config.py; then
    echo "⚠️  Warning: config.py still has placeholder values"
    echo ""
    echo "Please edit ../config.py and set:"
    echo "  BASE_URL = \"https://yourpd.policetocitizen.com\""
    echo "  AGENCY_ID = 123"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

CHECK_INTERVAL="${CHECK_INTERVAL:-5}"

echo "🤖 Starting Telegram CAD Bot"
echo "================================"
echo "Python Script: $PYTHON_SCRIPT"
echo "Config: ../config.py"
echo "Interval: $CHECK_INTERVAL minutes"
echo "================================"
echo ""
echo "⚠️  Make sure config.py is configured!"
echo ""

# Run the bot
./cadbot \
  -token "$TELEGRAM_BOT_TOKEN" \
  -interval "$CHECK_INTERVAL"
