package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ScheduleInput is the single live-demo tool. All fields are optional.
type ScheduleInput struct {
	Day         string `json:"day,omitempty" jsonschema:"optional day filter: Friday, Saturday, Viernes, or Sabado"`
	Time        string `json:"time,omitempty" jsonschema:"optional time filter, e.g. 18:05, 6:05 PM"`
	Person      string `json:"person,omitempty" jsonschema:"optional speaker name substring, e.g. Alberto"`
	Description string `json:"description,omitempty" jsonschema:"optional text searched in title and description"`
}

type ScheduleToolResult struct {
	Event    string        `json:"event"`
	Location string        `json:"location"`
	Source   string        `json:"source"`
	Timezone string        `json:"timezone"`
	Query    ScheduleInput `json:"query"`
	Count    int           `json:"count"`
	Sessions []Session     `json:"sessions"`
}

func registerScheduleTool(server *mcp.Server, cfg runtimeConfig) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "schedule",
		Description: "Fetch the live DevFest Valencia agenda and filter sessions. " +
			"Optional filters: day, time, person (speaker), description (matched against title and description).",
	}, makeScheduleHandler(cfg))
}

func makeScheduleHandler(cfg runtimeConfig) func(context.Context, *mcp.CallToolRequest, ScheduleInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in ScheduleInput) (*mcp.CallToolResult, any, error) {
		started := time.Now()
		log.Printf("tool schedule: call day=%q time=%q person=%q description=%q slug=%q event_id=%d",
			in.Day, in.Time, in.Person, in.Description, cfg.Slug, cfg.EventID)

		payload, err := FetchEvent(cfg.Slug, cfg.EventID)
		if err != nil {
			log.Printf("tool schedule: fetch error after %s: %v",
				time.Since(started).Round(time.Millisecond), err)
			return toolError(fmt.Sprintf("live fetch failed: %v", err)), nil, nil
		}
		schedule, err := ParseSchedule(payload)
		if err != nil {
			log.Printf("tool schedule: parse error after %s: %v",
				time.Since(started).Round(time.Millisecond), err)
			return toolError(fmt.Sprintf("parse failed: %v", err)), nil, nil
		}

		matched := schedule.Filter(ScheduleQuery{
			Day:         in.Day,
			Time:        in.Time,
			Person:      in.Person,
			Description: in.Description,
		})

		if matched == nil {
			matched = []Session{} // encode "no matches" as [] rather than null
		}
		out := ScheduleToolResult{
			Event:    schedule.Event,
			Location: schedule.Location,
			Source:   schedule.Source,
			Timezone: schedule.Timezone,
			Query:    in,
			Count:    len(matched),
			Sessions: matched,
		}
		log.Printf("tool schedule: done count=%d elapsed=%s",
			out.Count, time.Since(started).Round(time.Millisecond))
		return textJSON(out)
	}
}

func textJSON(v any) (*mcp.CallToolResult, any, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, v, nil
}

func toolError(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}
