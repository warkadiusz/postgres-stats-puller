package queries

import (
	"database/sql"
	"log"
)

type NumberOfRunningQueries struct {
	name string
	db   *sql.DB
}

func CreateNumberOfRunningQueries(db *sql.DB) Query {
	return NumberOfRunningQueries{
		db:   db,
		name: "NumberOfRunningQueries",
	}
}

func (q NumberOfRunningQueries) GetValue() int {
	var value int
	err := q.db.QueryRow("SELECT " +
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

func (q NumberOfRunningQueries) GetName() string {
	return q.name
}
