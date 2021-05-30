package main

import (
	context "context"
	"fmt"
	"github.com/saha/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function invoked with %v \n", request)
	greetRequest := request.GetGreeting()
	if greetRequest == nil {
		return nil, nil
	}
	firstName := greetRequest.GetFirstName()
	lastName := greetRequest.GetLastName()
	greetResponse := &greetpb.GreetResponse{
		Result: "Hello " + firstName + " " + lastName,
	}
	return greetResponse, nil
}

func main() {
	fmt.Println("Starting gRPC Server")

	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatal("Failed to listen :- ", err)
	}

	grpcServer := grpc.NewServer()
	localServer := server{}

	greetpb.RegisterGreetServiceServer(grpcServer, &localServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server :- ", err)
	}
}
