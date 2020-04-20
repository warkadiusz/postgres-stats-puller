package queries

import (
	"database/sql"
	"log"
)

func NumberOfConnections(db *sql.DB) int {
	var value int
	err := db.QueryRow("select count(*) as \"value\" from pg_stat_activity").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}