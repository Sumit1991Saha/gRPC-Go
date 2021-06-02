package main

import (
	context "context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/saha/grpc-go-course/greet"
	"github.com/saha/grpc-go-course/greet/greetpb"
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

func (*server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function invoked with %v \n", request)
	greetRequest := request.GetGreeting()
	if greetRequest == nil {
		return nil
	}
	firstName := greetRequest.GetFirstName()
	lastName := greetRequest.GetLastName()

	for i := 0; i < 10; i++ {
		response := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " " + lastName + " " + strconv.Itoa(i),
		}
		_ = stream.Send(response)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("Starting gRPC Server")

	lis, err := net.Listen(greet.Protocol, greet.Host)

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
