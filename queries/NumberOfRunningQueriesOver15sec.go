package queries

import (
	"database/sql"
	"log"
)

func NumberOfRunningQueriesOver15sec(db *sql.DB) int {
	var value int

	row := db.QueryRow("SELECT count(pid) AS \"value\"  FROM pg_stat_activity WHERE query != '<IDLE>' AND query NOT ILIKE '%pg_stat_activity%' AND usename IS NOT NULL AND query IS NOT NULL")
	err := row.Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}
