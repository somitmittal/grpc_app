package main

import (
	"context"
	"fmt"
	"grpc_app/proto"
	"log"
	"time"

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

	for {
		var name string
		log.Print("Enter your name (or 'quit' to exit): ")
		_, err := fmt.Scanln(&name)
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}

		if name == "quit" {
			log.Println("Goodbye!")
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		r, err := c.SayHello(ctx, &proto.HelloRequest{Name: name})
		cancel()

		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}
}
