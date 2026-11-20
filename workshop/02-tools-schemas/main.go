package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
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
		Version: "0.2.0",
	}, nil)

	// Register the three tools. The input schema of each one is inferred from its
	// input struct (ListSessionsInput, GetSessionInput, FindSpeakerInput) in tools.go.
	mcp.AddTool(server, &mcp.Tool{Name: "list_sessions", Description: "List DevFest sessions, optionally filtered by day (Friday/Saturday)"}, makeListSessions(schedule))
	mcp.AddTool(server, &mcp.Tool{Name: "get_session", Description: "Get a single session by id (e.g. s1)"}, makeGetSession(schedule))
	mcp.AddTool(server, &mcp.Tool{Name: "find_speaker", Description: "Find sessions by speaker name (substring match)"}, makeFindSpeaker(schedule))

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
