package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"
)

// UserCreate creates a new user
func (x Repo) UserCreate(ctx context.Context, username, password string) (id int) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	q := `INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id;`

	if err := x.db.QueryRowContext(ctx, q, username, hash).Scan(&id); err != nil {
		panic(err)
	}

	return
}

// UserMatchPassword returns a session token if the given user's password matches the database entry's hash
func (x Repo) UserMatchPassword(ctx context.Context, username, password string) string {
	tx, err := x.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() // nolint: errcheck

	q := `SELECT id, password_hash
		FROM users
		WHERE username = $1;`

	var userID int
	var hash []byte
	if err := tx.QueryRowContext(ctx, q, username).Scan(&userID, &hash); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return ``
	}

	exp := time.Now().Add(time.Hour * 12)

	claims := &jwt.StandardClaims{
		ExpiresAt: exp.Unix(),
		Issuer:    username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(x.signingKey))
	if err != nil {
		panic(err)
	}

	q = `INSERT INTO sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3);`

	if _, err := tx.Exec(q, userID, ss, exp); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return ss
}

// UserDelete deletes the given user by username
func (x Repo) UserDelete(ctx context.Context, username string) {
	q := `DELETE FROM users WHERE username = $1;`

	if _, err := x.db.ExecContext(ctx, q, username); err != nil {
		panic(err)
	}
}
