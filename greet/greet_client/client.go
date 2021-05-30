package main

import (
	"context"
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

	doUnary(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sumit",
			LastName: "Saha",
		},
	}
	response, err := client.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling greet rpc : %v", err)
	}
	log.Printf("Response from greet : %v", response.Result)
}
