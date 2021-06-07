package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/saha/grpc-go-course/blog/blogpb"
	"github.com/saha/grpc-go-course/greet"
)

func main() {
	log.Println("Starting Blog Client")

	clientConnection, err := grpc.Dial(greet.Host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func(clientConnection *grpc.ClientConn) {
		err = clientConnection.Close()
		if err != nil {
		}
	}(clientConnection)

	client := blogpb.NewBlogServiceClient(clientConnection)

	log.Println("Creating a blog")
	blogData := &blogpb.Blog{
		AuthorId: "Sumit Saha",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}
	response, err := client.CreateBlog(context.Background(),
		&blogpb.CreateBlogRequest{
			Blog: blogData,
		},
	)
	if err != nil {
		log.Fatalf("Unexpected error : %v \n", err)
	}
	fmt.Printf("Blog has been created : %v \n", response)
}
