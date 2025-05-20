package handlers

import (
	"context"
	"fmt"
	"net"

	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func startGRPCServer(handler pb.URLShortenerServer) (*grpc.ClientConn, func(), error) {
	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()
	pb.RegisterURLShortenerServer(s, handler)

	go func() {
		_ = s.Serve(lis)
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("failed to close bufnet connection: %v", err)
		}

		s.Stop()
		if err := lis.Close(); err != nil {
			fmt.Printf("failed to close bufnet connection: %v", err)
		}
	}

	return conn, cleanup, nil
}
