package queries

import (
	"database/sql"
	"log"
)

type NumberOfRunningQueriesOver15sec struct {
	name string
	db   *sql.DB
}

func CreateNumberOfRunningQueriesOver15sec(db *sql.DB) Query {
	return NumberOfRunningQueriesOver15sec{
		db:   db,
		name: "NumberOfRunningQueriesOver15sec",
	}
}

func (q NumberOfRunningQueriesOver15sec) GetValue() int {
	var value int

	row := q.db.QueryRow("SELECT count(pid) AS \"value\"  FROM pg_stat_activity WHERE query != '<IDLE>' AND query NOT ILIKE '%pg_stat_activity%' AND usename IS NOT NULL AND query IS NOT NULL")
	err := row.Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (q NumberOfRunningQueriesOver15sec) GetName() string {
	return q.name
}
