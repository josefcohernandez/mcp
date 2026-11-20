# Step 02 — Tools & schemas

**Timebox:** ~10 minutes

## Goal

Grow the server from one toy tool to a small **DevFest schedule** API:

| Tool | Purpose |
|------|---------|
| `list_sessions` | List sessions (optional `day` filter: Friday / Saturday) |
| `get_session` | Get one session by `id` |
| `find_speaker` | Find sessions by speaker name (substring match) |

Data lives in [`../data/schedule.json`](../data/schedule.json) — a **local snapshot** of the real DevFest Valencia agenda (no network during the exercise).

Refresh that snapshot anytime from the live GDG API:

```bash
cd ../../demo-example && go run . -dump=../workshop/data/schedule.json
```

`mcp.AddTool` builds JSON Schema from your Go structs (`json` / `jsonschema` tags). That is the same idea production Go MCP servers rely on: typed inputs the model can call safely.

## What you will edit

1. Finish the three handlers in `tools.go` (`TODO`)
2. Register them in `main.go`

## Run

```bash
cd workshop/02-tools-schemas
go mod tidy
go run .
```

Optional: `npx @modelcontextprotocol/inspector go run .`

## Done when

- Three tools appear in the inspector
- `list_sessions` with `{"day":"Friday"}` returns only Friday sessions

## Stuck?

Use [`solutions/02-tools-schemas`](../../solutions/02-tools-schemas).
