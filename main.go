package main

import (
	"database/sql"
	"log"

	"github.com/hhong0326/goPostgresqlDocker.git/api"
	db "github.com/hhong0326/goPostgresqlDocker.git/db/sqlc"
	"github.com/hhong0326/goPostgresqlDocker.git/util"
	_ "github.com/lib/pq"
)

// main setup db first and start api server
func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(" create create server: ", err)
	}

	server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
