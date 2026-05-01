# Quick Install (Docker)

## Fast path

1. Build Linux binary:

```bash
cd telegram-cad-bot
```bash
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o cadbot ./cmd/bot
```

> brew install filosottile/musl-cross/musl-cross is needed to run this

2. Copy to server deploy folder:

- `cadbot` binary
- `deploy/docker/docker-compose.yml` as `docker-compose.yml`
- `deploy/docker/setup.sh` as `setup.sh`
- `deploy/docker/.env.example` as `.env` (then edit values)
- `cad_calls/` folder (must include `direct_api_post.py`)

3. Start:

```bash
docker-compose up -d
docker-compose logs -f cadbot
```

## Full docs

- [INSTALL.md](INSTALL.md)
- [deploy/docker/README.md](deploy/docker/README.md)
