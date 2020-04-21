package queries

import (
	"database/sql"
	"errors"
)

type QueryFactory struct {
	db *sql.DB
}

func CreateQueryFactory(db *sql.DB) QueryFactory {
	return QueryFactory{
		db: db,
	}
}

func (qf QueryFactory) Create(name string) (Query, error) {
	switch name {
	case "DatabaseSize":
		return CreateDatabaseSizeQuery(qf.db), nil
	case "NumberOfConnections":
		return CreateNumberOfConnections(qf.db), nil
	case "NumberOfRunningQueries":
		return CreateNumberOfRunningQueries(qf.db), nil
	case "NumberOfRunningQueriesOver15sec":
		return CreateNumberOfRunningQueriesOver15sec(qf.db), nil
	case "TotalNumberOfTransactions":
		return CreateTotalNumberOfTransactions(qf.db), nil
	}

	return nil, errors.New("Query of type " + name + " not found")
}
