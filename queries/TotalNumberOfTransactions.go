package queries

import (
	"database/sql"
	"log"
)

type TotalNumberOfTransactions struct {
	name string
	db   *sql.DB
}

func CreateTotalNumberOfTransactions(db *sql.DB) Query {
	return TotalNumberOfTransactions{
		db:   db,
		name: "TotalNumberOfTransactions",
	}
}

func (q TotalNumberOfTransactions) GetValue() int {
	var value int
	err := q.db.QueryRow("SELECT sum(xact_commit+xact_rollback) as \"value\" FROM pg_stat_database").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (q TotalNumberOfTransactions) GetName() string {
	return q.name
}
