package main

import (
	"log"
	"net"

	"github.com/kpauljoseph/test/internal/server"
	"github.com/kpauljoseph/test/internal/storage"
	proto "github.com/kpauljoseph/test/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	log.Printf("Starting gRPC blog server on port %s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	storage := storage.NewMemoryStorage()
	blogServer := server.NewBlogServer(storage)

	s := grpc.NewServer()
	proto.RegisterBlogServiceServer(s, blogServer)

	log.Printf("Blog server listening at %v", lis.Addr())
	log.Println("Server ready to accept connections...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}