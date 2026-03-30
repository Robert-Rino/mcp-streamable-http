package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var verboseLogging = true

func init() {
	if val, ok := os.LookupEnv("MCP_VERBOSE_LOGGING"); ok {
		verboseLogging = strings.ToLower(val) == "true"
	}
}

type AddArgs struct {
	A float64 `json:"a" jsonschema:"The first number"`
	B float64 `json:"b" jsonschema:"The second number"`
}

type LongRunningArgs struct {
	Duration int `json:"duration" jsonschema:"Duration of the task in seconds"`
}

type responseLogger struct {
	http.ResponseWriter
	wroteHeader bool
}

func (rl *responseLogger) logHeaders(statusCode int) {
	if rl.wroteHeader {
		return
	}
	rl.wroteHeader = true
	if verboseLogging {
		log.Printf("<-- RESPONSE %d", statusCode)
		for k, v := range rl.ResponseWriter.Header() {
			log.Printf("  Response Header %s: %v", k, v)
		}
	}
}

func (rl *responseLogger) WriteHeader(statusCode int) {
	rl.logHeaders(statusCode)
	rl.ResponseWriter.WriteHeader(statusCode)
}

func (rl *responseLogger) Write(b []byte) (int, error) {
	if !rl.wroteHeader {
		rl.logHeaders(http.StatusOK)
	}
	if verboseLogging {
		log.Printf("<-- RESPONSE DATA: %s", string(b))
	}
	return rl.ResponseWriter.Write(b)
}

func (rl *responseLogger) Flush() {
	if f, ok := rl.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if verboseLogging {
			log.Printf("--> %s %s", r.Method, r.URL.Path)
			for k, v := range r.Header {
				log.Printf("  Header %s: %v", k, v)
			}
		}

		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err == nil {
				if len(body) > 0 && verboseLogging {
					log.Printf("--> REQUEST PAYLOAD: %s", string(body))
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		next.ServeHTTP(&responseLogger{ResponseWriter: w}, r)
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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "long_running_task",
		Description: "Simulates a long-running task with progress notifications",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args LongRunningArgs) (*mcp.CallToolResult, any, error) {
		log.Printf("[TOOL] long_running_task: starting for %d seconds", args.Duration)

		// Access the session from the request
		session := req.GetSession().(*mcp.ServerSession)

		// 1. Send an info log notification
		session.Log(ctx, &mcp.LoggingMessageParams{
			Level: "info",
			Data:  "Starting the long-running simulation...",
		})

		// 2. Simulate progress
		steps := 4
		for i := 1; i <= steps; i++ {
			time.Sleep(time.Duration(args.Duration*1000/steps) * time.Millisecond)
			progress := float64(i) / float64(steps) * 100
			log.Printf("[PROGRESS] %.0f%%", progress)

			// Notify progress if a token was provided in Meta
			if req.Params.Meta != nil {
				if token, ok := req.Params.Meta["progressToken"]; ok {
					session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
						Progress:      progress,
						Total:         100,
						ProgressToken: token,
					})
				}
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Task completed successfully after %d seconds", args.Duration),
				},
			},
		}, nil, nil
	})

	log.Println("Server initialized with tools")

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

