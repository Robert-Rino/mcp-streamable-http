# MCP Streamable HTTP Server Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an MCP server in Go that uses the Streamable HTTP transport and provides an `add` tool, with detailed HTTP/SSE logging for protocol exploration.

**Architecture:** A standalone Go web server using the official MCP Go SDK. It includes a logging middleware to intercept and print all JSON-RPC and SSE traffic.

**Tech Stack:** Go 1.22+, `github.com/modelcontextprotocol/go-sdk`

---

### Task 1: Project Initialization

**Files:**
- Create: `go.mod`

- [ ] **Step 1: Initialize Go module**
Run: `go mod init mcp-streamable-http`

- [ ] **Step 2: Add dependencies**
Run: `go get github.com/modelcontextprotocol/go-sdk/mcp`

- [ ] **Step 3: Commit**
```bash
git add go.mod go.sum
git commit -m "chore: initialize go module"
```

---

### Task 2: Core MCP Server & Add Tool

**Files:**
- Create: `main.go`

- [ ] **Step 1: Implement basic server structure and `add` tool**
Write `main.go` with `mcp.NewServer` and `mcp.AddTool`.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
				&mcp.TextContent{Text: fmt.Sprintf("The result is %.2f", result)},
			},
		}, nil, nil
	})

    // Placeholder for transport
    log.Println("Server initialized")
}
```

- [ ] **Step 2: Verify compilation**
Run: `go build -o server .`

- [ ] **Step 3: Commit**
```bash
git add main.go
git commit -m "feat: implement basic mcp server with add tool"
```

---

### Task 3: Streamable HTTP Transport & Logging Middleware

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Implement logging middleware**
Add a wrapper for `http.Handler` that logs request details and intercepts SSE responses.

- [ ] **Step 2: Set up Streamable HTTP Handler**
Use `mcp.NewStreamableHTTPHandler` and attach it to the logging middleware.

```go
// In main.go (refined)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--> %s %s", r.Method, r.URL.Path)
		for k, v := range r.Header {
			log.Printf("  Header %s: %v", k, v)
		}
        // For POST, log body (simplified for now)
		next.ServeHTTP(w, r)
	})
}
```

- [ ] **Step 3: Start the server**
Implement `http.ListenAndServe(":8080", mux)`.

- [ ] **Step 4: Commit**
```bash
git add main.go
git commit -m "feat: add streamable http transport and logging middleware"
```

---

### Task 4: Documentation & Usage Instructions

**Files:**
- Create: `README.md`

- [ ] **Step 1: Write README.md**
Include instructions for running the server and testing it with `curl`.

```markdown
# MCP Streamable HTTP Explorer

## Setup
1. `go mod tidy`
2. `go run main.go`

## Testing
1. Initiate SSE: `curl -v http://localhost:8080/mcp`
2. Call tool (requires session ID from Step 1): ...
```

- [ ] **Step 2: Commit**
```bash
git add README.md
git commit -m "docs: add readme with usage instructions"
```

---

### Task 5: Final Verification

- [ ] **Step 1: Run the server**
Run: `go run main.go`

- [ ] **Step 2: Test basic connectivity**
Use `curl -I http://localhost:8080/mcp` to check if it responds.

- [ ] **Step 3: (Optional) Run a test client script**
Create a small Go client if needed to verify the full flow.
