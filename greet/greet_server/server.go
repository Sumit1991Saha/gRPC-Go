package main

import (
	"fmt"
	"log"
	"net"

	"github.com/saha/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
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
