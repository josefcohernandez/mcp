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
		Version: "0.3.0",
	}, nil)
	registerTools(server, schedule)

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

	// TODO: create Streamable HTTP handler and mount /mcp + /health
	//
	// mcpHandler := mcp.NewStreamableHTTPHandler(
	// 	func(_ *http.Request) *mcp.Server { return server },
	// 	&mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true},
	// )
	// started := time.Now()
	// mux := http.NewServeMux()
	// mux.Handle("/mcp", mcpHandler)
	// mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodGet {
	// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// 	_ = json.NewEncoder(w).Encode(map[string]any{
	// 		"status": "ok", "service": "devfest-mcp",
	// 		"uptime_s": int64(time.Since(started).Seconds()),
	// 	})
	// })
	// log.Printf("MCP Streamable HTTP: http://127.0.0.1%s/mcp", addr)
	// log.Fatal(http.ListenAndServe(addr, mux))

	_, _, _, _ = addr, json.Marshal, time.Now, http.ListenAndServe
	log.Fatal("TODO: wire Streamable HTTP on /mcp and /health — see README.md")
}

func registerTools(server *mcp.Server, schedule *Schedule) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_sessions",
		Description: "List DevFest sessions, optionally filtered by day (Friday/Saturday)",
	}, makeListSessions(schedule))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_session",
		Description: "Get a single session by id (e.g. s1)",
	}, makeGetSession(schedule))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "find_speaker",
		Description: "Find sessions by speaker name (substring match)",
	}, makeFindSpeaker(schedule))
}
