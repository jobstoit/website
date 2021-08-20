// Package repo is responsible for the application data transactions
//
package repo

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"

	"github.com/jobstoit/website/dbc"
)

// Repo consists of all the repository functions and is responsible for data storage
type Repo struct {
	signingKey string
	db         *sql.DB
}

// New returns an initialized repository
func New(dbConnectionString, signingKey string) *Repo {
	x := new(Repo)

	x.db = dbc.Open(dbConnectionString)

	if signingKey == `` {
		signingKey = randomString()
	}
	x.signingKey = signingKey

	return x
}

// NewTest returns a new initialized test repository
func NewTest() *Repo {
	x := new(Repo)

	x.db = dbc.OpenTest()

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

func randomString() string {
	var output strings.Builder
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	length := 20
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}
