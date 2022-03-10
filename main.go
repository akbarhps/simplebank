package main

import (
	"database/sql"
	"github.com/akbarhps/simplebank/api"
	db "github.com/akbarhps/simplebank/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("error creating db conn:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err = server.Start(serverAddress); err != nil {
		log.Fatal("error starting server:", err)
	}
}
