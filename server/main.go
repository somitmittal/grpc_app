package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"grpc_app/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}