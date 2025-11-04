# Quick Start Guide

Get your Telegram CAD Monitor Bot up and running in 5 minutes!

## 1. Get a Telegram Bot Token

1. Open Telegram and search for `@BotFather`
2. Send `/newbot` command
3. Choose a name for your bot (e.g., "My CAD Monitor")
4. Choose a username (must end in 'bot', e.g., "mycadmonitor_bot")
5. Copy the bot token (looks like: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)

## 2. Find Your Police Department Info

### Option A: Use Default (South Miami PD)
The bot comes preconfigured with South Miami PD:
- Base URL: `https://southmiamipdfl.policetocitizen.com`
- Agency ID: `386`

### Option B: Configure Your Own Department
1. Go to your department's Police-to-Citizen portal (e.g., `https://yourtown.policetocitizen.com`)
2. Click on "CAD Calls" or "Active Calls"
3. Open browser Developer Tools (F12)
4. Go to Network tab
5. Look for API requests to `/api/CADCalls/[NUMBER]`
6. Note the base URL and the NUMBER (that's your Agency ID)

## 3. Run the Bot

### Method 1: Using Docker (Recommended)

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and add your bot token
nano .env  # or use any text editor

# Start the bot
docker-compose up -d

# View logs
docker-compose logs -f
```

### Method 2: Using the Binary

```bash
# If you haven't built it yet
go build -o cadbot cmd/bot/main.go

# Run with your token
./cadbot -token "YOUR_BOT_TOKEN_HERE"

# Or with custom settings
./cadbot \
  -token "YOUR_BOT_TOKEN_HERE" \
  -base-url "https://yourpd.policetocitizen.com" \
  -agency-id 123 \
  -interval 5
```

### Method 3: Using Environment Variable

```bash
# Set your token as an environment variable
export TELEGRAM_BOT_TOKEN="YOUR_BOT_TOKEN_HERE"

# Run the bot
./cadbot
```

## 4. Start Using the Bot

1. Open Telegram
2. Search for your bot by username (e.g., `@mycadmonitor_bot`)
3. Click "Start" or send `/start`
4. Add streets to monitor:
   ```
   /add Main Street
   /add Oak Avenue
   /add Maple Drive
   ```
5. Check your watch list:
   ```
   /list
   ```
6. Wait for alerts! The bot will automatically check for new calls every 5 minutes (or whatever interval you set)

## 5. Test It Out

Force an immediate check:
```
/check
```

View your settings:
```
/status
```

Change check interval to 10 minutes:
```
/interval 10
```

## Example Session

```
You:  /start
Bot:  👋 Welcome to CAD Calls Monitor Bot! ...

You:  /add Main Street
Bot:  ✅ Added 'Main Street' to your watch list!

You:  /check
Bot:  🔍 Checking for new calls...
Bot:  ✅ No new calls found on your monitored streets.

[Later...]
Bot:  🚨 New Call
      📍 Address: 123 Main Street
      🏷 Type: Emergency
      📋 Nature: Traffic Accident
      🕐 Time: 2025-11-04 14:23
      🆔 Incident: INC-2025-001234
```

## Troubleshooting

**Bot doesn't respond?**
- Check that the bot is running
- Verify your token is correct
- Make sure you clicked "Start" in Telegram

**No alerts coming?**
- Verify you added streets: `/list`
- Check if there are actually calls: `/check`
- Make sure the police department portal is accessible

**Build errors?**
```bash
# Make sure you have Go 1.19+
go version

# Install dependencies
go mod download

# Try building again
go build -o cadbot cmd/bot/main.go
```

**Docker errors?**
```bash
# Rebuild the image
docker-compose build --no-cache

# Check logs
docker-compose logs
```

## Next Steps

- ⚙️ Read the full [README.md](README.md) for detailed documentation
- 🔧 Configure custom check intervals with `/interval`
- 📊 Monitor multiple streets for comprehensive coverage
- 🛠️ Set up as a systemd service for production use

## Support

Need help? Check:
1. Full README.md for detailed docs
2. GitHub issues for common problems
3. Logs for error messages (`docker-compose logs` or check terminal output)

---

Enjoy your new CAD monitoring bot! 🎉
