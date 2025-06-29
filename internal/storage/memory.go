package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	proto "github.com/kpauljoseph/test/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MemoryStorage struct {
	mu    sync.RWMutex
	posts map[string]*proto.BlogPost
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts: make(map[string]*proto.BlogPost),
	}
}

func (s *MemoryStorage) CreatePost(title, content, author string, tags []string) (*proto.BlogPost, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := &proto.BlogPost{
		PostId:          uuid.New().String(),
		Title:           title,
		Content:         content,
		Author:          author,
		PublicationDate: timestamppb.New(time.Now()),
		Tags:            tags,
	}

	s.posts[post.PostId] = post
	return post, nil
}

func (s *MemoryStorage) GetPost(postID string) (*proto.BlogPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, exists := s.posts[postID]
	if !exists {
		return nil, fmt.Errorf("post with ID %s not found", postID)
	}

	return post, nil
}

func (s *MemoryStorage) UpdatePost(postID, title, content, author string, tags []string) (*proto.BlogPost, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[postID]
	if !exists {
		return nil, fmt.Errorf("post with ID %s not found", postID)
	}

	post.Title = title
	post.Content = content
	post.Author = author
	post.Tags = tags

	return post, nil
}

func (s *MemoryStorage) DeletePost(postID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.posts[postID]
	if !exists {
		return fmt.Errorf("post with ID %s not found", postID)
	}

	delete(s.posts, postID)
	return nil
}
