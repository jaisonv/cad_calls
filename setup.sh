#!/bin/bash
#
# Telegram CAD Calls Monitor Bot - Interactive Installer
# This script sets up the bot from scratch on a Linux machine
#

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}"
echo "=================================================="
echo "  Telegram CAD Calls Monitor Bot - Setup"
echo "=================================================="
echo -e "${NC}"
echo ""

# Function to print status messages
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[i]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Check if running as root for system-wide installation
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_warning "Running as root. This will install system-wide."
        INSTALL_SYSTEMD=true
    else
        print_info "Running as regular user. Systemd service will be optional."
        INSTALL_SYSTEMD=false
    fi
}

# Check and install Python 3
install_python() {
    print_info "Checking for Python 3..."

    if command -v python3 &> /dev/null; then
        PYTHON_VERSION=$(python3 --version | awk '{print $2}')
        print_status "Python 3 is already installed (version $PYTHON_VERSION)"
    else
        print_warning "Python 3 not found. Installing..."

        if command -v apt-get &> /dev/null; then
            sudo apt-get update
            sudo apt-get install -y python3 python3-pip python3-venv
        elif command -v yum &> /dev/null; then
            sudo yum install -y python3 python3-pip
        elif command -v dnf &> /dev/null; then
            sudo dnf install -y python3 python3-pip
        else
            print_error "Could not detect package manager. Please install Python 3 manually."
            exit 1
        fi

        print_status "Python 3 installed successfully"
    fi
}

# Check and install Go
install_go() {
    print_info "Checking for Go..."

    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_status "Go is already installed ($GO_VERSION)"
    else
        print_warning "Go not found. Installing Go 1.21..."

        # Detect architecture
        ARCH=$(uname -m)
        if [ "$ARCH" = "x86_64" ]; then
            GO_ARCH="amd64"
        elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
            GO_ARCH="arm64"
        elif [ "$ARCH" = "armv7l" ] || [ "$ARCH" = "armv6l" ]; then
            GO_ARCH="armv6l"
            print_info "Detected Raspberry Pi (32-bit ARM)"
        else
            print_error "Unsupported architecture: $ARCH"
            exit 1
        fi

        # Download and install Go
        GO_VERSION="1.21.5"
        GO_TAR="go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"

        cd /tmp
        wget -q "https://go.dev/dl/${GO_TAR}" || {
            print_error "Failed to download Go"
            exit 1
        }

        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf "${GO_TAR}"
        rm "${GO_TAR}"

        # Add Go to PATH
        if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        fi
        export PATH=$PATH:/usr/local/go/bin

        print_status "Go installed successfully"
    fi
}

# Interactive configuration
configure_bot() {
    echo ""
    echo -e "${BLUE}=================================================="
    echo "  Configuration Setup"
    echo "==================================================${NC}"
    echo ""

    # Get Telegram Bot Token
    echo -e "${YELLOW}1. Telegram Bot Token${NC}"
    echo "   Get your token from @BotFather on Telegram"
    echo "   (Create a new bot with /newbot if you don't have one)"
    echo ""
    read -p "Enter your Telegram Bot Token: " TELEGRAM_TOKEN
    while [ -z "$TELEGRAM_TOKEN" ]; do
        print_error "Token cannot be empty!"
        read -p "Enter your Telegram Bot Token: " TELEGRAM_TOKEN
    done
    print_status "Token saved"
    echo ""

    # Get Police Department URL
    echo -e "${YELLOW}2. Police Department Configuration${NC}"
    echo "   Example: https://southmiamipdfl.policetocitizen.com"
    echo ""
    read -p "Enter Police Department Base URL: " BASE_URL
    while [ -z "$BASE_URL" ]; do
        print_error "Base URL cannot be empty!"
        read -p "Enter Police Department Base URL: " BASE_URL
    done
    print_status "Base URL saved"
    echo ""

    # Get Agency ID
    echo -e "${YELLOW}3. Agency ID${NC}"
    echo "   Find this in browser DevTools Network tab"
    echo "   Look for: /api/CADCalls/[NUMBER]"
    echo "   Example: 386 (South Miami PD)"
    echo ""
    read -p "Enter Agency ID: " AGENCY_ID
    while [ -z "$AGENCY_ID" ] || ! [[ "$AGENCY_ID" =~ ^[0-9]+$ ]]; do
        print_error "Agency ID must be a number!"
        read -p "Enter Agency ID: " AGENCY_ID
    done
    print_status "Agency ID saved"
    echo ""

    # Check interval
    echo -e "${YELLOW}4. Check Interval (optional)${NC}"
    read -p "Check interval in minutes [default: 5]: " CHECK_INTERVAL
    CHECK_INTERVAL=${CHECK_INTERVAL:-5}
    print_status "Check interval set to $CHECK_INTERVAL minutes"
    echo ""
}

# Create config.py
create_config() {
    print_info "Creating config.py..."

    cat > "$SCRIPT_DIR/config.py" <<EOF
"""
Configuration for CAD Calls API.

This file contains configuration parameters for accessing the Police-to-Citizen portal.
Generated by setup script.
"""

# Police to Citizen API Configuration
BASE_URL = "$BASE_URL"
AGENCY_ID = $AGENCY_ID

# API Endpoints - derived from base URL
API_ENDPOINTS = {
    "cadcalls": f"{BASE_URL}/api/CADCalls/{AGENCY_ID}"
}

# Request settings
API_SETTINGS = {
    "verify_ssl": False,
    "timeout": 30,
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "request_method": "POST"
}

# Default request parameters
DEFAULT_PARAMS = {
    "include_open": True,
    "include_closed": False,
    "take": 30,
    "skip": 0,
    "search_text": ""
}
EOF

    print_status "config.py created"
}

# Install Python dependencies
install_python_deps() {
    print_info "Installing Python dependencies..."

    cd "$SCRIPT_DIR"

    # Create virtual environment if desired
    read -p "Create Python virtual environment? (recommended) [Y/n]: " CREATE_VENV
    CREATE_VENV=${CREATE_VENV:-Y}

    if [[ "$CREATE_VENV" =~ ^[Yy]$ ]]; then
        if [ ! -d "venv" ]; then
            python3 -m venv venv
            print_status "Virtual environment created"
        fi
        source venv/bin/activate
        PYTHON_ACTIVATE="source venv/bin/activate &&"
    else
        PYTHON_ACTIVATE=""
    fi

    pip3 install -q --upgrade pip
    pip3 install -q -r requirements.txt

    print_status "Python dependencies installed"
}

# Build Go bot
build_bot() {
    print_info "Building Telegram bot..."

    cd "$SCRIPT_DIR/telegram-cad-bot"

    # Download Go dependencies
    go mod download

    # Build the bot
    go build -o cadbot cmd/bot/main.go

    print_status "Bot built successfully"
}

# Create start script
create_start_script() {
    print_info "Creating start script..."

    cat > "$SCRIPT_DIR/start-bot.sh" <<EOF
#!/bin/bash
# Start the Telegram CAD Calls Monitor Bot

cd "\$(dirname "\$0")"

export TELEGRAM_BOT_TOKEN="$TELEGRAM_TOKEN"
export CHECK_INTERVAL="$CHECK_INTERVAL"

# Activate Python virtual environment if it exists
if [ -d "venv" ]; then
    source venv/bin/activate
fi

# Start the bot
cd telegram-cad-bot
./cadbot -interval "$CHECK_INTERVAL"
EOF

    chmod +x "$SCRIPT_DIR/start-bot.sh"
    print_status "Start script created: start-bot.sh"
}

# Create systemd service
create_systemd_service() {
    read -p "Install as systemd service (auto-start on boot)? [Y/n]: " INSTALL_SERVICE
    INSTALL_SERVICE=${INSTALL_SERVICE:-Y}

    if [[ ! "$INSTALL_SERVICE" =~ ^[Yy]$ ]]; then
        return
    fi

    print_info "Creating systemd service..."

    SERVICE_NAME="telegram-cad-bot"
    SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

    sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=Telegram CAD Calls Monitor Bot
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$SCRIPT_DIR/telegram-cad-bot
Environment="TELEGRAM_BOT_TOKEN=$TELEGRAM_TOKEN"
Environment="CHECK_INTERVAL=$CHECK_INTERVAL"
ExecStart=$SCRIPT_DIR/telegram-cad-bot/cadbot -interval $CHECK_INTERVAL
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

    sudo systemctl daemon-reload
    sudo systemctl enable "$SERVICE_NAME"

    print_status "Systemd service created: $SERVICE_NAME"
    print_info "Start with: sudo systemctl start $SERVICE_NAME"
    print_info "View logs: sudo journalctl -u $SERVICE_NAME -f"
}

# Test the setup
test_setup() {
    echo ""
    echo -e "${BLUE}=================================================="
    echo "  Testing Setup"
    echo "==================================================${NC}"
    echo ""

    print_info "Testing Python script..."
    cd "$SCRIPT_DIR"

    if [ -d "venv" ]; then
        source venv/bin/activate
    fi

    python3 direct_api_post.py --take 5 --quiet

    if [ $? -eq 0 ]; then
        print_status "Python script works!"
    else
        print_error "Python script test failed. Check your configuration."
        exit 1
    fi
}

# Print final instructions
print_instructions() {
    echo ""
    echo -e "${GREEN}=================================================="
    echo "  ✓ Installation Complete!"
    echo "==================================================${NC}"
    echo ""
    echo -e "${BLUE}Quick Start:${NC}"
    echo ""
    echo "  1. Start the bot:"
    echo "     cd $SCRIPT_DIR"
    echo "     ./start-bot.sh"
    echo ""
    echo "  2. In Telegram, find your bot and send:"
    echo "     /start"
    echo "     /add Main Street"
    echo "     /check"
    echo ""

    if [ -f "/etc/systemd/system/telegram-cad-bot.service" ]; then
        echo -e "${BLUE}Systemd Service:${NC}"
        echo ""
        echo "  Start:   sudo systemctl start telegram-cad-bot"
        echo "  Stop:    sudo systemctl stop telegram-cad-bot"
        echo "  Status:  sudo systemctl status telegram-cad-bot"
        echo "  Logs:    sudo journalctl -u telegram-cad-bot -f"
        echo ""
    fi

    echo -e "${BLUE}Configuration Files:${NC}"
    echo ""
    echo "  config.py          - Police department settings"
    echo "  start-bot.sh       - Start script"
    echo "  telegram-cad-bot/  - Bot application"
    echo ""
    echo -e "${YELLOW}Need Help?${NC}"
    echo "  Check README.md for detailed documentation"
    echo ""
}

# Main installation flow
main() {
    check_root
    install_python
    install_go
    configure_bot
    create_config
    install_python_deps
    build_bot
    create_start_script
    test_setup

    if [ "$INSTALL_SYSTEMD" = true ] || [ "$EUID" -eq 0 ]; then
        create_systemd_service
    fi

    print_instructions
}

# Run main installation
main
