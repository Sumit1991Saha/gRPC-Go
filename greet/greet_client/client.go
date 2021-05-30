package main

import (
	"fmt"
	"log"

	"github.com/saha/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting gRPC Client")

	clientConnection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer clientConnection.Close()

	client := greetpb.NewGreetServiceClient(clientConnection)
	log.Println(client)

}
