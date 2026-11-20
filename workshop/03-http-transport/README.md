# Step 03 — Streamable HTTP transport

**Timebox:** ~8 minutes

## Goal

Serve the same tools over **Streamable HTTP** so browser-based and remote clients can connect without spawning a subprocess.

You will:

1. Keep stdio as the default (no `-http` flag)
2. Add `-http=:8080` to expose `POST/GET/DELETE /mcp`
3. Add a simple `GET /health` endpoint

This is the transport pattern used by many hosted MCP servers (and by Go examples in the official SDK).

## What you will edit

Complete the `TODO` blocks in `main.go` around `mcp.NewStreamableHTTPHandler`.

## Run

```bash
cd workshop/03-http-transport
go mod tidy

# HTTP mode
go run . -http=:8080
```

Check health:

```bash
curl -s http://127.0.0.1:8080/health
```

## Done when

- Server logs the `/mcp` URL
- `/health` returns JSON `{"status":"ok",...}`

## Next

Wire Cursor to `http://127.0.0.1:8080/mcp` in [Step 04](../04-connect-cursor).

## Stuck?

Use [`solutions/03-http-transport`](../../solutions/03-http-transport).
