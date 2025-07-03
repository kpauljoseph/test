package server

import (
	"context"
	"testing"
	"time"

	"github.com/kpauljoseph/test/internal/storage"
	proto "github.com/kpauljoseph/test/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestBlogServer_CreatePost(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	server := NewBlogServer(memoryStorage)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *proto.CreatePostRequest
		wantErr bool
	}{
		{
			name: "valid post creation",
			req: &proto.CreatePostRequest{
				Title:           "Test Post",
				Content:         "This is test content",
				Author:          "Test Author",
				PublicationDate: timestamppb.New(time.Now()),
				Tags:            []string{"test", "golang"},
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: &proto.CreatePostRequest{
				Content:         "This is test content",
				Author:          "Test Author",
				PublicationDate: timestamppb.New(time.Now()),
				Tags:            []string{"test"},
			},
			wantErr: true,
		},
		{
			name: "missing content",
			req: &proto.CreatePostRequest{
				Title:           "Test Post",
				Author:          "Test Author",
				PublicationDate: timestamppb.New(time.Now()),
				Tags:            []string{"test"},
			},
			wantErr: true,
		},
		{
			name: "missing author",
			req: &proto.CreatePostRequest{
				Title:           "Test Post",
				Content:         "This is test content",
				PublicationDate: timestamppb.New(time.Now()),
				Tags:            []string{"test"},
			},
			wantErr: true,
		},
		{
			name: "missing publication date",
			req: &proto.CreatePostRequest{
				Title:   "Test Post",
				Content: "This is test content",
				Author:  "Test Author",
				Tags:    []string{"test"},
			},
			wantErr: true,
		},
		{
			name: "valid post without tags",
			req: &proto.CreatePostRequest{
				Title:           "Test Post",
				Content:         "This is test content",
				Author:          "Test Author",
				PublicationDate: timestamppb.New(time.Now()),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.CreatePost(ctx, tt.req)
			if err != nil {
				t.Errorf("CreatePost() error = %v", err)
				return
			}

			if tt.wantErr {
				if resp.Error == "" {
					t.Error("CreatePost() expected error but got none")
				}
				if resp.Post != nil {
					t.Error("CreatePost() expected no post on error")
				}
			} else {
				if resp.Error != "" {
					t.Errorf("CreatePost() unexpected error: %s", resp.Error)
				}
				if resp.Post == nil {
					t.Error("CreatePost() expected post but got nil")
				} else {
					if resp.Post.PostId == "" {
						t.Error("CreatePost() post should have PostId")
					}
					if resp.Post.Title != tt.req.Title {
						t.Errorf("CreatePost() title = %v, want %v", resp.Post.Title, tt.req.Title)
					}
					if resp.Post.PublicationDate == nil {
						t.Error("CreatePost() post should have PublicationDate")
					}
				}
			}
		})
	}
}

func TestBlogServer_ReadPost(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	server := NewBlogServer(memoryStorage)
	ctx := context.Background()

	post, err := memoryStorage.CreatePost("Test Post", "Test Content", "Test Author", timestamppb.New(time.Now()), []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		req     *proto.ReadPostRequest
		wantErr bool
	}{
		{
			name: "read existing post",
			req: &proto.ReadPostRequest{
				PostId: post.PostId,
			},
			wantErr: false,
		},
		{
			name: "read non-existent post",
			req: &proto.ReadPostRequest{
				PostId: "non-existent-id",
			},
			wantErr: true,
		},
		{
			name: "empty post ID",
			req: &proto.ReadPostRequest{
				PostId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.ReadPost(ctx, tt.req)
			if err != nil {
				t.Errorf("ReadPost() error = %v", err)
				return
			}

			if tt.wantErr {
				if resp.Error == "" {
					t.Error("ReadPost() expected error but got none")
				}
				if resp.Post != nil {
					t.Error("ReadPost() expected no post on error")
				}
			} else {
				if resp.Error != "" {
					t.Errorf("ReadPost() unexpected error: %s", resp.Error)
				}
				if resp.Post == nil {
					t.Error("ReadPost() expected post but got nil")
				} else {
					if resp.Post.PostId != tt.req.PostId {
						t.Errorf("ReadPost() postId = %v, want %v", resp.Post.PostId, tt.req.PostId)
					}
				}
			}
		})
	}
}

func TestBlogServer_UpdatePost(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	server := NewBlogServer(memoryStorage)
	ctx := context.Background()

	post, err := memoryStorage.CreatePost("Original Title", "Original Content", "Original Author", timestamppb.New(time.Now()), []string{"original"})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		req     *proto.UpdatePostRequest
		wantErr bool
	}{
		{
			name: "valid update",
			req: &proto.UpdatePostRequest{
				PostId:  post.PostId,
				Title:   "Updated Title",
				Content: "Updated Content",
				Author:  "Updated Author",
				Tags:    []string{"updated", "test"},
			},
			wantErr: false,
		},
		{
			name: "update non-existent post",
			req: &proto.UpdatePostRequest{
				PostId:  "non-existent-id",
				Title:   "Title",
				Content: "Content",
				Author:  "Author",
				Tags:    []string{"tag"},
			},
			wantErr: true,
		},
		{
			name: "empty post ID",
			req: &proto.UpdatePostRequest{
				PostId:  "",
				Title:   "Title",
				Content: "Content",
				Author:  "Author",
			},
			wantErr: true,
		},
		{
			name: "missing title",
			req: &proto.UpdatePostRequest{
				PostId:  post.PostId,
				Content: "Content",
				Author:  "Author",
			},
			wantErr: true,
		},
		{
			name: "missing content",
			req: &proto.UpdatePostRequest{
				PostId: post.PostId,
				Title:  "Title",
				Author: "Author",
			},
			wantErr: true,
		},
		{
			name: "missing author",
			req: &proto.UpdatePostRequest{
				PostId:  post.PostId,
				Title:   "Title",
				Content: "Content",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.UpdatePost(ctx, tt.req)
			if err != nil {
				t.Errorf("UpdatePost() error = %v", err)
				return
			}

			if tt.wantErr {
				if resp.Error == "" {
					t.Error("UpdatePost() expected error but got none")
				}
				if resp.Post != nil {
					t.Error("UpdatePost() expected no post on error")
				}
			} else {
				if resp.Error != "" {
					t.Errorf("UpdatePost() unexpected error: %s", resp.Error)
				}
				if resp.Post == nil {
					t.Error("UpdatePost() expected post but got nil")
				} else {
					if resp.Post.Title != tt.req.Title {
						t.Errorf("UpdatePost() title = %v, want %v", resp.Post.Title, tt.req.Title)
					}
					if resp.Post.Author != tt.req.Author {
						t.Errorf("UpdatePost() author = %v, want %v", resp.Post.Author, tt.req.Author)
					}
				}
			}
		})
	}
}

func TestBlogServer_DeletePost(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	server := NewBlogServer(memoryStorage)
	ctx := context.Background()

	post, err := memoryStorage.CreatePost("Test Post", "Test Content", "Test Author", timestamppb.New(time.Now()), []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name        string
		req         *proto.DeletePostRequest
		wantErr     bool
		wantSuccess bool
	}{
		{
			name: "delete existing post",
			req: &proto.DeletePostRequest{
				PostId: post.PostId,
			},
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name: "delete non-existent post",
			req: &proto.DeletePostRequest{
				PostId: "non-existent-id",
			},
			wantErr:     true,
			wantSuccess: false,
		},
		{
			name: "empty post ID",
			req: &proto.DeletePostRequest{
				PostId: "",
			},
			wantErr:     true,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.DeletePost(ctx, tt.req)
			if err != nil {
				t.Errorf("DeletePost() error = %v", err)
				return
			}

			if tt.wantErr {
				if resp.Error == "" {
					t.Error("DeletePost() expected error but got none")
				}
				if resp.Success != tt.wantSuccess {
					t.Errorf("DeletePost() success = %v, want %v", resp.Success, tt.wantSuccess)
				}
			} else {
				if resp.Error != "" {
					t.Errorf("DeletePost() unexpected error: %s", resp.Error)
				}
				if resp.Success != tt.wantSuccess {
					t.Errorf("DeletePost() success = %v, want %v", resp.Success, tt.wantSuccess)
				}
			}
		})
	}
}

func TestBlogServer_Integration(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	server := NewBlogServer(memoryStorage)
	ctx := context.Background()

	t.Run("full CRUD workflow", func(t *testing.T) {
		createReq := &proto.CreatePostRequest{
			Title:           "Integration Test Post",
			Content:         "This is an integration test",
			Author:          "Test Author",
			PublicationDate: timestamppb.New(time.Now()),
			Tags:            []string{"integration", "test"},
		}
		createResp, err := server.CreatePost(ctx, createReq)
		if err != nil {
			t.Fatalf("CreatePost() error = %v", err)
		}
		if createResp.Error != "" {
			t.Fatalf("CreatePost() error = %s", createResp.Error)
		}
		if createResp.Post == nil {
			t.Fatal("CreatePost() returned nil post")
		}

		postId := createResp.Post.PostId

		readReq := &proto.ReadPostRequest{PostId: postId}
		readResp, err := server.ReadPost(ctx, readReq)
		if err != nil {
			t.Fatalf("ReadPost() error = %v", err)
		}
		if readResp.Error != "" {
			t.Fatalf("ReadPost() error = %s", readResp.Error)
		}
		if readResp.Post.Title != createReq.Title {
			t.Errorf("ReadPost() title = %v, want %v", readResp.Post.Title, createReq.Title)
		}

		updateReq := &proto.UpdatePostRequest{
			PostId:  postId,
			Title:   "Updated Integration Test Post",
			Content: "This is an updated integration test",
			Author:  "Updated Test Author",
			Tags:    []string{"updated", "integration", "test"},
		}
		updateResp, err := server.UpdatePost(ctx, updateReq)
		if err != nil {
			t.Fatalf("UpdatePost() error = %v", err)
		}
		if updateResp.Error != "" {
			t.Fatalf("UpdatePost() error = %s", updateResp.Error)
		}
		if updateResp.Post.Title != updateReq.Title {
			t.Errorf("UpdatePost() title = %v, want %v", updateResp.Post.Title, updateReq.Title)
		}

		deleteReq := &proto.DeletePostRequest{PostId: postId}
		deleteResp, err := server.DeletePost(ctx, deleteReq)
		if err != nil {
			t.Fatalf("DeletePost() error = %v", err)
		}
		if deleteResp.Error != "" {
			t.Fatalf("DeletePost() error = %s", deleteResp.Error)
		}
		if !deleteResp.Success {
			t.Error("DeletePost() should return success = true")
		}

		readAfterDelete, err := server.ReadPost(ctx, readReq)
		if err != nil {
			t.Fatalf("ReadPost() after delete error = %v", err)
		}
		if readAfterDelete.Error == "" {
			t.Error("ReadPost() after delete should return error")
		}
		if readAfterDelete.Post != nil {
			t.Error("ReadPost() after delete should return nil post")
		}
	})
}
