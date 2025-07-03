package storage

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMemoryStorage_CreatePost(t *testing.T) {
	storage := NewMemoryStorage()

	tests := []struct {
		name    string
		title   string
		content string
		author  string
		tags    []string
	}{
		{
			name:    "create post with tags",
			title:   "Test Post",
			content: "This is test content",
			author:  "Test Author",
			tags:    []string{"tech", "golang"},
		},
		{
			name:    "create post without tags",
			title:   "Another Post",
			content: "More content",
			author:  "Another Author",
			tags:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubDate := timestamppb.New(time.Now())
			post, err := storage.CreatePost(tt.title, tt.content, tt.author, pubDate, tt.tags)
			if err != nil {
				t.Errorf("CreatePost() error = %v", err)
				return
			}

			if post.PostId == "" {
				t.Error("CreatePost() should generate a PostId")
			}
			if post.Title != tt.title {
				t.Errorf("CreatePost() title = %v, want %v", post.Title, tt.title)
			}
			if post.Content != tt.content {
				t.Errorf("CreatePost() content = %v, want %v", post.Content, tt.content)
			}
			if post.Author != tt.author {
				t.Errorf("CreatePost() author = %v, want %v", post.Author, tt.author)
			}
			if post.PublicationDate == nil {
				t.Error("CreatePost() should set PublicationDate")
			}
			if len(post.Tags) != len(tt.tags) {
				t.Errorf("CreatePost() tags length = %v, want %v", len(post.Tags), len(tt.tags))
			}
		})
	}
}

func TestMemoryStorage_GetPost(t *testing.T) {
	storage := NewMemoryStorage()

	post, err := storage.CreatePost("Test", "Content", "Author", timestamppb.New(time.Now()), []string{"tag"})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		postID  string
		wantErr bool
	}{
		{
			name:    "get existing post",
			postID:  post.PostId,
			wantErr: false,
		},
		{
			name:    "get non-existent post",
			postID:  "non-existent-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.GetPost(tt.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.PostId != tt.postID {
				t.Errorf("GetPost() PostId = %v, want %v", got.PostId, tt.postID)
			}
		})
	}
}

func TestMemoryStorage_UpdatePost(t *testing.T) {
	storage := NewMemoryStorage()

	post, err := storage.CreatePost("Original", "Original Content", "Original Author", timestamppb.New(time.Now()), []string{"original"})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name       string
		postID     string
		title      string
		content    string
		author     string
		tags       []string
		wantErr    bool
		wantTitle  string
		wantAuthor string
	}{
		{
			name:       "update existing post",
			postID:     post.PostId,
			title:      "Updated Title",
			content:    "Updated Content",
			author:     "Updated Author",
			tags:       []string{"updated", "test"},
			wantErr:    false,
			wantTitle:  "Updated Title",
			wantAuthor: "Updated Author",
		},
		{
			name:    "update non-existent post",
			postID:  "non-existent-id",
			title:   "Title",
			content: "Content",
			author:  "Author",
			tags:    []string{"tag"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.UpdatePost(tt.postID, tt.title, tt.content, tt.author, tt.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Title != tt.wantTitle {
					t.Errorf("UpdatePost() title = %v, want %v", got.Title, tt.wantTitle)
				}
				if got.Author != tt.wantAuthor {
					t.Errorf("UpdatePost() author = %v, want %v", got.Author, tt.wantAuthor)
				}
			}
		})
	}
}

func TestMemoryStorage_DeletePost(t *testing.T) {
	storage := NewMemoryStorage()

	post, err := storage.CreatePost("Test", "Content", "Author", timestamppb.New(time.Now()), []string{"tag"})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		postID  string
		wantErr bool
	}{
		{
			name:    "delete existing post",
			postID:  post.PostId,
			wantErr: false,
		},
		{
			name:    "delete non-existent post",
			postID:  "non-existent-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.DeletePost(tt.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := storage.GetPost(tt.postID)
				if err == nil {
					t.Error("DeletePost() should remove post from storage")
				}
			}
		})
	}
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	storage := NewMemoryStorage()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			post, err := storage.CreatePost("Title", "Content", "Author", timestamppb.New(time.Now()), []string{"tag"})
			if err != nil {
				t.Errorf("Concurrent CreatePost() failed: %v", err)
				return
			}

			_, err = storage.GetPost(post.PostId)
			if err != nil {
				t.Errorf("Concurrent GetPost() failed: %v", err)
			}

			_, err = storage.UpdatePost(post.PostId, "New Title", "New Content", "New Author", []string{"new"})
			if err != nil {
				t.Errorf("Concurrent UpdatePost() failed: %v", err)
			}

			err = storage.DeletePost(post.PostId)
			if err != nil {
				t.Errorf("Concurrent DeletePost() failed: %v", err)
			}
		}(i)
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out - possible deadlock")
		}
	}
}