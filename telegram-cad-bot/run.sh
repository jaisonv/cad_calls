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

# Default values (change these to your department)
BASE_URL="${CAD_BASE_URL:-https://southmiamipdfl.policetocitizen.com}"
AGENCY_ID="${CAD_AGENCY_ID:-386}"
CHECK_INTERVAL="${CHECK_INTERVAL:-5}"

echo "🤖 Starting Telegram CAD Bot"
echo "================================"
echo "Base URL:  $BASE_URL"
echo "Agency ID: $AGENCY_ID"
echo "Interval:  $CHECK_INTERVAL minutes"
echo "================================"
echo ""

# Run the bot
./cadbot \
  -token "$TELEGRAM_BOT_TOKEN" \
  -base-url "$BASE_URL" \
  -agency-id "$AGENCY_ID" \
  -interval "$CHECK_INTERVAL"
