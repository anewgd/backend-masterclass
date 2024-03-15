package main

import (
	"database/sql"
	"log"
	"net"

	db "github.com/anewgd/simple-bank/db/sqlc"
	"github.com/anewgd/simple-bank/gapi"
	"github.com/anewgd/simple-bank/pb"
	"github.com/anewgd/simple-bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {

		log.Fatal("Cannot load configuration file: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cant connect to db", err)
	}

	store := db.NewStore(conn)
	runGRPCServer(config, store)

}

// func runGinServer(config util.Config, store db.Store) {
// 	server, err := api.NewServer(config, store)

// 	if err != nil {
// 		log.Fatal("cannot create server: ", err)
// 	}
// 	err = server.Start(config.HTTPServerAddress)
// 	log.Fatal("Can not start server:", err)
// }

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	gRPCServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(gRPCServer, server)

	// Allows the gRPC client to explore what RPCs
	// are available on the server.	
	reflection.Register(gRPCServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener:%v", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())

	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start gRPC server:%v", err)
	}

}
