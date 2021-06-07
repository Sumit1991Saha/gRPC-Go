package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

func (*server) CreateBlog(ctx context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Create Blog Request")
	blogData := request.GetBlog()
	data := &blogItem{
		AuthorID: blogData.GetAuthorId(),
		Content:  blogData.GetContent(),
		Title:    blogData.GetTitle(),
	}
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error : %v", err),
		)
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID: %v", err),
		)
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			AuthorId: blogData.GetAuthorId(),
			Content:  blogData.GetContent(),
			Title:    blogData.GetTitle(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Read Blog Request")

	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": oid}
	res := collection.FindOne(ctx, filter)
	if err = res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")
	blogData := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blogData.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)
	if err = res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	// we update our internal struct
	data.AuthorID = blogData.GetAuthorId()
	data.Content = blogData.GetContent()
	data.Title = blogData.GetTitle()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr),
		)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
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

func startServices() (*mongo.Client, *grpc.Server) {
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
