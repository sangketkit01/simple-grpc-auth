package api

import (
	"context"
	"regexp"
	"strings"

	"github.com/lib/pq"
	db "github.com/sangketkit01/simple-grpc-auth/db/sqlc"
	"github.com/sangketkit01/simple-grpc-auth/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hashed password %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		HashedPassword: string(hashedPassword),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// Check for unique violation error
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				if strings.Contains(pqErr.Message, "username") {
					return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", req.GetUsername())
				} else if strings.Contains(pqErr.Message, "email") {
					return nil, status.Errorf(codes.AlreadyExists, "email already exists: %s", req.GetEmail())
				}
				return nil, status.Error(codes.AlreadyExists, "unique constraint violation")
			}
		}
		return nil, status.Error(codes.Internal, "failed to create user: "+err.Error())
	}

	response := &pb.CreateUserResponse{
		User: &pb.User{
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	return response, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) error {
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

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(req.GetEmail()) {
		return status.Error(codes.InvalidArgument, "invalid email format")
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

	if req.GetFullName() != "" && len(req.GetFullName()) > 100 {
		return status.Error(codes.InvalidArgument, "full name must not exceed 100 characters")
	}

	return nil
}
