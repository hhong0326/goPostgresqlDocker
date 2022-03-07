package main

import (
	"database/sql"
	"log"

	"github.com/hhong0326/goPostgresqlDocker.git/api"
	db "github.com/hhong0326/goPostgresqlDocker.git/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

// main setup db first and start api server
func main() {

	var err error
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
