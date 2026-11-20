package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var httpListen = flag.String("http", "", "listen address for Streamable HTTP (e.g. :8080); omit for stdio")

func main() {
	flag.Parse()

	dataPath := os.Getenv("SCHEDULE_PATH")
	if dataPath == "" {
		dataPath = filepath.Join("..", "data", "schedule.json")
	}

	schedule, err := LoadSchedule(dataPath)
	if err != nil {
		log.Fatalf("load schedule: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "devfest-mcp",
		Version: "0.5.0",
	}, nil)
	registerTools(server, schedule)
	registerResources(server, schedule, dataPath)
	registerPrompts(server)

	if *httpListen == "" {
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
		}
		return
	}

	addr := *httpListen
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	// Streamable HTTP: a single endpoint (/mcp) that accepts POST (JSON-RPC messages),
	// GET (optional server->client SSE stream) and DELETE (end session).
	// Stateless: no Mcp-Session-Id is kept between requests (fine for pure, read-only tools).
	// JSONResponse: reply with application/json instead of an SSE stream.
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true},
	)

	started := time.Now()
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":   "ok",
			"service":  "devfest-mcp",
			"uptime_s": int64(time.Since(started).Seconds()),
		})
	})

	log.Printf("MCP Streamable HTTP: http://127.0.0.1%s/mcp", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// readOnly marks a tool as safe for Hosts: it never modifies anything (ReadOnlyHint),
// repeating a call gives the same result (IdempotentHint) and it only touches the local
// snapshot, not the outside world (OpenWorldHint=false). Hosts may use these hints to
// skip confirmation prompts; they are hints, not a security boundary.
func readOnly(title string) *mcp.ToolAnnotations {
	openWorld := false
	return &mcp.ToolAnnotations{
		Title:          title,
		ReadOnlyHint:   true,
		IdempotentHint: true,
		OpenWorldHint:  &openWorld,
	}
}

func registerTools(server *mcp.Server, schedule *Schedule) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_sessions",
		Description: "List DevFest sessions, optionally filtered by day (Friday/Saturday)",
		Annotations: readOnly("List sessions"),
	}, makeListSessions(schedule))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_session",
		Description: "Get a single session by id (e.g. s1)",
		Annotations: readOnly("Get session"),
	}, makeGetSession(schedule))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "find_speaker",
		Description: "Find sessions by speaker name (substring match)",
		Annotations: readOnly("Find speaker"),
	}, makeFindSpeaker(schedule))
}
