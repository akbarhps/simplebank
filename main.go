package main

import (
	"database/sql"
	"github.com/akbarhps/simplebank/api"
	db "github.com/akbarhps/simplebank/db/sqlc"
	"github.com/akbarhps/simplebank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("error read config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("error creating db conn:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("error starting server:", err)
	}
}
