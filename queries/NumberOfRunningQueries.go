package queries

import (
	"database/sql"
	"log"
)

func NumberOfRunningQueries(db *sql.DB) int {
	var value int
	err := db.QueryRow("SELECT " +
		"count(pid)  as \"value\" " +
		"FROM pg_stat_activity " +
		"WHERE query != '<IDLE>' " +
		"AND query NOT LIKE '%pg_stat_activity%' " +
		"AND usename IS NOT NULL " +
		"AND query IS NOT NULL").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}
