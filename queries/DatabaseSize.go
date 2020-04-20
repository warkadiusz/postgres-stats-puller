package queries

import (
	"database/sql"
	"log"
)

func DatabaseSize(db *sql.DB) int {
	var value int
	db.Exec("ANALYZE");
	err := db.QueryRow("select sum(pg_database_size(datname)) as \"value\" from pg_database").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}
