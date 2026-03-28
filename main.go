package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AddArgs struct {
	A float64 `json:"a" jsonschema:"The first number"`
	B float64 `json:"b" jsonschema:"The second number"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--> %s %s", r.Method, r.URL.Path)
		for k, v := range r.Header {
			log.Printf("  Header %s: %v", k, v)
		}
		next.ServeHTTP(w, r)
	})
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

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", loggingMiddleware(handler))

	log.Println("MCP Server starting on :8080/mcp")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

