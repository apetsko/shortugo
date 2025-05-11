package grpc

import (
	grpch "github.com/apetsko/shortugo/internal/server/grpc/handlers"
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RouterGRPC sets up and returns a new gRPC server instance.
// It registers the URLShortener service implementation and enables server reflection
// for easier testing and introspection (e.g., via grpcurl).
//
// Parameters:
//   - h: pointer to the shared HTTP URLHandler, reused for business logic.
//
// Returns:
//   - *grpc.Server: the fully initialized gRPC server.
func RouterGRPC(h *httph.URLHandler) *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterURLShortenerServer(server, grpch.NewHandler(h))
	reflection.Register(server)
	return server
}
