# Installation Guide (Docker Only)

This project is deployed with Docker Compose only.

## 1) Build Linux binary

From `telegram-cad-bot/`:

```bash
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o cadbot ./cmd/bot
```

> brew install filosottile/musl-cross/musl-cross is needed to run this

Use `GOARCH=arm64` for ARM servers.

## 2) Prepare deploy bundle

Use files in `deploy/docker/`:

- `deploy/docker/docker-compose.yml`
- `deploy/docker/setup.sh`
- `deploy/docker/.env.example`

Create `.env` from the example and set:

- `TELEGRAM_BOT_TOKEN`
- `CAD_BASE_URL`
- `CAD_AGENCY_ID`
- `CHECK_INTERVAL` (optional)
- `VERIFY_SSL` (optional)

## 3) Server layout

In your server deploy directory (example: `server-setup/cad-bot/`):

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

Notes:
- `cadbot` is the compiled Linux binary
- `cad_calls/` must contain `direct_api_post.py`
- Put an existing `cadbot.db` in `data/` to migrate state

## 4) Start

```bash
docker-compose up -d
```

## 5) Logs

```bash
docker-compose logs -f cadbot
```

## 6) Update workflow

1. Pull latest code
2. Rebuild `cadbot` binary
3. Replace deploy files from `deploy/docker/` if changed
4. Restart:

```bash
docker-compose up -d --force-recreate
```
