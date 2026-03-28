# MCP Streamable HTTP Explorer

This is an MCP (Model Context Protocol) server built in Go using the **Streamable HTTP** transport. 

It provides an `add` tool and logs all HTTP request headers and methods to stdout, allowing for inspection of the protocol's transport layer.

## Setup

1. Install Go 1.22+.
2. Initialize dependencies:
   ```bash
   go mod tidy
   ```

## Running the Server

Start the server on port 8080:
```bash
go run main.go
```

The server will listen on `http://localhost:8080/mcp`.

## Testing the Protocol

### 1. Initiate SSE (GET)
To start a session, use `curl` to open an SSE stream:
```bash
curl -v http://localhost:8080/mcp
```
The server will respond with an `Mcp-Session-Id` header. Keep this ID for subsequent requests.

### 2. List Tools (POST)
In a separate terminal, use the session ID to list available tools:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Mcp-Session-Id: <YOUR-SESSION-ID>" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "params": {},
    "id": 1
  }'
```

### 3. Call Tool (POST)
Call the `add` tool:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Mcp-Session-Id: <YOUR-SESSION-ID>" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {
        "a": 10.5,
        "b": 20.2
      }
    },
    "id": 2
  }'
```

The response will be streamed back through the SSE connection initiated in Step 1.
