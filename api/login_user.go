package api

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/sangketkit01/simple-grpc-auth/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if err := validateLoginUserRequest(req); err != nil {
		return nil, err
	}

	user, err := server.store.LoginUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "cannot get user: %s", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, status.Errorf(codes.Unauthenticated, "password does not match")
		}

		return nil, status.Errorf(codes.Internal, "compared password failed: %s", err)
	}

	userData, err := server.store.GetUser(ctx, user.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	token, payload, err := server.tokenMaker.CreateToken(userData.Username, server.config.AccessKeyDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %s", err)
	}

	reponse := &pb.LoginUserResponse{
		User: &pb.User{
			Username:  userData.Username,
			FullName:  userData.FullName,
			Email:     userData.Email,
			CreatedAt: timestamppb.New(userData.CreatedAt),
		},
		SessionId: payload.TokenID.String(),
		AccessToken: token,
		AccessTokenIssuedAt: timestamppb.New(payload.IssuedAt),
		AccessTokenExpiredAt: timestamppb.New(payload.ExpiredAt),
	}

	return reponse, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}
	if len(req.GetUsername()) < 3 || len(req.GetUsername()) > 30 {
		return status.Error(codes.InvalidArgument, "username must be between 3 and 30 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.GetUsername()) {
		return status.Error(codes.InvalidArgument, "username can only contain letters, numbers and underscores")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if len(req.GetPassword()) < 8 {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(req.GetPassword()) ||
		!regexp.MustCompile(`[a-z]`).MatchString(req.GetPassword()) ||
		!regexp.MustCompile(`[0-9]`).MatchString(req.GetPassword()) {
		return status.Error(codes.InvalidArgument, "password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}
