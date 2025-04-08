# gRPC Chat Application

A real-time chat application demonstrating the integration of gRPC and WebSockets in Go. This application showcases how to build a modern client-server architecture using gRPC for backend communication and WebSockets for real-time web frontend interaction.

## Architecture

This application follows a three-tier architecture:

1. **gRPC Server**: Core backend service implementing the greeting functionality
2. **WebSocket Proxy**: Middleware that translates between WebSocket and gRPC protocols
3. **Clients**:
   - Web Client: Browser-based UI using WebSockets
   - CLI Client: Command-line interface using direct gRPC calls

```
+----------------+     +-----------------+     +----------------+
|                |     |                 |     |                |
|  Web Client    |<--->|  WebSocket     |<--->|  gRPC Server   |
|  (Browser)     |     |  Proxy         |     |                |
|                |     |                 |     |                |
+----------------+     +-----------------+     +----------------+
                                               ^
                                               |
                       +----------------+      |
                       |                |      |
                       |  CLI Client    |<-----+
                       |                |      
                       |                |      
                       +----------------+      
```

## Features

- Real-time communication between clients and server
- Protocol Buffer-based message serialization
- WebSocket support for browser clients
- Command-line interface for direct gRPC interaction
- Simple and clean web UI

## Prerequisites

- Go 1.16 or higher
- Web browser (for web client)

## Installation

1. Clone the repository

```bash
git clone <repository-url>
cd grpc_app
```

2. Install dependencies

```bash
go mod download
```

## Usage

### Starting the Server

Run the server which hosts both the gRPC service and WebSocket proxy:

```bash
cd server
go run .
```

This will start:
- gRPC server on port 50051
- Web/WebSocket server on port 8080

### Using the Web Client

1. Open your browser and navigate to `http://localhost:8080`
2. Enter your name in the input field and click "Send"
3. You'll receive a greeting message from the server

### Using the CLI Client

In a separate terminal:

```bash
cd client
go run .
```

Follow the prompts to enter your name and receive greetings from the server.

## Project Structure

```
├── client/             # CLI client implementation
│   └── main.go
├── proto/              # Protocol Buffer definitions
│   ├── hello.proto     # Service definition
│   ├── hello.pb.go     # Generated code
│   └── hello_grpc.pb.go
├── server/             # Server implementation
│   ├── main.go         # gRPC server and HTTP handlers
│   └── websocket.go    # WebSocket implementation
└── web/                # Web client
    ├── index.html      # Web UI
    └── chat.js         # WebSocket client
```

## How It Works

1. The gRPC server implements the `Greeter` service defined in `hello.proto`
2. The WebSocket proxy:
   - Accepts WebSocket connections from browsers
   - Translates WebSocket messages to gRPC calls
   - Returns gRPC responses back through the WebSocket
3. The web client connects to the WebSocket endpoint and displays messages
4. The CLI client directly connects to the gRPC server

## Development

### Regenerating Protocol Buffers

If you modify the `.proto` files, regenerate the Go code with:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/hello.proto
```

## License

[MIT](LICENSE)
