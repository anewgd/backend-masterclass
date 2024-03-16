package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	db "github.com/anewgd/simple-bank/db/sqlc"
	"github.com/anewgd/simple-bank/gapi"
	"github.com/anewgd/simple-bank/pb"
	"github.com/anewgd/simple-bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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

	go runGatewayServer(config, store)
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

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("cannot register handler server:%v", err)
	}

	//Create a new HTTP serve mux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener:%v", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("cannot start HTTP gateway server:%v", err)
	}

}
