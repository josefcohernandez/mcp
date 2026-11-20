# Step 05 — Bonus (at home)

Continue after DevFest. Pick any of these — ordered from easiest to harder.

## A. Add a resource

Expose `schedule.json` as an MCP **resource** (read-only context) in addition to tools.

Hints in the [Go SDK examples](https://github.com/modelcontextprotocol/go-sdk/tree/main/examples/server).

## B. Add a prompt

Register a prompt template such as `plan-my-day` that guides the model to call `list_sessions` and build an itinerary.

## C. Read-only annotations

Mark tools with `ReadOnlyHint` so Hosts can treat them as safe.

## D. Deploy HTTP somewhere

Containerize Step 03 (`Dockerfile` sketch below) and deploy to Cloud Run / Fly.io / a VM. Point Cursor at the HTTPS `/mcp` URL.

```dockerfile
FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN go build -o /out/devfest-mcp ./solutions/03-http-transport

FROM gcr.io/distroless/static
COPY --from=build /out/devfest-mcp /devfest-mcp
COPY workshop/data/schedule.json /data/schedule.json
ENV SCHEDULE_PATH=/data/schedule.json
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/devfest-mcp", "-http=:8080"]
```

## E. Read production code (optional)

Browse [`containers/kubernetes-mcp-server`](https://github.com/containers/kubernetes-mcp-server) and map concepts you already know:

| Workshop | Production cousin |
|----------|-------------------|
| `mcp.NewServer` | server bootstrap in `pkg/mcp` |
| `mcp.AddTool` | tool registration / Go SDK adapters |
| Streamable HTTP `/mcp` | `pkg/http` |
| (not in workshop) | toolsets, RBAC, confirmation, OTEL |

Do **not** try to copy that layout into a first server — grow features only when you need them.
