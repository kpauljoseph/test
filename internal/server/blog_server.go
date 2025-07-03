package server

import (
	"context"
	"log"

	"github.com/kpauljoseph/test/internal/storage"
	proto "github.com/kpauljoseph/test/proto"
)

type BlogServer struct {
	proto.UnimplementedBlogServiceServer
	storage *storage.MemoryStorage
}

func NewBlogServer(storage *storage.MemoryStorage) *BlogServer {
	return &BlogServer{
		storage: storage,
	}
}

func (s *BlogServer) CreatePost(ctx context.Context, req *proto.CreatePostRequest) (*proto.CreatePostResponse, error) {
	log.Printf("Creating post: title=%s, author=%s", req.Title, req.Author)

	if req.Title == "" {
		return &proto.CreatePostResponse{
			Error: "title is required",
		}, nil
	}
	if req.Content == "" {
		return &proto.CreatePostResponse{
			Error: "content is required",
		}, nil
	}
	if req.Author == "" {
		return &proto.CreatePostResponse{
			Error: "author is required",
		}, nil
	}
	if req.PublicationDate == nil {
		return &proto.CreatePostResponse{
			Error: "publication_date is required",
		}, nil
	}

	post, err := s.storage.CreatePost(req.Title, req.Content, req.Author, req.PublicationDate, req.Tags)
	if err != nil {
		log.Printf("Failed to create post: %v", err)
		return &proto.CreatePostResponse{
			Error: err.Error(),
		}, nil
	}

	log.Printf("Post created successfully: postId=%s", post.PostId)
	return &proto.CreatePostResponse{
		Post: post,
	}, nil
}

func (s *BlogServer) ReadPost(ctx context.Context, req *proto.ReadPostRequest) (*proto.ReadPostResponse, error) {
	log.Printf("Reading post: postId=%s", req.PostId)

	if req.PostId == "" {
		return &proto.ReadPostResponse{
			Error: "post_id is required",
		}, nil
	}

	post, err := s.storage.GetPost(req.PostId)
	if err != nil {
		log.Printf("Post not found: postId=%s, error=%v", req.PostId, err)
		return &proto.ReadPostResponse{
			Error: err.Error(),
		}, nil
	}

	log.Printf("Post found: postId=%s, title=%s", post.PostId, post.Title)
	return &proto.ReadPostResponse{
		Post: post,
	}, nil
}

func (s *BlogServer) UpdatePost(ctx context.Context, req *proto.UpdatePostRequest) (*proto.UpdatePostResponse, error) {
	log.Printf("Updating post: postId=%s", req.PostId)

	if req.PostId == "" {
		return &proto.UpdatePostResponse{
			Error: "post_id is required",
		}, nil
	}
	if req.Title == "" {
		return &proto.UpdatePostResponse{
			Error: "title is required",
		}, nil
	}
	if req.Content == "" {
		return &proto.UpdatePostResponse{
			Error: "content is required",
		}, nil
	}
	if req.Author == "" {
		return &proto.UpdatePostResponse{
			Error: "author is required",
		}, nil
	}

	post, err := s.storage.UpdatePost(req.PostId, req.Title, req.Content, req.Author, req.Tags)
	if err != nil {
		log.Printf("Failed to update post: postId=%s, error=%v", req.PostId, err)
		return &proto.UpdatePostResponse{
			Error: err.Error(),
		}, nil
	}

	log.Printf("Post updated successfully: postId=%s", post.PostId)
	return &proto.UpdatePostResponse{
		Post: post,
	}, nil
}

func (s *BlogServer) DeletePost(ctx context.Context, req *proto.DeletePostRequest) (*proto.DeletePostResponse, error) {
	log.Printf("Deleting post: postId=%s", req.PostId)

	if req.PostId == "" {
		return &proto.DeletePostResponse{
			Success: false,
			Error:   "post_id is required",
		}, nil
	}

	err := s.storage.DeletePost(req.PostId)
	if err != nil {
		log.Printf("Failed to delete post: postId=%s, error=%v", req.PostId, err)
		return &proto.DeletePostResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	log.Printf("Post deleted successfully: postId=%s", req.PostId)
	return &proto.DeletePostResponse{
		Success: true,
	}, nil
}
