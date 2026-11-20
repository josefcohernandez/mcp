# demo-example — live DevFest MCP server

MCP server that exposes one tool, **`schedule`**, which **fetches the real DevFest Valencia 2026 agenda live** from the GDG Community API and filters it.

Event page:

https://gdg.community.dev/events/details/google-gdg-valencia-presents-devfest-valencia-2026-1/

API:

```text
GET https://gdg.community.dev/api/event_slim/google-gdg-valencia-presents-devfest-valencia-2026-1/
```

## Tool: `schedule`

All parameters are **optional**. Combine freely.

| Param | Matches |
|-------|---------|
| `day` | `Friday` / `Saturday` (also `Viernes` / `Sábado`) |
| `time` | Session clock — accepts `18:05`, `6:05 PM`, etc. |
| `person` | Substring of the **speaker** name |
| `description` | Substring searched in **title** and **description** |

Examples the model (or you) can call:

```json
{}
```

```json
{ "day": "Friday" }
```

```json
{ "person": "Alberto" }
```

```json
{ "description": "MCP" }
```

```json
{ "day": "Saturday", "time": "10:00" }
```

Each call hits the live API, parses Friday + Saturday, then applies filters.


### MCP Inspector (`npx`)

[MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is a local UI to list tools and call them without Cursor/Gemini. Requires Node.js.

**Stdio** (Inspector spawns the server):

```bash
cd demo-example
npx @modelcontextprotocol/inspector go run .
```

**HTTP** (start the server yourself first):

```bash
# terminal 1
go run . -http=:8080

# terminal 2
npx @modelcontextprotocol/inspector --transport http --server-url http://127.0.0.1:8080/mcp
```

In the UI, open tool **`schedule`** and try e.g. `{"person":"Alberto"}` or `{"day":"Friday","description":"MCP"}`.

Try prompts like:

- “What is on the DevFest schedule on Friday?”
- “Find Alberto’s session with the schedule tool.”
- “Any talk about MCP or agents?”

## Run as MCP (stdio)

```bash
cd demo-example
go mod tidy
go run .
```

Wire Cursor / Claude to this folder with a stdio server, for example `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "devfest-live": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/ABS/PATH/TO/devfestValencia2026/demo-example"
    }
  }
}
```

## Run as MCP (HTTP) — great for live demos

```bash
go run . -http=:8080
```

### Cursor

Project `.cursor/mcp.json` or `~/.cursor/mcp.json` — use **`url`** for Streamable HTTP:

```json
{
  "mcpServers": {
    "devfest-live": {
      "url": "http://127.0.0.1:8080/mcp"
    }
  }
}
```

### Gemini CLI

Edit `~/.gemini/settings.json` (user-wide) or project `.gemini/settings.json`.  
For Streamable HTTP Gemini expects **`httpUrl`** (not `url` — that is SSE):

```json
{
  "mcpServers": {
    "devfest-live": {
      "httpUrl": "http://127.0.0.1:8080/mcp"
    }
  }
}
```

Or add it from the terminal:

```bash
gemini mcp add --transport http --scope user devfest-live http://127.0.0.1:8080/mcp
```

**What `--scope` means**

| Scope | Writes to | Who sees the server |
|-------|-----------|---------------------|
| `user` | `~/.gemini/settings.json` | Every project on your machine |
| `project` (default if omitted) | `.gemini/settings.json` in the current repo | Only this project |

Use `--scope user` for a live demo laptop setup; use project scope when you want the MCP config committed (or local) to this workshop repo only.

Keep `go run . -http=:8080` running, then in Gemini CLI check with `/mcp`.

**Stdio** alternative for Gemini (no `-http` flag on the server):

```json
{
  "mcpServers": {
    "devfest-live": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/ABS/PATH/TO/devfestValencia2026/demo-example"
    }
  }
}
```


## Refresh the workshop snapshot

Workshop steps 02–04 use a **local** file so attendees do not need network. Refresh it from live data:

```bash
go run . -dump=../workshop/data/schedule.json
```

## Flags

| Flag | Meaning |
|------|---------|
| *(none)* | MCP over **stdio** |
| `-http=:8080` | MCP Streamable HTTP on `/mcp` (+ `/health`) |
| `-dump=PATH` | One-shot fetch → write JSON → exit |
| `-slug=...` | Override event slug |
| `-event-id=133556` | Use `/api/event/{id}/` instead of `event_slim` |
