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

	blog := createBlog(client)

	readBlog(client, blog)
	updateBlog(client, blog)
}

func createBlog(client blogpb.BlogServiceClient) *blogpb.Blog{
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
	log.Printf("Blog has been created : %v \n", response)
	return response.Blog
}

func readBlog(client blogpb.BlogServiceClient, blog *blogpb.Blog) {
	// read Blog
	log.Println("Reading the blog")

	//Reading random blog which is not there in db
	_, err := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5bdc29e661b75adcac496cf4"})
	if err != nil {
		log.Printf("Error happened while reading: %v \n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blog.Id}
	readBlogRes, err := client.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		log.Printf("Error happened while reading: %v \n", err)
	}

	log.Printf("Blog was read: %v \n", readBlogRes)
}

func updateBlog(client blogpb.BlogServiceClient, blog *blogpb.Blog) {
	// update Blog
	newBlog := &blogpb.Blog{
		Id:       blog.GetId(),
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := client.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)
}
