package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sangketkit01/simple-grpc-auth/api"
	"github.com/sangketkit01/simple-grpc-auth/config"
	db "github.com/sangketkit01/simple-grpc-auth/db/sqlc"
	"github.com/sangketkit01/simple-grpc-auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil{
		log.Fatalln(err)
	}

	database, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil{
		log.Fatalln(err)
	}

	store, err := db.NewStore(database)
	if err != nil{
		log.Fatalln(err)
	}

	server, err := api.NewServer(store,config)
	if err != nil{
		log.Fatalln(err)
	}

	go runGatewayServer(server)
	runGrpcServer(server)
}

func runGrpcServer(server *api.Server){
	grpcServer := grpc.NewServer()
	pb.RegisterGrpcSimpleAuthServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp","0.0.0.0:9090")
	if err != nil{
		log.Fatalln(err)
	}

	err = grpcServer.Serve(listener)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("server started at port :9090")
}

func runGatewayServer(server *api.Server){
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterGrpcSimpleAuthHandlerServer(ctx, grpcMux, server)
	if err != nil{
		log.Println(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp","0.0.0.0:8080")
	if err != nil{
		log.Println(err)
	}

	err = http.Serve(listener, mux)
	if err != nil{
		log.Fatalln(err)
	}
	
}