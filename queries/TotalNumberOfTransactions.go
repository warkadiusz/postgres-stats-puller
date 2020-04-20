package queries

import (
	"database/sql"
	"log"
)

func TotalNumberOfTransactions(db *sql.DB) int {
	var value int
	err := db.QueryRow("SELECT sum(xact_commit+xact_rollback) as \"value\" FROM pg_stat_database").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}
