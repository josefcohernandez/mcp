# Step 04 — Connect Cursor

**Timebox:** ~7 minutes

## Goal

Point an MCP **Host** (Cursor) at the server you built and ask questions that force tool calls.

## 1. Start the HTTP server

Prefer the completed Step 03 (or the solution):

```bash
cd solutions/03-http-transport   # or workshop/03-http-transport once finished
go mod tidy
go run . -http=:8080
```

Leave it running.

## 2. Configure Cursor

Create or edit **project** config:

`.cursor/mcp.json` (in the repo root) **or** global `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "devfest-mcp": {
      "url": "http://127.0.0.1:8080/mcp"
    }
  }
}
```

Notes:

- For Streamable HTTP, Cursor expects the key **`url`** (not `command` / `args`).
- Include the **`/mcp`** path.
- Reload MCP servers in Cursor (or restart the window) after saving.

### Alternative: stdio (subprocess)

If you prefer launching the binary as a child process:

```json
{
  "mcpServers": {
    "devfest-mcp": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/ABS/PATH/TO/devfestValencia2026/solutions/03-http-transport"
    }
  }
}
```

Set `SCHEDULE_PATH` in `env` if the working directory cannot see `workshop/data/schedule.json`.

## 3. Try these prompts

- “List Friday sessions from the DevFest MCP server.”
- “Who is speaking in session s5?”
- “Find sessions for speaker Alberto.”

You should see tool calls to `list_sessions`, `get_session`, or `find_speaker`.

## Claude Desktop (stdio)

Use the stdio form with `command` pointing at a built binary:

```bash
cd solutions/03-http-transport && go build -o ../../bin/devfest-mcp .
```

Then configure Claude Desktop with that binary path (no `-http`).

## Done when

The model answers using **live tool results** from your local schedule JSON — not from training data alone.
