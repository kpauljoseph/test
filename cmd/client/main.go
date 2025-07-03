package main

import (
	"context"
	"log"
	"time"

	proto "github.com/kpauljoseph/test/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	address = "localhost:50051"
)

func main() {
	log.Println("Starting gRPC blog client...")

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connected to server at %s", address)

	log.Println("\n1. Creating blog posts...")
	post1 := createPost(ctx, client, "Go Programming Best Practices", "Learn the best practices for writing Go code...", "John Doe", []string{"golang", "programming", "best-practices"})
	post2 := createPost(ctx, client, "gRPC Tutorial", "A comprehensive guide to gRPC in Go...", "Jane Smith", []string{"grpc", "golang", "tutorial"})

	log.Println("\n2. Reading blog posts...")
	if post1 != nil {
		readPost(ctx, client, post1.PostId)
	}
	if post2 != nil {
		readPost(ctx, client, post2.PostId)
	}

	log.Println("\n3. Reading non-existent post...")
	readPost(ctx, client, "non-existent-id")

	log.Println("\n4. Updating blog post...")
	if post1 != nil {
		updatePost(ctx, client, post1.PostId, "Go Programming Advanced Techniques", "Advanced techniques for Go development...", "John Doe Updated", []string{"golang", "advanced", "techniques"})
	}

	log.Println("\n5. Deleting blog post...")
	if post2 != nil {
		deletePost(ctx, client, post2.PostId)
	}

	log.Println("\n6. Verifying deletion...")
	if post2 != nil {
		readPost(ctx, client, post2.PostId)
	}
}

func createPost(ctx context.Context, client proto.BlogServiceClient, title, content, author string, tags []string) *proto.BlogPost {
	log.Printf("Creating post: title='%s', author='%s'", title, author)

	req := &proto.CreatePostRequest{
		Title:           title,
		Content:         content,
		Author:          author,
		PublicationDate: timestamppb.New(time.Now()),
		Tags:            tags,
	}

	resp, err := client.CreatePost(ctx, req)
	if err != nil {
		log.Printf("CreatePost failed: %v", err)
		return nil
	}

	if resp.Error != "" {
		log.Printf("CreatePost error: %s", resp.Error)
		return nil
	}

	post := resp.Post
	log.Printf("Post created successfully!")
	log.Printf("  PostID: %s", post.PostId)
	log.Printf("  Title: %s", post.Title)
	log.Printf("  Content: %s", post.Content)
	log.Printf("  Author: %s", post.Author)
	log.Printf("  Publication Date: %s", post.PublicationDate.AsTime().Format(time.RFC3339))
	log.Printf("  Tags: %v", post.Tags)

	return post
}

func readPost(ctx context.Context, client proto.BlogServiceClient, postID string) {
	log.Printf("Reading post: postID='%s'", postID)

	req := &proto.ReadPostRequest{
		PostId: postID,
	}

	resp, err := client.ReadPost(ctx, req)
	if err != nil {
		log.Printf("ReadPost failed: %v", err)
		return
	}

	if resp.Error != "" {
		log.Printf("ReadPost error: %s", resp.Error)
		return
	}

	post := resp.Post
	log.Printf("Post found!")
	log.Printf("  PostID: %s", post.PostId)
	log.Printf("  Title: %s", post.Title)
	log.Printf("  Content: %s", post.Content)
	log.Printf("  Author: %s", post.Author)
	log.Printf("  Publication Date: %s", post.PublicationDate.AsTime().Format(time.RFC3339))
	log.Printf("  Tags: %v", post.Tags)
}

func updatePost(ctx context.Context, client proto.BlogServiceClient, postID, title, content, author string, tags []string) {
	log.Printf("Updating post: postID='%s'", postID)

	req := &proto.UpdatePostRequest{
		PostId:  postID,
		Title:   title,
		Content: content,
		Author:  author,
		Tags:    tags,
	}

	resp, err := client.UpdatePost(ctx, req)
	if err != nil {
		log.Printf("UpdatePost failed: %v", err)
		return
	}

	if resp.Error != "" {
		log.Printf("UpdatePost error: %s", resp.Error)
		return
	}

	post := resp.Post
	log.Printf("Post updated successfully!")
	log.Printf("  PostID: %s", post.PostId)
	log.Printf("  Title: %s", post.Title)
	log.Printf("  Content: %s", post.Content)
	log.Printf("  Author: %s", post.Author)
	log.Printf("  Publication Date: %s", post.PublicationDate.AsTime().Format(time.RFC3339))
	log.Printf("  Tags: %v", post.Tags)
}

func deletePost(ctx context.Context, client proto.BlogServiceClient, postID string) {
	log.Printf("Deleting post: postID='%s'", postID)

	req := &proto.DeletePostRequest{
		PostId: postID,
	}

	resp, err := client.DeletePost(ctx, req)
	if err != nil {
		log.Printf("DeletePost failed: %v", err)
		return
	}

	if resp.Error != "" {
		log.Printf("DeletePost error: %s", resp.Error)
		return
	}

	if resp.Success {
		log.Printf("Post deleted successfully!")
	} else {
		log.Printf("Failed to delete post")
	}
}
