package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	scheduleURI     = "devfest://schedule"
	sessionTemplate = "devfest://sessions/{id}"
)

// registerResources exposes read-only context (Bonus A).
//
// Tools are *model-controlled* (the LLM decides to call them); resources are
// *application-controlled*: the Host or the user chooses to attach them to the
// conversation, e.g. via an "@" menu, without spending a tool call.
func registerResources(server *mcp.Server, schedule *Schedule, dataPath string) {
	// A fixed resource: the whole snapshot, read from disk on every request so an
	// updated schedule.json (e.g. refreshed with demo-example -dump) is picked up.
	server.AddResource(&mcp.Resource{
		URI:         scheduleURI,
		Name:        "schedule",
		Title:       "DevFest Valencia 2026 schedule",
		Description: "Full DevFest schedule snapshot (JSON)",
		MIMEType:    "application/json",
	}, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		raw, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{
			URI:      scheduleURI,
			MIMEType: "application/json",
			Text:     string(raw),
		}}}, nil
	})

	// A resource template (RFC 6570 URI template): one addressable resource per session.
	server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: sessionTemplate,
		Name:        "session",
		Title:       "DevFest session",
		Description: "A single DevFest session by id, e.g. devfest://sessions/s5",
		MIMEType:    "application/json",
	}, func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		uri := req.Params.URI
		id := strings.TrimPrefix(uri, "devfest://sessions/")
		sess, err := schedule.ByID(id)
		if err != nil {
			return nil, mcp.ResourceNotFoundError(uri)
		}
		b, err := json.MarshalIndent(sess, "", "  ")
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{
			URI:      uri,
			MIMEType: "application/json",
			Text:     string(b),
		}}}, nil
	})
}

// registerPrompts adds a reusable prompt template (Bonus B).
//
// Prompts are *user-controlled*: Hosts usually surface them as slash commands
// (e.g. /plan-my-day). The server only returns messages; the model then decides
// which tools to call to fulfil them.
func registerPrompts(server *mcp.Server) {
	server.AddPrompt(&mcp.Prompt{
		Name:        "plan-my-day",
		Title:       "Plan my DevFest day",
		Description: "Build a personal itinerary for one DevFest day using the schedule tools",
		Arguments: []*mcp.PromptArgument{
			{Name: "day", Description: "Friday or Saturday (Viernes / Sábado also work)", Required: true},
			{Name: "interests", Description: "Topics you care about, e.g. \"agents, Flutter, careers\""},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		day := strings.TrimSpace(req.Params.Arguments["day"])
		if day == "" {
			return nil, fmt.Errorf("argument %q is required", "day")
		}
		interests := strings.TrimSpace(req.Params.Arguments["interests"])
		if interests == "" {
			interests = "anything interesting; give a balanced mix"
		}
		text := fmt.Sprintf(`Plan my day at DevFest Valencia 2026 for %s.
My interests: %s.

1. Call the list_sessions tool with {"day": %q} to get the real agenda. Do not invent sessions.
2. Use get_session or find_speaker if you need more detail about a talk or speaker.
3. Return a chronological itinerary: time, title, speaker and one line on why it matches my interests.
4. Keep breaks and networking slots, and point out any time where two sessions I would like overlap.`,
			day, interests, day)
		return &mcp.GetPromptResult{
			Description: "Itinerary for " + day,
			Messages: []*mcp.PromptMessage{{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			}},
		}, nil
	})
}
