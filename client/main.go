package main

import (
	"context"
	"log"
	"time"

	"grpc_app/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}