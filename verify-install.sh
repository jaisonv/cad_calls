#!/bin/bash
#
# Verify installation of Telegram CAD Bot
#

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Verifying installation..."
echo ""

# Check Python
if command -v python3 &> /dev/null; then
    echo -e "${GREEN}✓${NC} Python 3: $(python3 --version)"
else
    echo -e "${RED}✗${NC} Python 3: Not installed"
fi

# Check Go
if command -v go &> /dev/null; then
    echo -e "${GREEN}✓${NC} Go: $(go version | awk '{print $3}')"
else
    echo -e "${RED}✗${NC} Go: Not installed"
fi

# Check config.py
if [ -f "config.py" ]; then
    if grep -q "REPLACE_WITH_YOUR_BASE_URL" config.py; then
        echo -e "${YELLOW}⚠${NC} config.py: Exists but not configured"
    else
        echo -e "${GREEN}✓${NC} config.py: Configured"
        grep "BASE_URL" config.py | head -1
        grep "AGENCY_ID" config.py | head -1
    fi
else
    echo -e "${RED}✗${NC} config.py: Not found"
fi

# Check Python dependencies
if python3 -c "import requests" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Python dependencies: Installed"
else
    echo -e "${RED}✗${NC} Python dependencies: Not installed"
fi

# Check bot binary
if [ -f "telegram-cad-bot/cadbot" ]; then
    echo -e "${GREEN}✓${NC} Bot binary: Built"
else
    echo -e "${RED}✗${NC} Bot binary: Not built"
fi

# Check start script
if [ -f "start-bot.sh" ]; then
    echo -e "${GREEN}✓${NC} Start script: Created"
else
    echo -e "${YELLOW}⚠${NC} Start script: Not found"
fi

# Check TELEGRAM_BOT_TOKEN
if [ -n "$TELEGRAM_BOT_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} TELEGRAM_BOT_TOKEN: Set"
else
    echo -e "${YELLOW}⚠${NC} TELEGRAM_BOT_TOKEN: Not set in environment"
fi

# Check systemd service
if [ -f "/etc/systemd/system/telegram-cad-bot.service" ]; then
    echo -e "${GREEN}✓${NC} Systemd service: Installed"
else
    echo -e "${YELLOW}⚠${NC} Systemd service: Not installed"
fi

echo ""
echo "Verification complete!"
