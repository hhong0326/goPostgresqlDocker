package gapi

import (
	"context"

	db "github.com/hhong0326/goPostgresqlDocker.git/db/sqlc"
	"github.com/hhong0326/goPostgresqlDocker.git/pb"
	"github.com/hhong0326/goPostgresqlDocker.git/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	// internal validator exec

	hashedPw, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPw,
		Fullname:       req.GetFullname(),
		Email:          req.GetEmail(),
	}

	// test must be error
	// arg = db.CreateUserParams{}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation": // check duplicated
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}
