// Package handlers implements the gRPC server for the URL shortening service.
//
// It defines the Handler struct, which delegates business logic to an underlying
// URLHandler instance shared with the HTTP layer. This design ensures consistent
// behavior and centralized logic for storage operations and user authentication.
package handlers
