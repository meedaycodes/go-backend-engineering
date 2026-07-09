// Package main is the entry point for the gRPC server. It creates a TCP
// listener, registers the UserService implementation, and starts serving.
// Dependencies flow: listener → grpc.Server → UserServer → handler.
// Port 50051 is the conventional default for gRPC services.
package main

import (
	"log"
	"net"

	"github.com/meedaycodes/day12-grpc/internal/server"
	"github.com/meedaycodes/day12-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// net.Listen creates a TCP socket but does not accept connections yet.
	// Serving begins when grpcServer.Serve(lis) is called below.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	// RegisterUserServiceServer binds our implementation to the gRPC server.
	// The generated function wires the server so incoming RPCs are dispatched
	// to the correct method on UserServer.
	proto.RegisterUserServiceServer(grpcServer, server.NewUserServer())
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :50051")

	// Serve blocks until the server is stopped or encounters a fatal error.
	// log.Fatal ensures the process exits with a non-zero code on failure.
	log.Fatal(grpcServer.Serve(lis))
}
