# Contributing

Bug reports and small patches are welcome. Open an issue first if the change is more than a few files.

## Development

```bash
cp .env.example .env
go test -race ./...
go build ./...
```

Docker:

```bash
docker compose up -d --build
```

The compose file expects an external Docker network named `immich_default` unless you switch `IMMICH_URL` to `http://host.docker.internal:2283` and remove the `networks:` block.

## Pull requests

- Keep the change focused. Match existing style.
- Add a test when the behavior is non-trivial (MQTT state, filename sanitization, HTTP client).
- Do not commit `.env`, printer credentials, or sample timelapse videos.

## Release (maintainers)

1. Tag `vX.Y.Z` and push the tag. GitHub Actions builds `linux/amd64` and `linux/arm64` and pushes `X.Y.Z`, `X.Y`, and `latest` to GHCR.
2. After the git repo is public, set the GHCR package visibility to **public** as well (Package settings → Change visibility). Making the repo public does not publish the image.
3. Create a GitHub Release from the tag so the Releases page has notes.
