# Telegram CAD Calls Monitor Bot

A Telegram bot written in Go that monitors active CAD (Computer Aided Dispatch) calls from police departments and sends real-time alerts for incidents on streets you're watching.

## Features

- 🔍 **Street Monitoring**: Add specific streets to your watch list
- 🔔 **Real-time Alerts**: Get instant notifications for new calls on monitored streets
- 📊 **Smart Filtering**: Case-insensitive matching with street abbreviation support
- 💾 **Persistent Storage**: SQLite database for user preferences and seen calls
- ⏱️ **Configurable Intervals**: Set custom check intervals (1-60 minutes)
- 🚀 **Concurrent Processing**: Efficient polling for multiple users
- 🧹 **Automatic Cleanup**: Removes old seen calls to keep database lean

## Prerequisites

- Go 1.19 or higher
- Telegram Bot Token (get one from [@BotFather](https://t.me/botfather))
- Access to a Police-to-Citizen CAD portal

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
cd telegram-cad-bot

# Download dependencies
go mod download

# Build the binary
go build -o cadbot cmd/bot/main.go

# Run the bot
./cadbot -token YOUR_BOT_TOKEN
```

### Option 2: Using Docker

```bash
# Build the Docker image
docker build -t telegram-cad-bot .

# Run the container
docker run -d \
  -e TELEGRAM_BOT_TOKEN=your_token_here \
  -v $(pwd)/data:/data \
  --name cadbot \
  telegram-cad-bot
```

## Configuration

The bot can be configured using command-line flags or environment variables:

### Command-Line Flags

```bash
./cadbot \
  -token YOUR_BOT_TOKEN \
  -base-url https://yourpd.policetocitizen.com \
  -agency-id 386 \
  -db ./cadbot.db \
  -interval 5 \
  -verify-ssl false
```

### Environment Variables

```bash
export TELEGRAM_BOT_TOKEN="your_token_here"
export CAD_BASE_URL="https://yourpd.policetocitizen.com"
export CAD_AGENCY_ID="386"
export CHECK_INTERVAL="5"
```

### Configuration Options

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `-token` | `TELEGRAM_BOT_TOKEN` | - | Telegram bot token (required) |
| `-base-url` | `CAD_BASE_URL` | `https://southmiamipdfl.policetocitizen.com` | CAD API base URL |
| `-agency-id` | `CAD_AGENCY_ID` | `386` | Agency ID for the police department |
| `-db` | `DB_PATH` | `./cadbot.db` | SQLite database file path |
| `-interval` | `CHECK_INTERVAL` | `5` | Check interval in minutes |
| `-verify-ssl` | `VERIFY_SSL` | `false` | Verify SSL certificates |

## Getting Your Bot Token

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the bot token provided by BotFather
5. Use the token with the `-token` flag or `TELEGRAM_BOT_TOKEN` environment variable

## Finding Your Police Department

To use this bot with your local police department:

1. Search for your department at: `https://yourtown.policetocitizen.com`
2. Navigate to the CAD Calls section
3. Open browser developer tools (F12)
4. Look for API requests to `/api/CADCalls/[number]`
5. The number is your Agency ID

**Note**: Not all Police-to-Citizen portals support automated access. Some may have security measures that block bots.

## Bot Commands

Once the bot is running, interact with it using these commands:

| Command | Description | Example |
|---------|-------------|---------|
| `/start` | Welcome message and instructions | `/start` |
| `/help` | Show help message | `/help` |
| `/add <street>` | Add a street to watch list | `/add Main Street` |
| `/remove <street>` | Remove a street from watch list | `/remove Main Street` |
| `/list` | Show all monitored streets | `/list` |
| `/clear` | Clear all monitored streets | `/clear` |
| `/status` | Show bot status and settings | `/status` |
| `/interval <min>` | Set check interval (1-60 minutes) | `/interval 10` |
| `/check` | Manually check for new calls | `/check` |

## Usage Example

```
You: /start
Bot: 👋 Welcome to CAD Calls Monitor Bot! ...

You: /add Main Street
Bot: ✅ Added 'Main Street' to your watch list!

You: /add Oak Avenue
Bot: ✅ Added 'Oak Avenue' to your watch list!

You: /list
Bot: 📋 Your Watch List (2 streets):
     1. Main Street
     2. Oak Avenue

You: /status
Bot: 📊 Bot Status
     👤 User ID: 123456789
     📍 Monitored Streets: 2
     ⏱ Check Interval: 5 minutes
     ...

[Later, when a new call appears]
Bot: 🚨 New Call
     📍 Address: 123 Main Street
     🏷 Type: Emergency
     📋 Nature: Traffic Accident
     🕐 Time: 2025-11-04 14:23
     🆔 Incident: INC-2025-001234
```

## Street Matching

The bot uses intelligent street matching:

- **Case-insensitive**: "main street" matches "Main Street"
- **Partial matching**: "Main St" matches "123 Main Street"
- **Abbreviation support**: Automatically handles common abbreviations
  - Street/St/Str
  - Avenue/Ave/Av
  - Road/Rd
  - Drive/Dr
  - Boulevard/Blvd
  - Lane/Ln
  - Court/Ct

## Architecture

```
telegram-cad-bot/
├── cmd/
│   └── bot/
│       └── main.go              # Application entry point
├── internal/
│   ├── bot/
│   │   ├── bot.go               # Bot initialization
│   │   └── handlers.go          # Command handlers
│   ├── cad/
│   │   ├── client.go            # CAD API client
│   │   └── models.go            # Data structures
│   ├── storage/
│   │   ├── db.go                # Database operations
│   │   └── models.go            # DB models
│   ├── filter/
│   │   └── street.go            # Street matching logic
│   └── scheduler/
│       └── poller.go            # Periodic polling
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Database Schema

The bot uses SQLite with the following schema:

```sql
-- Users table
CREATE TABLE users (
    telegram_id INTEGER PRIMARY KEY,
    check_interval INTEGER DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Monitored streets
CREATE TABLE monitored_streets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id INTEGER NOT NULL,
    street_name TEXT NOT NULL,
    UNIQUE(telegram_id, street_name)
);

-- Seen calls (to prevent duplicate alerts)
CREATE TABLE seen_calls (
    telegram_id INTEGER NOT NULL,
    incident_id TEXT NOT NULL,
    seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(telegram_id, incident_id)
);
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o cadbot-linux cmd/bot/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o cadbot-macos cmd/bot/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o cadbot.exe cmd/bot/main.go
```

## Deployment

### Using systemd (Linux)

Create `/etc/systemd/system/cadbot.service`:

```ini
[Unit]
Description=Telegram CAD Calls Monitor Bot
After=network.target

[Service]
Type=simple
User=cadbot
WorkingDirectory=/opt/cadbot
Environment="TELEGRAM_BOT_TOKEN=your_token_here"
ExecStart=/opt/cadbot/cadbot
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable cadbot
sudo systemctl start cadbot
sudo systemctl status cadbot
```

### Using Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  cadbot:
    build: .
    container_name: telegram-cad-bot
    restart: unless-stopped
    environment:
      - TELEGRAM_BOT_TOKEN=your_token_here
      - CAD_BASE_URL=https://yourpd.policetocitizen.com
      - CAD_AGENCY_ID=386
      - CHECK_INTERVAL=5
    volumes:
      - ./data:/data
```

Run:

```bash
docker-compose up -d
```

## Troubleshooting

### Bot doesn't respond

- Verify your bot token is correct
- Check that the bot is running: `systemctl status cadbot` (if using systemd)
- Check logs for errors

### No alerts received

- Verify you have streets in your watch list: `/list`
- Check that there are active calls on your monitored streets
- Manually trigger a check: `/check`
- Verify the CAD API is accessible

### API connection errors

- Check that the base URL and agency ID are correct
- Some departments may block automated requests
- Try increasing the timeout value
- Check if SSL verification is needed (`-verify-ssl true`)

## Legal and Ethical Considerations

⚠️ **Important**: This bot accesses publicly available data from Police-to-Citizen portals. Please use responsibly:

- ✅ Use for personal safety awareness
- ✅ Use for research and journalism
- ✅ Use for community monitoring
- ❌ Do NOT use for harassment or stalking
- ❌ Do NOT make excessive requests
- ❌ Do NOT violate terms of service

The data represents real emergency situations and should be treated with sensitivity.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is provided for educational and informational purposes. Users are responsible for ensuring compliance with applicable laws and terms of service.

## Future Features (Not in MVP)

- 🗺️ Geographic perimeter filtering
- 📈 Statistics and analytics
- 🔕 Custom notification schedules
- 🌐 Multi-agency support
- 📱 Location-based monitoring
- 🤖 Natural language queries

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing issues for solutions
- Review the troubleshooting section

---

Built with ❤️ using Go and Telebot
