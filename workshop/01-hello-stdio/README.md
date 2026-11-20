# Step 01 — Hello MCP (stdio)

**Timebox:** ~10 minutes

## Goal

Create the smallest useful MCP server in Go:

1. Construct an `mcp.Server`
2. Register one tool (`greet`)
3. Run it over **stdio** (`mcp.StdioTransport`)

This matches how many local MCP clients launch servers as a subprocess.

## What you will edit

Open `main.go`. Complete the `TODO` markers.

## Run

```bash
cd workshop/01-hello-stdio
go mod tidy
go run .
```

The process waits on stdin/stdout — that is expected. Stop with `Ctrl+C`.

### Smoke-test without an AI client (optional)

With the [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) (Node.js required):

```bash
npx @modelcontextprotocol/inspector go run .
```

## Done when

- `go run .` starts without errors
- Inspector (or Cursor later) lists a tool named `greet`

## Stuck?

Copy [`solutions/01-hello-stdio`](../../solutions/01-hello-stdio) and move on.
