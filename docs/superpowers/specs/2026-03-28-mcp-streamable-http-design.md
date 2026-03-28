# Design: MCP Streamable HTTP Server in Go

## 1. Goal
Create a Model Context Protocol (MCP) server in Go that uses the **Streamable HTTP** transport. The goal is to explore how the protocol works by inspecting HTTP requests/responses and testing a simple `add` tool.

## 2. Architecture
The server will be a standalone Go web application using the official `github.com/modelcontextprotocol/go-sdk`.

### Transport: Streamable HTTP
- **Endpoint:** `/mcp`
- **GET:** Initiates a Server-Sent Events (SSE) stream.
- **POST:** Receives JSON-RPC messages from the client.
- **DELETE:** Terminates the session.
- **Session Management:** Uses `Mcp-Session-Id` header to track client sessions.

### Tool: `add`
- **Name:** `add`
- **Description:** Adds two numbers together.
- **Arguments:**
  - `a` (float64): The first number.
  - `b` (float64): The second number.
- **Response:** A text message containing the result.

## 3. Implementation Details

### Technology Stack
- **Language:** Go 1.22+
- **SDK:** `github.com/modelcontextprotocol/go-sdk`

### Logging Strategy
To facilitate inspection, we will implement a custom `http.Handler` wrapper that logs:
1.  **Request Method & URL**
2.  **Request Headers** (especially `Mcp-Session-Id` and `Content-Type`)
3.  **Request Body** (JSON-RPC payload)
4.  **Response Status & Headers**
5.  **SSE Events** (Intercepting the stream output)

### Project Structure
```
.
├── go.mod
├── main.go
└── README.md
```

## 4. Usage Flow
1.  **Initialize:** `go mod init mcp-streamable-http` and `go mod tidy`.
2.  **Run Server:** `go run main.go`.
3.  **Test with Client:** Use an MCP client (e.g., Claude Desktop or a test script) pointing to `http://localhost:8080/mcp`.
4.  **Inspect Logs:** Observe the HTTP handshake and message flow in the server terminal.

## 5. Success Criteria
- The server starts and listens on `:8080`.
- The `add` tool is correctly registered and discoverable by an MCP client.
- The `add` tool returns the correct sum.
- All HTTP traffic is clearly logged to stdout.
