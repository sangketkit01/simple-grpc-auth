package api

import (
	"fmt"

	"github.com/sangketkit01/simple-grpc-auth/db"
	"github.com/sangketkit01/simple-grpc-auth/pb"
	"github.com/sangketkit01/simple-grpc-auth/token"
)

type Server struct {
	pb.UnimplementedGrpcSimpleAuthServer
	tokenMaker token.Maker
	store db.Store
}

func NewServer(store db.Store) (*Server, error){
	tokenMaker, err := token.NewPasetoMaker("12345678901234567890123456789012345678901234")
	if err != nil{
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store: store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}