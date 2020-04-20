package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var postgresConnection *sql.DB

func CreatePostgresConnection() {
	var err error

	connStr := fmt.Sprintf(
		"user=%s password=%s port=%s dbname=postgres search_path=pg_catalog host=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_HOST"))
	postgresConnection, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
}

func GetPostgresConnection() *sql.DB {
	return postgresConnection
}

func ClosePostgresConnection() {
	postgresConnection.Close()
}
