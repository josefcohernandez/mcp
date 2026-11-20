# Deploy to a Linux VM (Docker + existing HTTPS reverse proxy)

Claude (web, desktop and the Android/iOS apps) reaches custom connectors **from Anthropic's
cloud**, so the server must be public, over **HTTPS on port 443**, and the connector URL must
end in **`/mcp`**.

## Quick install (script)

```bash
curl -fsSL https://raw.githubusercontent.com/josefcohernandez/mcp/main/deploy/install.sh -o install.sh
less install.sh                                  # read it first
bash install.sh                                  # you already have an HTTPS proxy -> then step 2
bash install.sh --domain mcp.example.com         # no proxy: bundled Caddy + Let's Encrypt
```

It installs Docker if missing (asks first), clones to `/opt/devfest-mcp`, builds and starts the
container, checks `/health` and `tools/list`, and prints the connector URL.
`--uninstall` stops everything. Re-running it updates to the latest commit.

## 1. On the VM (manual)

```bash
git clone https://github.com/josefcohernandez/mcp.git && cd mcp   # or git pull
docker compose -f deploy/compose.yaml up -d --build
curl -s http://127.0.0.1:8080/health        # {"status":"ok",...}
```

The container listens only on `127.0.0.1:8080` (change it with `MCP_PORT=...`).

## 2. Reverse proxy (pick yours)

Use a dedicated subdomain such as `mcp.dev.example.com`, or a path on an existing one.

**Caddy**

```caddyfile
mcp.dev.example.com {
    reverse_proxy 127.0.0.1:8080
}
```

**Nginx** (inside the `server { listen 443 ssl; ... }` block)

```nginx
location /mcp {
    proxy_pass http://127.0.0.1:8080;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_buffering off;          # Streamable HTTP may stream responses
    proxy_read_timeout 300s;
}
location = /health {
    proxy_pass http://127.0.0.1:8080;
}
```

**Traefik** (labels on the service in `compose.yaml`)

```yaml
labels:
  - traefik.enable=true
  - traefik.http.routers.devfest-mcp.rule=Host(`mcp.dev.example.com`)
  - traefik.http.routers.devfest-mcp.entrypoints=websecure
  - traefik.http.routers.devfest-mcp.tls.certresolver=letsencrypt
  - traefik.http.services.devfest-mcp.loadbalancer.server.port=8080
```

## 3. Check from outside, then add it to Claude

```bash
curl -s https://mcp.dev.example.com/health
curl -s https://mcp.dev.example.com/mcp \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

claude.ai → Settings → Connectors → **Add custom connector** → `https://mcp.dev.example.com/mcp`.
It then shows up in the Android app after you sign in again.

> The server has **no authentication**: anyone with the URL can call it. It only serves the
> public DevFest agenda read-only; do not reuse this setup for private data without OAuth.

## Update / remove

```bash
git pull && docker compose -f deploy/compose.yaml up -d --build
docker compose -f deploy/compose.yaml down
```
