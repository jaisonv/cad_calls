# Installation Guide

Complete installation guide for the Telegram CAD Calls Monitor Bot on a fresh Linux machine.

## 🚀 Quick Install (Recommended)

For a brand new Linux machine with nothing installed:

```bash
# 1. Clone or download this repository
git clone <your-repo-url>
cd cad_calls

# 2. Run the interactive installer
./setup.sh
```

The installer will:
- ✅ Check and install Python 3 (if needed)
- ✅ Check and install Go (if needed)
- ✅ Prompt for your configuration
- ✅ Set up the bot automatically
- ✅ Test the installation
- ✅ Optionally create a systemd service

### What You'll Be Asked

The installer will prompt you for:

1. **Telegram Bot Token**
   - Get from [@BotFather](https://t.me/botfather) on Telegram
   - Create a new bot with `/newbot`
   - Example: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`

2. **Police Department Base URL**
   - Your department's Police-to-Citizen portal
   - Example: `https://southmiamipdfl.policetocitizen.com`

3. **Agency ID**
   - Found in browser DevTools (F12) → Network tab
   - Look for `/api/CADCalls/[NUMBER]`
   - Example: `386`

4. **Check Interval** (optional)
   - How often to check for new calls (in minutes)
   - Default: `5` minutes

---

## 📋 Manual Installation

If you prefer to install manually or the automatic script doesn't work:

### Prerequisites

Install on Ubuntu/Debian:
```bash
sudo apt-get update
sudo apt-get install -y python3 python3-pip python3-venv wget git
```

Install on CentOS/RHEL/Fedora:
```bash
sudo yum install -y python3 python3-pip wget git
# or
sudo dnf install -y python3 python3-pip wget git
```

### Install Go

```bash
# Download Go 1.21.5
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# Install
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### Set Up Python Environment

```bash
cd cad_calls

# Create virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt
```

### Configure the Bot

Edit `config.py`:

```bash
nano config.py
```

Set your values:
```python
BASE_URL = "https://yourpd.policetocitizen.com"
AGENCY_ID = 123
```

### Build the Go Bot

```bash
cd telegram-cad-bot

# Download dependencies
go mod download

# Build
go build -o cadbot cmd/bot/main.go
```

### Create Start Script

Create `start-bot.sh`:

```bash
#!/bin/bash
export TELEGRAM_BOT_TOKEN="your_token_here"
cd telegram-cad-bot
./cadbot -interval 5
```

Make it executable:
```bash
chmod +x start-bot.sh
```

---

## 🔧 Running the Bot

### Option 1: Run Directly

```bash
cd cad_calls
./start-bot.sh
```

### Option 2: Systemd Service (Auto-start on boot)

Create `/etc/systemd/system/telegram-cad-bot.service`:

```ini
[Unit]
Description=Telegram CAD Calls Monitor Bot
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/path/to/cad_calls/telegram-cad-bot
Environment="TELEGRAM_BOT_TOKEN=your_token_here"
ExecStart=/path/to/cad_calls/telegram-cad-bot/cadbot -interval 5
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable telegram-cad-bot
sudo systemctl start telegram-cad-bot

# Check status
sudo systemctl status telegram-cad-bot

# View logs
sudo journalctl -u telegram-cad-bot -f
```

### Option 3: Screen/tmux (Background session)

Using screen:
```bash
screen -S cadbot
cd cad_calls
./start-bot.sh

# Detach: Ctrl+A, then D
# Reattach: screen -r cadbot
```

Using tmux:
```bash
tmux new -s cadbot
cd cad_calls
./start-bot.sh

# Detach: Ctrl+B, then D
# Reattach: tmux attach -t cadbot
```

---

## 📱 Testing the Bot

1. **Open Telegram** and find your bot by username

2. **Send commands:**
   ```
   /start
   /add Main Street
   /add Oak Avenue
   /list
   /check
   ```

3. **Wait for alerts** - The bot checks every 5 minutes (or your configured interval)

---

## 🔍 Finding Your Police Department Info

### Find Base URL

Your police department's portal should look like:
- `https://[city][state].policetocitizen.com`
- `https://[county]sheriff.policetocitizen.com`

Examples:
- Tyler, TX: `https://tylertx.policetocitizen.com`
- South Miami, FL: `https://southmiamipdfl.policetocitizen.com`
- Wood County, OH: `https://woodcountyohsheriff.policetocitizen.com`

### Find Agency ID

1. Go to your department's portal in a browser
2. Click "CAD Calls" or "Active Calls"
3. Press F12 to open Developer Tools
4. Go to "Network" tab
5. Refresh the page
6. Look for a request to `/api/CADCalls/[NUMBER]`
7. The NUMBER is your Agency ID

**Example:**
```
Request URL: https://southmiamipdfl.policetocitizen.com/api/CADCalls/386
                                                                    ^^^
                                                              Agency ID
```

---

## 🐛 Troubleshooting

### Bot doesn't start

**Check logs:**
```bash
# If using systemd
sudo journalctl -u telegram-cad-bot -n 50

# If running manually
./start-bot.sh
# Look for error messages
```

**Common issues:**
- Token incorrect → Check with @BotFather
- Python script fails → Test with: `python3 direct_api_post.py --take 5`
- Go build fails → Ensure Go 1.19+ is installed: `go version`

### No alerts received

1. **Check if streets are added:**
   ```
   /list
   ```

2. **Manually trigger check:**
   ```
   /check
   ```

3. **Check if there are actual calls:**
   ```bash
   python3 direct_api_post.py --take 10
   # Look in cadcalls_results/ directory
   ```

### Python script fails

**Test your configuration:**
```bash
python3 direct_api_post.py --take 5
```

**Check error messages:**
- 400 Bad Request → Agency ID might be wrong
- 404 Not Found → Base URL or Agency ID incorrect
- 403 Forbidden → Department blocks automated requests
- Connection errors → Check BASE_URL

### Permission errors

**If systemd service fails:**
```bash
# Check service status
sudo systemctl status telegram-cad-bot

# Fix permissions
sudo chown -R yourusername:yourusername /path/to/cad_calls

# Restart service
sudo systemctl restart telegram-cad-bot
```

---

## 📦 System Requirements

### Minimum Requirements
- **OS:** Linux (Ubuntu 18.04+, Debian 10+, CentOS 7+, etc.)
- **RAM:** 512 MB
- **Disk:** 500 MB free space
- **Network:** Internet connection

### Software Requirements (installed by setup.sh)
- Python 3.6+
- Go 1.19+
- pip
- wget (for Go installation)

---

## 🔄 Updating

To update the bot:

```bash
cd cad_calls

# Pull latest changes
git pull

# Rebuild the bot
cd telegram-cad-bot
go build -o cadbot cmd/bot/main.go

# Restart
sudo systemctl restart telegram-cad-bot
# or
./start-bot.sh
```

---

## 🗑️ Uninstallation

```bash
# Stop the service
sudo systemctl stop telegram-cad-bot
sudo systemctl disable telegram-cad-bot

# Remove service file
sudo rm /etc/systemd/system/telegram-cad-bot.service
sudo systemctl daemon-reload

# Remove files
cd ~
rm -rf cad_calls

# Optional: Remove Go (if you don't need it)
sudo rm -rf /usr/local/go
```

---

## 📞 Support

- **Documentation:** See [README.md](README.md)
- **Troubleshooting:** See [TROUBLESHOOTING.md](telegram-cad-bot/TROUBLESHOOTING.md)
- **Issues:** Open a GitHub issue

---

## 🔒 Security Notes

- Keep your `TELEGRAM_BOT_TOKEN` secret
- Don't commit `config.py` with real credentials
- Use a dedicated user for the systemd service
- Consider using SSL verification (`verify_ssl: true`) in production

---

**Happy monitoring! 🚨📍**
