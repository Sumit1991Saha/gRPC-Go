package main

import (
	context "context"
	"fmt"
	"io"
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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Printf("LongGreet function invoked with streaming request \n")
	result := ""
	for {
		request, err := stream.Recv()
		log.Printf("Received request : %v", request)
		if err == io.EOF {
			response := &greetpb.LongGreetResponse{
				Result: result,
			}
			log.Printf("Sending response : %v", result)
			return stream.SendAndClose(response)
		} else if err != nil {
			log.Fatalf("Error while reading client Stream %v", err)
		}
		firstName := request.GetGreeting().GetFirstName()
		lastName := request.GetGreeting().GetLastName()
		result += "Hello " + firstName + " " + lastName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Printf("GreetEveryone function invoked with streaming request \n")

	for {
		request, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return nil
		}
		if recvErr != nil {
			log.Fatalf("Error while reading client Stream %v", recvErr)
		}
		log.Printf("Request received at server : %v \n", request)

		response := processRequestForBidirectionalStreams(request)
		sendErr := stream.Send(response)
		if sendErr != nil {
			log.Fatalf("Error while sending resposne to client %v", sendErr)
		}
		log.Printf("Response sent from server : %v \n\n", response)
	}
}

func processRequestForBidirectionalStreams(request *greetpb.GreetEveryoneRequest) *greetpb.GreetEveryoneResponse {
	log.Printf("Request being processed at server : %v \n", request)
	time.Sleep(5 * time.Second)
	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName + "! "
	response := &greetpb.GreetEveryoneResponse{
		Result: result,
	}
	return response
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
