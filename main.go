package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AddArgs struct {
	A float64 `json:"a" jsonschema:"description=The first number"`
	B float64 `json:"b" jsonschema:"description=The second number"`
}

func main() {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "explorer-server",
			Version: "1.0.0",
		},
		nil,
	)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "Adds two numbers together",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddArgs) (*mcp.CallToolResult, any, error) {
		result := args.A + args.B
		log.Printf("[TOOL] add: %.2f + %.2f = %.2f", args.A, args.B, result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("The result is %.2f", result),
				},
			},
		}, nil, nil
	})

	log.Println("Server initialized with add tool")
	// Transport to be added in next task
}
