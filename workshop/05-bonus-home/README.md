# Step 05 — Bonus (at home)

> **Implemented in this folder.** It starts from Step 03 and adds bonuses A–D:
>
> | Bonus | Where | What |
> |---|---|---|
> | A. Resource | `extras.go` | `devfest://schedule` (full JSON) + template `devfest://sessions/{id}` |
> | B. Prompt | `extras.go` | `plan-my-day` (`day` required, `interests` optional) |
> | C. Annotations | `main.go` | every tool: `readOnlyHint`, `idempotentHint`, `openWorldHint=false` |
> | D. Container | `Dockerfile` | static binary on `distroless/static`, non-root, ~18 MB |
>
> ```bash
> cd workshop/05-bonus-home
> go test -race ./...     # 9 end-to-end tests over Streamable HTTP
> go run . -http=:8080
> ```
>
> Notes against SDK v1.8.0:
> - `ReadResourceResult.Contents` is `[]*mcp.ResourceContents` (with `Text`/`Blob`); there is no
>   `TextResourceContents` type, so the sketch in `solutions/05-bonus-home` does not compile as is.
> - Build the container with `CGO_ENABLED=0`: otherwise the binary links against glibc and
>   `distroless/static` fails with `exec ...: no such file or directory`.

Continue after DevFest. Pick any of these — ordered from easiest to harder.

## A. Add a resource

Expose `schedule.json` as an MCP **resource** (read-only context) in addition to tools.

Hints in the [Go SDK examples](https://github.com/modelcontextprotocol/go-sdk/tree/main/examples/server).

## B. Add a prompt

Register a prompt template such as `plan-my-day` that guides the model to call `list_sessions` and build an itinerary.

## C. Read-only annotations

Mark tools with `ReadOnlyHint` so Hosts can treat them as safe.

## D. Deploy HTTP somewhere

Containerize Step 03 (`Dockerfile` sketch below; the tested version is [`Dockerfile`](Dockerfile), build it from the repo root with `docker build -f workshop/05-bonus-home/Dockerfile -t devfest-mcp .`) and deploy to Cloud Run / Fly.io / a VM. Point Cursor at the HTTPS `/mcp` URL.

```dockerfile
FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/devfest-mcp ./solutions/03-http-transport

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
