package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type runtimeConfig struct {
	Slug    string
	EventID int
}

func main() {
	httpListen := flag.String("http", "", "listen address for Streamable HTTP (e.g. :8080); omit for stdio MCP")
	slug := flag.String("slug", defaultSlug, "event slug on gdg.community.dev")
	eventID := flag.Int("event-id", 0, "optional numeric event id (uses /api/event/{id}/)")
	dump := flag.String("dump", "", "fetch once and write JSON snapshot then exit (e.g. ../workshop/data/schedule.json)")
	flag.Parse()

	cfg := runtimeConfig{Slug: *slug, EventID: *eventID}
	log.Printf("config: slug=%q event_id=%d", cfg.Slug, cfg.EventID)

	if *dump != "" {
		log.Printf("dump: writing schedule snapshot to %s", *dump)
		if err := dumpSchedule(cfg, *dump); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", *dump)
		return
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "devfest-live-mcp",
		Version: "0.1.0",
	}, nil)
	registerScheduleTool(server, cfg)
	log.Printf("mcp: registered tool %q", "schedule")

	if *httpListen == "" {
		log.Printf("MCP stdio server ready (tool: schedule) — live agenda from gdg.community.dev")
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
		}
		return
	}

	addr := *httpListen
	if p := os.Getenv("PORT"); p != "" {
		log.Printf("mcp: PORT env overrides listen addr → :%s", p)
		addr = ":" + p
	}

	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true},
	)

	started := time.Now()
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("http: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":   "ok",
			"service":  "devfest-live-mcp",
			"tool":     "schedule",
			"uptime_s": int64(time.Since(started).Seconds()),
		})
	})

	log.Printf("MCP Streamable HTTP: http://127.0.0.1%s/mcp", addr)
	log.Printf("Tool: schedule (optional day, time, person, description) — fetches live agenda")
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func dumpSchedule(cfg runtimeConfig, path string) error {
	started := time.Now()
	payload, err := FetchEvent(cfg.Slug, cfg.EventID)
	if err != nil {
		return err
	}
	schedule, err := ParseSchedule(payload)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(schedule, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return err
	}
	log.Printf("dump: wrote %d bytes (%d sessions) to %s in %s",
		len(b), len(schedule.Sessions), path, time.Since(started).Round(time.Millisecond))
	return nil
}
