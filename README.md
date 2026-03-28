# MCP Streamable HTTP Explorer

This is an MCP (Model Context Protocol) server built in Go using the **Streamable HTTP** transport. 

It provides tools and logs all HTTP request/response details (headers, payload, and streamed data) to stdout, allowing for deep inspection of the protocol's transport layer.

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

## Testing with MCP Inspector (Recommended)

The [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is the easiest way to test the server with a UI.

1. Ensure the server is running (`go run main.go`).
2. Open a new terminal and run:
   ```bash
   npx @modelcontextprotocol/inspector
   ```
3. Open the URL provided in the console (it will include an auth token).
4. In the "Connect" panel:
   - Choose **Streamable HTTP** as the transport.
   - Enter `http://localhost:8080/mcp` as the URL.
   - Click **Connect**.
5. You can now browse tools and call them via the UI.

## Manual Testing (curl)

### 1. Initiate SSE (GET)
To start a session, use `curl` to open an SSE stream:
```bash
curl -v http://localhost:8080/mcp
```
The server will respond with an `Mcp-Session-Id` header. Keep this ID for subsequent requests.

### 2. Initialize the Session (POST)
Before calling tools, you must initialize the session:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Mcp-Session-Id: <YOUR-SESSION-ID>" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {"name": "curl", "version": "1.0"}
    },
    "id": 1
  }'
```

### 3. Call a Tool (POST)
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
      "arguments": { "a": 10.5, "b": 20.2 }
    },
    "id": 2
  }'
```

Test the `long_running_task` with progress notifications:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Mcp-Session-Id: <YOUR-SESSION-ID>" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "long_running_task",
      "arguments": { "duration": 5 },
      "_meta": { "progressToken": "my-token" }
    },
    "id": 3
  }'
```

The responses and progress updates will be streamed back through the SSE connection initiated in Step 1.
