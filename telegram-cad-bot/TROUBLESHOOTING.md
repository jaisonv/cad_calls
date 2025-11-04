# Troubleshooting Guide

## Error: "API returned status 400"

This means the CAD API rejected the request. Here's how to fix it:

### Step 1: Test Your API Configuration

Use the test utility to verify your department's API works:

```bash
# Test with default (South Miami PD)
./test-api

# Test with your department
./test-api -base-url "https://yourpd.policetocitizen.com" -agency-id 123
```

**If the test succeeds**, your configuration is correct! Run the bot with those same settings:

```bash
./cadbot -token "YOUR_TOKEN" -base-url "https://yourpd.policetocitizen.com" -agency-id 123
```

**If the test fails**, see the error messages below.

### Step 2: Find Your Police Department Settings

#### Finding the Base URL

Your base URL should look like:
- `https://yourtown.policetocitizen.com`
- `https://yourcounty.policetocitizen.com`
- `https://departmentname.policetocitizen.com`

Common examples:
- Tyler PD: `https://tylertx.policetocitizen.com`
- South Miami: `https://southmiamipdfl.policetocitizen.com`
- Wood County: `https://woodcountyohsheriff.policetocitizen.com`

#### Finding the Agency ID

1. Open your department's portal in a browser
2. Click on "CAD Calls" or "Active Calls"
3. Open Developer Tools (F12)
4. Go to "Network" tab
5. Refresh the page
6. Look for a request to `/api/CADCalls/[NUMBER]`
7. That NUMBER is your Agency ID

**Example:**
```
Request URL: https://southmiamipdfl.policetocitizen.com/api/CADCalls/386
                                                                    ^^^
                                                              Agency ID = 386
```

### Step 3: Common Issues

#### Issue: "failed to visit main page"
**Cause**: Base URL is wrong or site is down
**Fix**:
- Verify the URL works in your browser
- Try without `https://` or with `http://` instead
- Check if there's a typo

#### Issue: "API returned status 404"
**Cause**: Wrong agency ID or API endpoint doesn't exist
**Fix**:
- Double-check the agency ID from browser dev tools
- Make sure you're looking at the POST request to `/api/CADCalls/[ID]`
- Try the agency ID ±1 (sometimes they're off by one)

#### Issue: "API returned status 400"
**Cause**: API request payload is wrong or department has different API structure
**Fix**:
- The department might use a different API version
- Check the Network tab to see what payload the website sends
- Some departments may not support automated access

#### Issue: "API returned status 403" or "Request Rejected"
**Cause**: Department blocks automated requests (WAF/firewall)
**Fix**:
- Some departments have security measures against bots
- Try setting `-verify-ssl true`
- Unfortunately, some departments can't be automated

### Step 4: Test with Python Script First

The Python script in the parent directory uses the same API. Test it first:

```bash
cd ..
# Edit config.py with your settings
nano config.py

# Run the Python script
python3 direct_api_post.py --take 10
```

If the Python script works, use the **exact same** `BASE_URL` and `AGENCY_ID` with the bot.

## Error: Bot doesn't respond in Telegram

### Check 1: Is the bot running?
```bash
ps aux | grep cadbot
```

### Check 2: Is the token correct?
- Go to @BotFather in Telegram
- Send `/mybots`
- Select your bot
- Check the token

### Check 3: Can the bot reach Telegram?
```bash
# Check logs for connection errors
# You should see "Bot is running" message
```

## Error: No alerts received

### Check 1: Do you have streets in your watch list?
```
/list
```

### Check 2: Are there actually calls on those streets?
```
/check
```
This manually triggers a check and will tell you if calls were found.

### Check 3: Are the calls new or already seen?
The bot only alerts for NEW calls. If you restart the bot, the database remembers which calls you've already seen.

To reset:
```bash
rm cadbot.db
# Restart bot
```

## Still Having Issues?

### Enable Debug Mode

Add verbose logging to see what's happening:

```bash
# Run with verbose output
./cadbot -token "TOKEN" 2>&1 | tee bot.log
```

### Check the Database

```bash
sqlite3 cadbot.db "SELECT * FROM users;"
sqlite3 cadbot.db "SELECT * FROM monitored_streets;"
sqlite3 cadbot.db "SELECT COUNT(*) FROM seen_calls;"
```

### Test API Directly

Use curl to test the API:

```bash
curl -X POST "https://yourpd.policetocitizen.com/api/CADCalls/386" \
  -H "Content-Type: application/json" \
  -d '{
    "IncludeOpenCalls": true,
    "IncludeClosedCalls": false,
    "IncludeCount": true,
    "PagingOptions": {
      "SortOptions": [{"Name": "StartTime", "SortDirection": "Descending", "Sequence": 1}],
      "Take": 10,
      "Skip": 0
    },
    "FilterOptionsParameters": {
      "IntersectionSearch": true,
      "SearchText": "",
      "Parameters": []
    }
  }'
```

## Known Working Departments

These departments are known to work (as of 2025):

| Department | Base URL | Agency ID | Notes |
|------------|----------|-----------|-------|
| Tyler PD | `https://tylertx.policetocitizen.com` | Check browser | ✅ Working |
| South Miami PD | `https://southmiamipdfl.policetocitizen.com` | 386 | ✅ Working |
| Wood County Sheriff | `https://woodcountyohsheriff.policetocitizen.com` | Check browser | ✅ Working |

## Known Non-Working Departments

These departments have security measures that block automated access:

- Lynchburg PD (WAF/firewall blocks bots)
- Departments with CAPTCHA enabled
- Departments requiring authentication

If your department blocks automated access, you'll need to use browser automation (Selenium) instead of this bot.

## Getting Help

1. Run the test utility: `./test-api -base-url "URL" -agency-id ID`
2. Check logs for detailed error messages
3. Verify the Python script works with same settings
4. Open an issue on GitHub with:
   - Error message
   - Department URL (if comfortable sharing)
   - Output from test utility

---

**Quick Test Checklist:**
- [ ] Base URL works in browser
- [ ] Found Agency ID in browser dev tools
- [ ] `./test-api` succeeds with your settings
- [ ] Bot token is correct
- [ ] Bot is running
- [ ] Added streets to watch list (`/add`)
- [ ] Manually checked for calls (`/check`)
