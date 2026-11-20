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

// TODO: improve the greeting message if you want — the wiring below already works
// once you uncomment the server setup in main.
func greet(_ context.Context, _ *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, any, error) {
	msg := fmt.Sprintf("Hello, %s! Welcome to DevFest Valencia.", input.Name)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
	}, nil, nil
}

func main() {
	// TODO(1): create the server
	// server := mcp.NewServer(&mcp.Implementation{Name: "devfest-mcp", Version: "0.1.0"}, nil)

	// TODO(2): register the tool
	// mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "Greet a DevFest attendee by name"}, greet)

	// TODO(3): run over stdio until the client disconnects
	// if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
	// 	log.Fatal(err)
	// }

	_ = greet
	log.Fatal("TODO: uncomment the three steps in main() — see README.md")
}
