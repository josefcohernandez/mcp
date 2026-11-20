package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newTestSession starts the devfest MCP server behind an httptest server using the same
// Streamable HTTP options as main(), and connects the official Go SDK client to /mcp.
// This is exactly what an MCP Host (Cursor, Claude, ...) does in Step 04.
func newTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	schedule, err := LoadSchedule(filepath.Join("..", "data", "schedule.json"))
	if err != nil {
		t.Fatalf("load schedule: %v", err)
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "devfest-mcp", Version: "test"}, nil)
	registerTools(server, schedule)

	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true},
	)
	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "test-host", Version: "0"}, nil)
	cs, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: ts.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func call(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	return res
}

func decodeList(t *testing.T, res *mcp.CallToolResult) SessionList {
	t.Helper()
	if res.IsError {
		t.Fatalf("unexpected tool error: %+v", res.Content)
	}
	// One content block only: no duplicated JSON fallback (see SessionList).
	if len(res.Content) != 1 {
		t.Fatalf("want 1 content block, got %d", len(res.Content))
	}
	var out SessionList
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return out
}

func TestListTools(t *testing.T) {
	cs := newTestSession(t)
	res, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, tool := range res.Tools {
		got[tool.Name] = true
	}
	for _, want := range []string{"list_sessions", "get_session", "find_speaker"} {
		if !got[want] {
			t.Errorf("tool %q not listed", want)
		}
	}
}

func TestListSessionsFriday(t *testing.T) {
	cs := newTestSession(t)
	for _, day := range []string{"Friday", "viernes"} {
		out := decodeList(t, call(t, cs, "list_sessions", map[string]any{"day": day}))
		if out.Count == 0 || out.Count != len(out.Sessions) {
			t.Fatalf("%s: bad count %d / %d", day, out.Count, len(out.Sessions))
		}
		for _, s := range out.Sessions {
			if s.Day != "Friday" {
				t.Errorf("%s: got session %s on %s", day, s.ID, s.Day)
			}
		}
	}
}

func TestGetSession(t *testing.T) {
	cs := newTestSession(t)
	res := call(t, cs, "get_session", map[string]any{"id": "s5"})
	var s Session
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &s); err != nil {
		t.Fatal(err)
	}
	if s.ID != "s5" || s.Speaker == "" {
		t.Fatalf("unexpected session: %+v", s)
	}

	missing := call(t, cs, "get_session", map[string]any{"id": "nope"})
	if !missing.IsError {
		t.Fatal("unknown id should return a tool error (IsError), not a protocol error")
	}
}

func TestFindSpeaker(t *testing.T) {
	cs := newTestSession(t)
	if out := decodeList(t, call(t, cs, "find_speaker", map[string]any{"name": "alberto"})); out.Count == 0 {
		t.Fatal("expected at least one session for 'alberto' (case-insensitive)")
	}
	out := decodeList(t, call(t, cs, "find_speaker", map[string]any{"name": "nobody-xyz"}))
	if out.Count != 0 || out.Sessions == nil {
		t.Fatalf("no match should be an empty list, got %+v", out)
	}
}

func TestInvalidArguments(t *testing.T) {
	cs := newTestSession(t)
	// The input schema is generated from GetSessionInput: "id" is required and unknown
	// properties are rejected, so the SDK refuses this call before our handler runs.
	if res := call(t, cs, "get_session", map[string]any{"session": "s5"}); !res.IsError {
		t.Fatal("expected schema validation error")
	}
}
