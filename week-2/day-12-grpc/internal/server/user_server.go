// Package server implements the gRPC UserService defined in proto/user.proto.
// It uses an in-memory map for storage, protected by an RWMutex so concurrent
// reads don't block each other while writes remain exclusive.
package server

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/meedaycodes/day12-grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements proto.UserServiceServer. Embedding
// UnimplementedUserServiceServer satisfies the interface for any RPC methods
// not explicitly implemented, and provides forward compatibility when new
// methods are added to the proto without breaking existing servers.
type UserServer struct {
	proto.UnimplementedUserServiceServer
	user map[string]*proto.User
	mu   sync.RWMutex
}

// NewUserServer creates a UserServer with an initialized user map.
func NewUserServer() *UserServer {
	userMap := make(map[string]*proto.User)
	return &UserServer{user: userMap}
}

// GetUser looks up a user by ID. Uses RLock so multiple concurrent reads can
// proceed in parallel — only writes need an exclusive lock. Returns a gRPC
// NotFound status error if the ID does not exist, which the framework
// automatically translates to the appropriate gRPC status code on the wire.
func (u *UserServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.User, error) {

	u.mu.RLock()
	defer u.mu.RUnlock()
	val, ok := u.user[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return val, nil

}

// CreateUser generates a UUID, constructs a User from the request fields,
// stores it under the new ID, and returns it. Lock (not RLock) is required
// because this is a write operation — no concurrent reads can proceed while
// the map is being modified.
func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {

	id := uuid.New().String()
	user := &proto.User{Id: id, Name: req.Name, Email: req.Email}
	u.mu.Lock()
	defer u.mu.Unlock()

	u.user[id] = user

	return user, nil
}
