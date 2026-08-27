# bambu2immich

[![CI](https://github.com/eugene-bert/bambu2immich/actions/workflows/ci.yml/badge.svg)](https://github.com/eugene-bert/bambu2immich/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/eugene-bert/bambu2immich)](https://goreportcard.com/report/github.com/eugene-bert/bambu2immich)
[![Go Version](https://img.shields.io/badge/Go-1.27+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Automatically download timelapse videos from your Bambu Lab 3D printer and upload them to [Immich](https://immich.app).

This project is unofficial and is not affiliated with Bambu Lab or Immich.

## How it works

1. Connects to your Bambu Lab printer via MQTT (port 8883)
2. Listens for print completion events
3. Downloads the timelapse video via FTPS (port 990)
4. Uploads to Immich via API
5. Optionally sends a Telegram notification

No Home Assistant required. Runs as a standalone Docker container.

## Supported printers

Any Bambu Lab printer with MQTT and FTPS support:
- A1 / A1 Mini
- P1P / P1S
- X1 / X1C / X1E

Enable **LAN mode** on the printer (Settings → Network) so MQTT and FTPS are reachable on your LAN.

## Quick start

```bash
mkdir -p bambu2immich/data
cd bambu2immich
curl -o docker-compose.yml https://raw.githubusercontent.com/eugene-bert/bambu2immich/main/docker-compose.yml
curl -o .env https://raw.githubusercontent.com/eugene-bert/bambu2immich/main/.env.example
# Edit .env with your printer and Immich details
docker compose up -d
```

No build step needed — pulls the pre-built image from `ghcr.io/eugene-bert/bambu2immich`.

The container runs as uid 1000. Make sure `./data` is writable by that user (`mkdir -p data && chown 1000:1000 data` on Linux).

### Build from source

```bash
git clone https://github.com/eugene-bert/bambu2immich.git
cd bambu2immich
cp .env.example .env
# Edit .env with your printer and Immich details
docker compose up -d --build
```

## Docker network

Compose attaches to an external network named `immich_default` so `IMMICH_URL=http://immich_server:2283` resolves. That is the default network created by the official Immich compose project named `immich`.

If `docker compose up` fails with `network immich_default not found`:

```bash
docker network ls | grep immich
```

Edit `docker-compose.yml` and replace `immich_default` with the network you found.

Alternatively, remove the `networks:` block and point Immich at a published port:

```env
IMMICH_URL=http://host.docker.internal:2283
```

`host.docker.internal` is mapped via `extra_hosts` / `host-gateway` so this also works on Linux.

## Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `BAMBU_IP` | Yes | Printer IP address |
| `BAMBU_ACCESS_CODE` | Yes | LAN access code (printer LCD: Settings → Network) |
| `BAMBU_SERIAL` | Yes | Printer serial number |
| `BAMBU_NAME` | No | Label used in Immich description and local filenames (default: `bambu-a1`) |
| `IMMICH_URL` | Yes | Immich server URL (e.g. `http://immich_server:2283`) |
| `IMMICH_API_KEY` | Yes | Immich API key (Immich → Account Settings → API Keys) |
| `TELEGRAM_BOT_TOKEN` | No | Telegram bot token for notifications |
| `TELEGRAM_CHAT_ID` | No | Telegram chat ID for notifications |
| `KEEP_LOCAL` | No | Keep local copies after upload (default: `false`) |
| `DOWNLOAD_DIR` | No | Download directory (default: `/data/timelapses`) |

## How to get the printer serial number

- On the printer LCD: Settings → Device Info
- In Bambu Studio: Device tab
- On the sticker on the printer

## Security

- Keep `.env` private. It holds the printer access code, Immich API key, and optional Telegram token.
- MQTT and FTPS skip TLS certificate verification. Bambu printers present a self-signed X.509 v1 certificate that standard verification rejects. Traffic stays on your LAN; a machine on the same network could intercept the access code. Do not expose printer MQTT/FTPS to the internet.
- Report vulnerabilities privately — see [SECURITY.md](SECURITY.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). This project follows a [Code of Conduct](CODE_OF_CONDUCT.md).

## License

MIT
