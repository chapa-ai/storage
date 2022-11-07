package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "storage/grpc-storage"
	"sync"
)

const (
	port = ":50051"
)

var mu sync.Mutex
var m = make(map[string][]byte)

type StorageManagementServer struct {
	pb.UnimplementedStorageManagementServer
}

func (s *StorageManagementServer) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetResponse, error) {
	mu.Lock()
	m[in.Key] = in.Value
	mu.Unlock()
	return &pb.SetResponse{Key: in.Key, ResultStored: in.Value}, nil
}

func (s *StorageManagementServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	mu.Lock()
	value, exists := m[in.Key]
	if !exists {
		return nil, fmt.Errorf("no such key in map")
	}
	mu.Unlock()
	return &pb.GetResponse{Key: in.Key, ResultOK: value}, nil
}

func (server *StorageManagementServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	mu.Lock()
	delete(m, in.Key)

	str := "Deleted"

	del := &pb.DeleteResponse{
		ResultDeleted: []byte(str),
	}

	mu.Unlock()
	return del, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStorageManagementServer(s, &StorageManagementServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
