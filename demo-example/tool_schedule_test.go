package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// fakeAgenda mimics the Bevy event_slim payload. The agenda is embedded as a JSON
// *string*, which also exercises ParseSchedule's fallback path.
func fakeAgenda(t *testing.T) []byte {
	t.Helper()
	agenda, _ := json.Marshal(AgendaRoot{Multiday: true, Days: []AgendaDay{
		{Title: "Viernes", Agenda: []AgendaItem{
			{Time: "4:00 PM", Activity: "Acceso al evento"},
			{Time: "18:05", Activity: "Alberto Gutiérrez - Taller Práctico: Construye tu primer servidor MCP paso a paso"},
		}},
		{Title: "Sábado", Agenda: []AgendaItem{
			{Time: "10:00", Activity: "Javier Gallego - De Junior a Tech Lead", Description: "Carrera profesional"},
		}},
	}})
	body, _ := json.Marshal(map[string]any{
		"id": 1, "title": "DevFest Valencia 2026", "url": "https://example.test/event",
		"agenda": string(agenda), "timezone_abbreviation": "CET",
		"venue_name": "La Harinera", "venue_city": "València",
	})
	return body
}

// connect wires the schedule tool to a fake upstream API and returns an MCP client
// session talking to the server over in-memory transports.
func connect(t *testing.T, upstream http.HandlerFunc) *mcp.ClientSession {
	t.Helper()
	api := httptest.NewServer(upstream)
	t.Cleanup(api.Close)
	old := apiBase
	apiBase = api.URL
	t.Cleanup(func() { apiBase = old })

	server := mcp.NewServer(&mcp.Implementation{Name: "devfest-live-mcp", Version: "test"}, nil)
	registerScheduleTool(server, runtimeConfig{Slug: defaultSlug})

	ct, st := mcp.NewInMemoryTransports()
	ctx := context.Background()
	ss, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ss.Close() })
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func callSchedule(t *testing.T, cs *mcp.ClientSession, args map[string]any) (*mcp.CallToolResult, ScheduleToolResult) {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "schedule", Arguments: args})
	if err != nil {
		t.Fatal(err)
	}
	var out ScheduleToolResult
	if !res.IsError {
		if len(res.Content) != 1 {
			t.Fatalf("want a single content block, got %d", len(res.Content))
		}
		if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &out); err != nil {
			t.Fatal(err)
		}
	}
	return res, out
}

func TestScheduleToolEndToEnd(t *testing.T) {
	var gotPath string
	cs := connect(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write(fakeAgenda(t))
	})

	tools, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if a := tools.Tools[0].Annotations; a == nil || !a.ReadOnlyHint || a.OpenWorldHint == nil || !*a.OpenWorldHint {
		t.Errorf("schedule should be annotated read-only and open-world, got %+v", a)
	}

	_, all := callSchedule(t, cs, map[string]any{})
	if gotPath != "/api/event_slim/"+defaultSlug+"/" {
		t.Errorf("unexpected upstream path %q", gotPath)
	}
	if all.Count != 3 || all.Event != "DevFest Valencia 2026" || all.Location != "La Harinera, València" {
		t.Fatalf("unexpected result: %+v", all)
	}

	// 12h and 24h clocks are normalized before comparing.
	_, byTime := callSchedule(t, cs, map[string]any{"time": "6:05 PM"})
	if byTime.Count != 1 || byTime.Sessions[0].Speaker != "Alberto Gutiérrez" {
		t.Fatalf("time filter: %+v", byTime)
	}
	_, byDay := callSchedule(t, cs, map[string]any{"day": "sabado", "description": "carrera"})
	if byDay.Count != 1 || byDay.Sessions[0].ID != "s3" {
		t.Fatalf("day+description filter: %+v", byDay)
	}

	_, none := callSchedule(t, cs, map[string]any{"person": "nobody"})
	if none.Count != 0 || none.Sessions == nil {
		t.Fatalf("no match should be an empty list, got %+v", none)
	}
}

func TestScheduleToolUpstreamError(t *testing.T) {
	cs := connect(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusBadGateway)
	})
	res, _ := callSchedule(t, cs, map[string]any{})
	if !res.IsError {
		t.Fatal("upstream failure should surface as a tool error the model can read")
	}
}
