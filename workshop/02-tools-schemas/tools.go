package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListSessionsInput struct {
	Day   string `json:"day,omitempty" jsonschema:"optional day filter: Friday, Saturday, Viernes, or Sabado"`
	Track string `json:"track,omitempty" jsonschema:"optional track filter if present in the data"`
}

type GetSessionInput struct {
	ID string `json:"id" jsonschema:"session id, e.g. s5"`
}

type FindSpeakerInput struct {
	Name string `json:"name" jsonschema:"speaker name or substring, e.g. Alberto"`
}

func textResult(v any) (*mcp.CallToolResult, any, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, v, nil
}

func makeListSessions(schedule *Schedule) func(context.Context, *mcp.CallToolRequest, ListSessionsInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in ListSessionsInput) (*mcp.CallToolResult, any, error) {
		// TODO: return schedule.Filter(in.Day, in.Track) via textResult
		_ = in
		_ = schedule
		return nil, nil, fmt.Errorf("TODO: implement list_sessions")
	}
}

func makeGetSession(schedule *Schedule) func(context.Context, *mcp.CallToolRequest, GetSessionInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in GetSessionInput) (*mcp.CallToolResult, any, error) {
		// TODO: schedule.ByID(in.ID) then textResult
		_ = in
		_ = schedule
		return nil, nil, fmt.Errorf("TODO: implement get_session")
	}
}

func makeFindSpeaker(schedule *Schedule) func(context.Context, *mcp.CallToolRequest, FindSpeakerInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in FindSpeakerInput) (*mcp.CallToolResult, any, error) {
		// TODO: schedule.BySpeaker(in.Name) then textResult
		_ = in
		_ = schedule
		return nil, nil, fmt.Errorf("TODO: implement find_speaker")
	}
}
