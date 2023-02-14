package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// Connecting to the database
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Creating the repository with the database connection
	store := db.NewStore(conn)
	// Starting the API with the 'store' object
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
