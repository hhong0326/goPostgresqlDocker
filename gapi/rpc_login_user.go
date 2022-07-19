package gapi

import (
	"context"
	"database/sql"

	"github.com/hhong0326/goPostgresqlDocker.git/pb"
	"github.com/hhong0326/goPostgresqlDocker.git/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", user)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to check password: %s", err)
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	res := &pb.LoginUserResponse{
		AccessToken: accessToken,
		User:        convertUser(user),
	}

	return res, nil
}
