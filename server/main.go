package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc_app/proto"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
}

type Message struct {
	Name string `json:"name"`
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Connect to gRPC server
	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return
	}
	defer grpcConn.Close()

	client := proto.NewGreeterClient(grpcConn)

	for {
		// Read message from WebSocket
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		// Parse the message
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Call gRPC service
		ctx := context.Background()
		resp, err := client.SayHello(ctx, &proto.HelloRequest{Name: message.Name})
		if err != nil {
			log.Printf("Error calling gRPC service: %v", err)
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

func main() {
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})

	log.Printf("gRPC Server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Set up HTTP server with WebSocket handler
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/ws", HandleWebsocket)

	log.Printf("Web Server starting at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
