# Running with Systemd (Production Deployment)

Complete guide for running the Telegram CAD Monitor Bot as a systemd service on Linux.

## Why Use Systemd?

- ✅ **Auto-start on boot** - Bot starts automatically when server reboots
- ✅ **Auto-restart on crash** - Service restarts if bot crashes
- ✅ **Background operation** - Runs in background, no terminal needed
- ✅ **Centralized logging** - Logs managed by journald
- ✅ **Resource management** - System monitors resource usage

---

## Quick Setup (Using Installer)

The easiest way is to use the automated installer:

```bash
cd cad_calls
./setup.sh
```

When prompted "Install as systemd service?", answer **Y**

The installer will:
1. Create the service file automatically
2. Use the current username for the service
3. Enable auto-start on boot
4. Start the service

**Then skip to [Managing the Service](#managing-the-service) section below.**

---

## Manual Setup

If you prefer to set up manually or need to customize:

### 1. Create the Service File

```bash
sudo nano /etc/systemd/system/telegram-cad-bot.service
```

Paste this configuration (adjust paths and username):

```ini
[Unit]
Description=Telegram CAD Calls Monitor Bot
Documentation=https://github.com/jaisonv/cad_calls
After=network-online.target
Wants=network-online.target

[Service]
Type=simple

# IMPORTANT: Change 'yourusername' to your actual Linux username
User=yourusername
Group=yourusername

# IMPORTANT: Update paths to match your installation
WorkingDirectory=/home/yourusername/cad_calls/telegram-cad-bot
ExecStart=/home/yourusername/cad_calls/telegram-cad-bot/cadbot -interval 5

# Environment variables
Environment="TELEGRAM_BOT_TOKEN=your_bot_token_here"
Environment="CHECK_INTERVAL=5"

# Restart policy
Restart=always
RestartSec=10
StartLimitInterval=0

# Security settings (optional but recommended)
NoNewPrivileges=true
PrivateTmp=true

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=cadbot

[Install]
WantedBy=multi-user.target
```

**Important:** Replace these values:
- `yourusername` → Your actual Linux username (e.g., `jaisonv`, `pi`)
- `/home/yourusername/cad_calls` → Actual path to your installation
- `your_bot_token_here` → Your Telegram bot token from @BotFather

### 2. Verify Paths Are Correct

```bash
# Check the bot binary exists
ls -la /home/yourusername/cad_calls/telegram-cad-bot/cadbot

# Verify ownership matches the User= in service file
stat /home/yourusername/cad_calls/telegram-cad-bot/cadbot
```

### 3. Reload Systemd

```bash
sudo systemctl daemon-reload
```

### 4. Enable Auto-Start on Boot

```bash
sudo systemctl enable telegram-cad-bot
```

This creates a symlink so the service starts automatically on boot.

### 5. Start the Service

```bash
sudo systemctl start telegram-cad-bot
```

---

## Managing the Service

### Check Service Status

```bash
systemctl status telegram-cad-bot
```

**Output examples:**

✅ **Good - Running:**
```
● telegram-cad-bot.service - Telegram CAD Calls Monitor Bot
   Loaded: loaded (/etc/systemd/system/telegram-cad-bot.service; enabled)
   Active: active (running) since Sat 2025-11-09 10:30:00 EST; 2h ago
 Main PID: 12345 (cadbot)
   CGroup: /system.slice/telegram-cad-bot.service
           └─12345 /home/jaisonv/cad_calls/telegram-cad-bot/cadbot -interval 5
```

❌ **Problem - Failed:**
```
● telegram-cad-bot.service - Telegram CAD Calls Monitor Bot
   Loaded: loaded (/etc/systemd/system/telegram-cad-bot.service; enabled)
   Active: failed (Result: exit-code) since Sat 2025-11-09 10:35:00 EST
```

### Start the Service

```bash
sudo systemctl start telegram-cad-bot
```

### Stop the Service

```bash
sudo systemctl stop telegram-cad-bot
```

### Restart the Service

```bash
sudo systemctl restart telegram-cad-bot
```

Use this after:
- Updating the bot code (`git pull` + rebuild)
- Changing configuration
- Editing the service file

### Disable Auto-Start on Boot

```bash
sudo systemctl disable telegram-cad-bot
```

The service won't start automatically anymore, but you can still start it manually.

### Re-enable Auto-Start

```bash
sudo systemctl enable telegram-cad-bot
```

---

## Monitoring Logs

### View Live Logs (Follow Mode)

```bash
sudo journalctl -u telegram-cad-bot -f
```

Press **Ctrl+C** to stop following.

**You'll see:**
- Bot startup messages
- When users add/remove streets
- When new calls are checked
- Errors and warnings

**Example output:**
```
Nov 09 10:30:00 raspberrypi cadbot[12345]: [CADBot] 2025/11/09 10:30:00 Initializing bot...
Nov 09 10:30:00 raspberrypi cadbot[12345]: [CADBot] 2025/11/09 10:30:00 Python script: /home/jaisonv/cad_calls/direct_api_post.py
Nov 09 10:30:00 raspberrypi cadbot[12345]: [CADBot] 2025/11/09 10:30:00 Poller started
Nov 09 10:30:00 raspberrypi cadbot[12345]: [CADBot] 2025/11/09 10:30:00 Bot is running. Press Ctrl+C to stop.
Nov 09 10:35:00 raspberrypi cadbot[12345]: [CADBot] 2025/11/09 10:35:00 Checking calls for 2 active users
Nov 09 10:35:01 raspberrypi cadbot[12345]: [CADBot] 2025/11/09 10:35:01 Found 1 new calls for user 121107313
```

### View Last N Lines

```bash
# Last 50 lines
sudo journalctl -u telegram-cad-bot -n 50

# Last 100 lines
sudo journalctl -u telegram-cad-bot -n 100
```

### View Logs Since a Time

```bash
# Since today
sudo journalctl -u telegram-cad-bot --since today

# Since 1 hour ago
sudo journalctl -u telegram-cad-bot --since "1 hour ago"

# Since specific date
sudo journalctl -u telegram-cad-bot --since "2025-11-09 08:00:00"
```

### View Logs in Specific Time Range

```bash
sudo journalctl -u telegram-cad-bot --since "2025-11-09 08:00" --until "2025-11-09 12:00"
```

### Search Logs for Errors

```bash
# Show only errors
sudo journalctl -u telegram-cad-bot -p err

# Search for specific text
sudo journalctl -u telegram-cad-bot | grep "error"
sudo journalctl -u telegram-cad-bot | grep "Permission denied"
```

### Export Logs to File

```bash
# Export all logs
sudo journalctl -u telegram-cad-bot > bot-logs.txt

# Export last 100 lines
sudo journalctl -u telegram-cad-bot -n 100 > bot-logs-recent.txt
```

### Clear Old Logs (Optional)

Logs can grow large over time. To clean up:

```bash
# Keep only last 7 days
sudo journalctl --vacuum-time=7d

# Keep only last 100MB
sudo journalctl --vacuum-size=100M
```

---

## Updating the Bot

When you pull new code or make changes:

```bash
cd ~/cad_calls

# 1. Pull latest code
git pull

# 2. Rebuild the bot
cd telegram-cad-bot
go build -o cadbot cmd/bot/main.go

# 3. Restart the service
sudo systemctl restart telegram-cad-bot

# 4. Verify it's running
systemctl status telegram-cad-bot

# 5. Check logs for errors
sudo journalctl -u telegram-cad-bot -n 20
```

---

## Troubleshooting

### Service Won't Start

**Check the status:**
```bash
systemctl status telegram-cad-bot
```

**Look for:**
- `Failed to start` - Check logs for why
- `Unit not found` - Service file doesn't exist
- `Permission denied` - User/path incorrect

**Common fixes:**

1. **Binary doesn't exist:**
   ```bash
   cd ~/cad_calls/telegram-cad-bot
   go build -o cadbot cmd/bot/main.go
   ```

2. **Wrong user in service file:**
   ```bash
   sudo nano /etc/systemd/system/telegram-cad-bot.service
   # Change User= to your actual username
   sudo systemctl daemon-reload
   sudo systemctl restart telegram-cad-bot
   ```

3. **Wrong paths:**
   ```bash
   # Check paths in service file match reality
   cat /etc/systemd/system/telegram-cad-bot.service | grep -E "WorkingDirectory|ExecStart"
   ```

### Service Keeps Restarting

```bash
# Check logs for errors
sudo journalctl -u telegram-cad-bot -n 50

# Common issues:
# - Missing TELEGRAM_BOT_TOKEN
# - Python script errors (permission, missing packages)
# - config.py not configured
```

### Permission Denied Errors

```bash
# Fix ownership of all files
sudo chown -R yourusername:yourusername ~/cad_calls

# Fix permissions
chmod -R 755 ~/cad_calls
chmod 644 ~/cad_calls/config.py

# Create cadcalls_results directory if missing
mkdir -p ~/cad_calls/cadcalls_results
```

### Can't Find Python Packages

The bot should auto-detect the venv. If it doesn't:

```bash
# Check if venv exists
ls ~/cad_calls/venv/bin/python3

# If not, create it
cd ~/cad_calls
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Rebuild bot (it will now find venv)
cd telegram-cad-bot
go build -o cadbot cmd/bot/main.go
sudo systemctl restart telegram-cad-bot
```

---

## Service File Reference

### Full Example with Comments

```ini
[Unit]
# Service description
Description=Telegram CAD Calls Monitor Bot
Documentation=https://github.com/jaisonv/cad_calls

# Wait for network before starting
After=network-online.target
Wants=network-online.target

[Service]
# Run as a simple service (doesn't fork)
Type=simple

# User and group to run as (CHANGE THIS)
User=jaisonv
Group=jaisonv

# Where to run from (CHANGE THIS)
WorkingDirectory=/home/jaisonv/cad_calls/telegram-cad-bot

# Command to execute (CHANGE THIS)
ExecStart=/home/jaisonv/cad_calls/telegram-cad-bot/cadbot -interval 5

# Environment variables
Environment="TELEGRAM_BOT_TOKEN=1234567890:ABCdef..."
Environment="CHECK_INTERVAL=5"

# Restart if it crashes
Restart=always
RestartSec=10
StartLimitInterval=0

# Security hardening (optional)
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=/home/jaisonv/cad_calls/cadcalls_results
ReadWritePaths=/home/jaisonv/cad_calls/telegram-cad-bot

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=cadbot

[Install]
# Start on boot (multi-user target = normal boot)
WantedBy=multi-user.target
```

### Environment Variables

You can add more environment variables:

```ini
Environment="TELEGRAM_BOT_TOKEN=your_token"
Environment="CHECK_INTERVAL=5"
Environment="PYTHON_SCRIPT=/path/to/custom/script.py"
```

Or use an environment file:

```ini
EnvironmentFile=/home/jaisonv/cad_calls/.env
```

Then create `/home/jaisonv/cad_calls/.env`:
```bash
TELEGRAM_BOT_TOKEN=your_token
CHECK_INTERVAL=5
```

---

## System Resource Monitoring

### Check Resource Usage

```bash
# CPU and memory usage
systemctl status telegram-cad-bot

# Detailed resource info
systemd-cgtop | grep cadbot
```

### Set Resource Limits (Optional)

Add to service file under `[Service]`:

```ini
# Limit memory to 256MB
MemoryMax=256M
MemoryHigh=200M

# Limit CPU to 50%
CPUQuota=50%
```

Then reload:
```bash
sudo systemctl daemon-reload
sudo systemctl restart telegram-cad-bot
```

---

## Multiple Bots

To run multiple bot instances (different police departments):

1. **Copy the project:**
   ```bash
   cp -r cad_calls cad_calls_tyler
   cp -r cad_calls cad_calls_miami
   ```

2. **Configure each separately:**
   ```bash
   nano cad_calls_tyler/config.py
   nano cad_calls_miami/config.py
   ```

3. **Create separate service files:**
   ```bash
   sudo cp /etc/systemd/system/telegram-cad-bot.service \
           /etc/systemd/system/telegram-cad-bot-tyler.service

   sudo cp /etc/systemd/system/telegram-cad-bot.service \
           /etc/systemd/system/telegram-cad-bot-miami.service
   ```

4. **Edit each service file** with different paths and tokens

5. **Enable and start:**
   ```bash
   sudo systemctl enable telegram-cad-bot-tyler
   sudo systemctl enable telegram-cad-bot-miami
   sudo systemctl start telegram-cad-bot-tyler
   sudo systemctl start telegram-cad-bot-miami
   ```

---

## Quick Reference

```bash
# Start/Stop/Restart
sudo systemctl start telegram-cad-bot
sudo systemctl stop telegram-cad-bot
sudo systemctl restart telegram-cad-bot

# Enable/Disable auto-start
sudo systemctl enable telegram-cad-bot
sudo systemctl disable telegram-cad-bot

# Status and logs
systemctl status telegram-cad-bot
sudo journalctl -u telegram-cad-bot -f
sudo journalctl -u telegram-cad-bot -n 50

# After editing service file
sudo systemctl daemon-reload
sudo systemctl restart telegram-cad-bot

# After updating code
cd ~/cad_calls/telegram-cad-bot
go build -o cadbot cmd/bot/main.go
sudo systemctl restart telegram-cad-bot
```

---

## See Also

- [INSTALL.md](INSTALL.md) - Complete installation guide
- [README.md](README.md) - General documentation
- [TROUBLESHOOTING.md](telegram-cad-bot/TROUBLESHOOTING.md) - Troubleshooting guide
