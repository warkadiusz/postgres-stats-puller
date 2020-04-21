package queries

import (
	"database/sql"
	"log"
)

type DatabaseSizeQuery struct {
	name string
	db   *sql.DB
}

func CreateDatabaseSizeQuery(db *sql.DB) Query {
	return DatabaseSizeQuery{
		db:   db,
		name: "DatabaseSize",
	}
}

func (q DatabaseSizeQuery) GetValue() int {
	var value int
	q.db.Exec("ANALYZE")
	err := q.db.QueryRow("select sum(pg_database_size(datname)) as \"value\" from pg_database").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (q DatabaseSizeQuery) GetName() string {
	return q.name
}
