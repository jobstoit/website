package repo

import (
	"context"
	"database/sql"

	"github.com/jobstoit/website/dbc"
)

// Repo consists of all the repository functions and is responsible for data storage
type Repo struct {
	db *sql.DB
}

func New(dbConnectionString string) *Repo {
	x := new(Repo)

	x.db = dbc.Open(dbConnectionString)

	return x
}

type querier interface {
	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
}
