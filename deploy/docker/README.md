# Docker Deploy Bundle

This folder is the canonical Docker deployment bundle for server installs and updates.

## Files

- `docker-compose.yml`: container runtime definition
- `setup.sh`: runtime bootstrap (installs Python deps, writes `config.py`, starts bot)
- `.env.example`: required environment values

## First-time server setup

1. Build Linux binary from `telegram-cad-bot/`:

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o cadbot ./cmd/bot
```

2. Copy these to your server deployment directory (for example `server-setup/cad-bot/`):
   - `deploy/docker/docker-compose.yml`
   - `deploy/docker/setup.sh`
   - `deploy/docker/.env.example` (rename to `.env`)
   - `telegram-cad-bot/cadbot` (compiled binary)
   - root `cad_calls/` source folder (must include `direct_api_post.py`)

3. Ensure runtime layout:

```text
cad-bot/
  cadbot
  docker-compose.yml
  setup.sh
  .env
  data/
  cad_calls/
    direct_api_post.py
    config.py
```

4. Start:

```bash
docker-compose up -d
docker-compose logs -f cadbot
```

## Admin visibility command

Set `ADMIN_USER_IDS` in `.env` with one or more Telegram user IDs (comma-separated) to enable `/users`.

Example:

```env
ADMIN_USER_IDS=123456789
```

Users in that list can run `/users` in Telegram to see all users and what streets they monitor.

## Updating after changes

1. Pull latest source
2. Rebuild Linux binary (`CGO_ENABLED=1`)
3. Replace `cadbot`
4. If deploy config changed, copy updated files from `deploy/docker/`
5. Restart:

```bash
docker-compose up -d --force-recreate
```
