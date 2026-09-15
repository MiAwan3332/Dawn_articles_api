# Dawn Articles API

A small Go API that reads Dawn's official latest-news RSS feed and exposes it as JSON.

## Endpoints

- `GET /api/v1/latest-news`
- `GET /api/v1/article/:id`
- `GET /swagger/index.html`

Article details are available while an article remains in Dawn's latest-news feed.

## Docker

```bash
docker compose up -d --build
docker compose ps
docker compose logs -f api
```

The Compose service listens on `127.0.0.1:3004` and is intended to be exposed through a reverse proxy.

## Production

The LCA deployment is installed at `/opt/dawn_articles_api` and served through Nginx at:

```text
https://dawn-api.lca-portal.com
```

To deploy updates:

```bash
cd /opt/dawn_articles_api
git pull --ff-only
docker compose up -d --build
```
