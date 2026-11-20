package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GreetInput is the structured input for the greet tool.
type GreetInput struct {
	Name string `json:"name" jsonschema:"the person to greet"`
}

// greet is the tool handler. The SDK has already unmarshalled and validated the input.
func greet(_ context.Context, _ *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, any, error) {
	msg := fmt.Sprintf("Hello, %s! Welcome to DevFest Valencia.", input.Name)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
	}, nil, nil
}

func main() {
	// 1. Create the server: the Implementation is what the client sees during "initialize".
	server := mcp.NewServer(&mcp.Implementation{Name: "devfest-mcp", Version: "0.1.0"}, nil)

	// 2. Register the tool. AddTool infers the JSON input schema from GreetInput
	//    (json / jsonschema struct tags) and validates arguments before calling greet.
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "Greet a DevFest attendee by name"}, greet)

	// 3. Serve over stdio (newline-delimited JSON-RPC on stdin/stdout) until the client disconnects.
	//    Never write to stdout yourself: it is the protocol channel. Logs go to stderr.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
