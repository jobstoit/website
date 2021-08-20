package repo

import (
	"context"
	"database/sql"
	"testing"
)

func TestUserCreate(t *testing.T) {
	x, as := initTest(t)

	username, password := `test3`, `SuperComplexPassword`
	id := x.UserCreate(context.Background(), username, password)
	as.True(id > 0)

	q := `SELECT id FROM users WHERE id = $1 AND username = $2`
	as.NoError(x.db.QueryRow(q, id, username).Scan(new(int)))
}

func TestUserMatchPassword(t *testing.T) {
	x, as := initTest(t)

	ctx := context.Background()
	username, password := `test1`, `SuperComplexPassword`
	as.Eq(``, x.UserMatchPassword(ctx, username, `falsePassword`))
	as.Ne(``, x.UserMatchPassword(ctx, username, password))
}

func TestUserDelete(t *testing.T) {
	x, as := initTest(t)

	username := `test2`
	x.UserDelete(context.Background(), username)

	q := `SELECT id FROM users WHERE username = $1;`
	as.Eq(sql.ErrNoRows, x.db.QueryRow(q, username).Scan(new(int)))
}
