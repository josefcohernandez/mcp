package main

import (
	"context"
	"encoding/json"

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
		return textResult(schedule.Filter(in.Day, in.Track))
	}
}

func makeGetSession(schedule *Schedule) func(context.Context, *mcp.CallToolRequest, GetSessionInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in GetSessionInput) (*mcp.CallToolResult, any, error) {
		sess, err := schedule.ByID(in.ID)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
				IsError: true,
			}, nil, nil
		}
		return textResult(sess)
	}
}

func makeFindSpeaker(schedule *Schedule) func(context.Context, *mcp.CallToolRequest, FindSpeakerInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in FindSpeakerInput) (*mcp.CallToolResult, any, error) {
		return textResult(schedule.BySpeaker(in.Name))
	}
}
