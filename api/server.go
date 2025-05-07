package api

import (
	"fmt"
	"log"

	"github.com/sangketkit01/simple-grpc-auth/config"
	db "github.com/sangketkit01/simple-grpc-auth/db/sqlc"
	"github.com/sangketkit01/simple-grpc-auth/pb"
	"github.com/sangketkit01/simple-grpc-auth/token"
)

type Server struct {
	pb.UnimplementedGrpcSimpleAuthServer
	tokenMaker token.Maker
	store *db.Store
	config config.Config
}

func NewServer(store *db.Store, config config.Config) (*Server, error){
	tokenMaker, err := token.NewPasetoMaker(config.SecretKey)
	if err != nil{
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store: store,
		tokenMaker: tokenMaker,
		config: config,
	}

	return server, nil
}