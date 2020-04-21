package queries

import (
	"database/sql"
	"log"
)

type NumberOfConnections struct {
	name string
	db   *sql.DB
}

func CreateNumberOfConnections(db *sql.DB) Query {
	return NumberOfConnections{
		db:   db,
		name: "NumberOfConnections",
	}
}

func (q NumberOfConnections) GetValue() int {
	var value int
	err := q.db.QueryRow("select count(*) as \"value\" from pg_stat_activity").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (q NumberOfConnections) GetName() string {
	return q.name
}
