package main

import (
	"database/sql"
	"log"

	"github.com/anewgd/simple-bank/api"
	db "github.com/anewgd/simple-bank/db/sqlc"
	"github.com/anewgd/simple-bank/util"

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
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	log.Fatal("Can not start server:", err)
}
