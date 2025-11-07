# Quick Install - 5 Minutes ⚡

Get the Telegram CAD Monitor Bot running in 5 minutes on a fresh Linux machine.

## One-Command Install

```bash
git clone <your-repo-url>
cd cad_calls
./setup.sh
```

That's it! The installer will:
1. ✅ Install Python 3 & Go (if needed)
2. ✅ Ask for your configuration
3. ✅ Build everything
4. ✅ Test it works

## What You'll Need Ready

Before running `./setup.sh`, have these ready:

### 1️⃣ Telegram Bot Token

1. Open Telegram → Search for `@BotFather`
2. Send: `/newbot`
3. Follow prompts to create your bot
4. **Copy the token** (looks like `1234567890:ABCdef...`)

### 2️⃣ Police Department Info

Find your department's **Base URL**:
```
Format: https://[city][state].policetocitizen.com
Example: https://southmiamipdfl.policetocitizen.com
```

Find the **Agency ID**:
1. Go to your department's portal
2. Press **F12** → **Network** tab
3. Click "CAD Calls"
4. Look for: `/api/CADCalls/[NUMBER]`
5. The **NUMBER** is your Agency ID

Example:
```
https://southmiamipdfl.policetocitizen.com/api/CADCalls/386
                                                        ^^^
                                                    Agency ID
```

## Run the Installer

```bash
./setup.sh
```

You'll be prompted for:
- Telegram Bot Token
- Police Department Base URL
- Agency ID
- Check interval (default: 5 minutes)

## Start the Bot

After installation completes:

```bash
./start-bot.sh
```

## Test in Telegram

1. Find your bot in Telegram
2. Send:
   ```
   /start
   /add Main Street
   /check
   ```

Done! You'll get alerts when calls happen on streets you're monitoring.

---

## Auto-Start on Boot (Optional)

The installer asks if you want to install as a systemd service.

If yes:
```bash
sudo systemctl start telegram-cad-bot
sudo systemctl status telegram-cad-bot
```

View logs:
```bash
sudo journalctl -u telegram-cad-bot -f
```

---

## Troubleshooting

**Problem:** Setup script fails

```bash
# Check what's missing
./verify-install.sh
```

**Problem:** Bot doesn't respond in Telegram

```bash
# Test the Python script
python3 direct_api_post.py --take 5

# Check if bot is running
ps aux | grep cadbot
```

**Problem:** No alerts

```
/list     - Check your watched streets
/check    - Manually trigger a check
/status   - Check bot status
```

---

## Full Documentation

- **Detailed Install:** [INSTALL.md](INSTALL.md)
- **Usage Guide:** [telegram-cad-bot/README.md](telegram-cad-bot/README.md)
- **Troubleshooting:** [telegram-cad-bot/TROUBLESHOOTING.md](telegram-cad-bot/TROUBLESHOOTING.md)

---

**That's it! Your bot is now monitoring CAD calls! 🚨**
