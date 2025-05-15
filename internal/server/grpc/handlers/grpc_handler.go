package handlers

import (
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
)

// Handler implements the gRPC server interface for the URL shortening service.
// It delegates business logic to the underlying HTTP URLHandler instance.
type Handler struct {
	pb.UnimplementedURLShortenerServer                   // Embeds the unimplemented server for forward compatibility.
	URLHandler                         *httph.URLHandler // Reference to the shared HTTP handler logic.
}

// NewHandler creates a new gRPC handler instance that wraps the existing HTTP handler.
// This allows reusing the same logic for both HTTP and gRPC endpoints.
func NewHandler(h *httph.URLHandler) *Handler {
	return &Handler{URLHandler: h}
}
