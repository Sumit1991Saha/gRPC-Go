package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/saha/grpc-go-course/blog"
	"github.com/saha/grpc-go-course/blog/blogpb"
)

var (
	collection *mongo.Collection
)
type server struct {
	blogpb.BlogServiceServer
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mongoClient, grpcServer := startServices()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	stopServices(mongoClient, grpcServer)
}

func startServices() ( *mongo.Client, *grpc.Server) {
	mongoClient := startMongoService()
	grpcServer := startGRPCService()
	return mongoClient, grpcServer
}

func startMongoService() *mongo.Client {
	// connect to MongoDB
	fmt.Println("Connecting to MongoDB")
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Blog Service Started")
	collection = mongoClient.Database("mydb").Collection("blog")
	return mongoClient
}

func startGRPCService() *grpc.Server {
	fmt.Println("Starting Blog gRPC Server")
	lis, err := net.Listen(blog.Protocol, blog.Host)
	if err != nil {
		log.Fatal("Failed to listen :- ", err)
	}
	grpcServer := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	go func(listener net.Listener) {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("Failed to server :- ", err)
		}
	}(lis)
	return grpcServer
}
func stopServices(mongoClient *mongo.Client, grpcServer *grpc.Server) {
	// First we close the connection with MongoDB:
	fmt.Print("\n Closing MongoDB Connection")
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error on disconnection with MongoDB : %v", err)
	}
	fmt.Println("\n Stopping Blog gRPC Server")
	grpcServer.Stop()
}

