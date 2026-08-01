# bambu2immich

Automatically download timelapse videos from your Bambu Lab 3D printer and upload them to [Immich](https://immich.app).

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

## Quick start

```bash
git clone https://github.com/eugene-bert/bambu2immich.git
cd bambu2immich
cp .env.example .env
# Edit .env with your printer and Immich details
docker compose up -d
```

## Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `BAMBU_IP` | Yes | Printer IP address |
| `BAMBU_ACCESS_CODE` | Yes | LAN access code (printer LCD: Settings > Network) |
| `BAMBU_SERIAL` | Yes | Printer serial number |
| `IMMICH_URL` | Yes | Immich server URL (e.g. `http://immich:2283`) |
| `IMMICH_API_KEY` | Yes | Immich API key (Immich > Account Settings > API Keys) |
| `TELEGRAM_BOT_TOKEN` | No | Telegram bot token for notifications |
| `TELEGRAM_CHAT_ID` | No | Telegram chat ID for notifications |
| `KEEP_LOCAL` | No | Keep local copies after upload (default: `false`) |
| `DOWNLOAD_DIR` | No | Download directory (default: `/data/timelapses`) |

## How to get the printer serial number

- On the printer LCD: Settings > Device Info
- In Bambu Studio: Device tab
- On the sticker on the printer

## License

MIT
