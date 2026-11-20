# DevFest Valencia 2026 — Build Your First MCP Server

**Hands-on workshop (~75 minutes):** build a Model Context Protocol (MCP) server in Go, step by step, using the official [`go-sdk`](https://github.com/modelcontextprotocol/go-sdk) — the same core library used by production servers such as [`kubernetes-mcp-server`](https://github.com/containers/kubernetes-mcp-server).

| | |
|---|---|
| **Slides** | [GitHub Pages](https://aljesusg.github.io/devfestValencia2026/) |
| **Event** | DevFest Valencia 2026 |
| **Language** | English |
| **Stack** | Go · `github.com/modelcontextprotocol/go-sdk` |

> We intentionally keep this workshop small. You will **not** replicate kubernetes-mcp-server architecture (toolsets, RBAC, OAuth, OpenTelemetry, etc.). You will learn the same methodology and core APIs: create a server, register tools, run over stdio, then expose Streamable HTTP and connect a client.

## Agenda (~75 minutes)

| Time | Block | What you do |
|------|--------|-------------|
| 0–15 min | Slides · intro | LLM context gap → what MCP is (Host / Client / Server, tools, transports, Go SDK) |
| 15–18 min | Setup | Clone this repo, confirm Go works |
| 18–30 min | Live demo | Live `schedule` MCP demo (shared during the talk), sample prompts, Inspector, token design + evals |
| 30–40 min | [Step 01](workshop/01-hello-stdio) | Minimal MCP server over stdio + one tool |
| 40–50 min | [Step 02](workshop/02-tools-schemas) | Multiple tools with input schemas |
| 50–58 min | [Step 03](workshop/03-http-transport) | Streamable HTTP (`/mcp`) |
| 58–68 min | [Step 04](workshop/04-connect-cursor) | Connect Cursor (or another MCP client) and ask live questions |
| 68–75 min | Wrap-up | Takeaways + [Step 05](workshop/05-bonus-home) for home |

If you fall behind, use the matching folder from `solutions/` once it is published for that step.

## Prerequisites

Install **before** the session:

1. **Go 1.25+** (Go 1.24+ usually works via the Go toolchain download: `GOTOOLCHAIN=auto`)
   ```bash
   go version
   ```
2. **Git**
3. An MCP-capable client — ideally **[Cursor](https://cursor.com/)** (or Claude Desktop / Gemini CLI)

Optional for Step 03+: a free terminal port (default `8080`).

## Quick start

```bash
git clone https://github.com/aljesusg/devfestValencia2026.git
cd devfestValencia2026

# Step 01 — hello MCP over stdio
cd workshop/01-hello-stdio
go mod tidy
go run .
```

## Repository layout

```text
docs/                 # Reveal.js slides (GitHub Pages)
workshop/
  01-hello-stdio/     # starter for live coding
  02-tools-schemas/
  03-http-transport/
  04-connect-cursor/
  05-bonus-home/      # homework
  data/               # local snapshot of the real agenda (for offline steps)
solutions/            # filled in step by step during the workshop
```

The live demo server (`demo-example/`) and completed solutions are published during the session, not in the starter clone.

## Domain of the demo server

Workshop tools read a **local snapshot** of the real DevFest Valencia 2026 agenda (`workshop/data/schedule.json`) — no Kubernetes, no cloud credentials.

During the talk we also run a live MCP server with a **`schedule`** tool (optional `day`, `time`, `person`, `description`) that fetches Friday + Saturday in real time from the GDG Community API.

Event page: [DevFest Valencia 2026](https://gdg.community.dev/events/details/google-gdg-valencia-presents-devfest-valencia-2026-1/).

## Methodology (inspired by kubernetes-mcp-server)

1. Create an `mcp.Server` with an `Implementation` name/version  
2. Register **tools** (the interesting part for agents)  
3. Serve over a **transport** (`StdioTransport` → Streamable HTTP)  
4. Point a **Host** (Cursor, Claude, …) at your server  

Same idea as production Go MCP servers — without the production surface area.

## Useful links

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [kubernetes-mcp-server](https://github.com/containers/kubernetes-mcp-server) (reference architecture, not a dependency of this workshop)
- [Workshop slides](https://aljesusg.github.io/devfestValencia2026/)

## License

See [LICENSE](LICENSE).
