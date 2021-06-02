package main

import (
	"context"
	"fmt"
	"io"

	"github.com/saha/grpc-go-course/greet"

	"log"

	"google.golang.org/grpc"

	"github.com/saha/grpc-go-course/greet/greetpb"
)

func main() {
	fmt.Println("Starting gRPC Client")

	clientConnection, err := grpc.Dial(greet.Host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer clientConnection.Close()

	client := greetpb.NewGreetServiceClient(clientConnection)

	//doUnary(client)
	doServerStreaming(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sumit",
			LastName:  "Saha",
		},
	}
	response, err := client.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Greet rpc : %v", err)
	}
	log.Printf("Response from greet : %v", response.Result)
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	request := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sumit",
			LastName:  "Saha",
		},
	}
	responseStream, err := client.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes rpc : %v", err)
	}

	for {
		response, err := responseStream.Recv()
		if err == io.EOF {
			// End of stream has been reached
			log.Println("Stream finished")
			break
		} else if err != nil {
			log.Fatalf("Error while reading stream : %v", err)
		}
		log.Printf("Response from greet : %v", response.Result)
	}
}
