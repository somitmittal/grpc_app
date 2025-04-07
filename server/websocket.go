// Package main implements a WebSocket server that connects to a gRPC backend
package main

import (
	"context"
	"encoding/json"
	"grpc_app/proto"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// HandleWebSocket handles WebSocket connections and communicates with the gRPC server
func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer func() {
		log.Printf("Client disconnected")
		conn.Close()
	}()

	// Send initial connection success message
	response := map[string]string{"message": "Connected to chat server"}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}

	// Connect to gRPC server
	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		conn.WriteJSON(map[string]string{"message": "Internal server error"})
		return
	}
	defer grpcConn.Close()

	client := proto.NewGreeterClient(grpcConn)

	for {
		// Read message from WebSocket
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading WebSocket message: %v", err)
			}
			break
		}

		// Parse the message
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			conn.WriteJSON(map[string]string{"message": "Invalid message format"})
			continue
		}

		// Call gRPC service
		ctx := context.Background()
		resp, err := client.SayHello(ctx, &proto.HelloRequest{Name: message.Name})
		if err != nil {
			log.Printf("Error calling gRPC service: %v", err)
			conn.WriteJSON(map[string]string{"message": "Failed to process message"})
			continue
		}

		// Send response back through WebSocket
		response := map[string]string{"message": resp.Message}
		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Error writing response: %v", err)
			break
		}
	}
}
