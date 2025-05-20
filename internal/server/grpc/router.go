package grpc

import (
	logger "github.com/apetsko/shortugo/internal/logging"
	grpch "github.com/apetsko/shortugo/internal/server/grpc/handlers"
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
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
	zaplogger := zap.NewExample()

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	// You can now create a server with logging instrumentation that e.g. logs when the unary or stream call is started or finished.
	_ = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(logger.InterceptorLogger(zaplogger), opts...),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(logger.InterceptorLogger(zaplogger), opts...),
			// Add any other interceptor you want.
		),
	)
	server := grpc.NewServer()
	pb.RegisterURLShortenerServer(server, grpch.NewHandler(h))
	reflection.Register(server)

	return server
}
