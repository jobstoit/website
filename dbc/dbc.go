// Package dbc contains the database context
//
package dbc

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"
)

//go:embed migrate/* test/*
var migration embed.FS

// Open opens a new instance of a initialized postgres database
func Open(cs string) *sql.DB {
	db := open(cs)

	return db
}

// Open test opens a test database
func OpenTest() *sql.DB {
	db := open(os.Getenv(`TEST_DB_CONNECTION_STRING`))

	migrate(db, `test`)

	return db
}

func open(cs string) *sql.DB {
	db, err := sql.Open(`postgres`, cs)
	if err != nil {
		panic(err)
	}

	migrate(db, `migrate`)

	return db
}

func migrate(db *sql.DB, folder string) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() // nolint: errcheck

	versionTable := folder + `_version`
	if _, err := tx.Exec(`CREATE TABLE IF NOT EXISTS ` + versionTable + ` (id int PRIMARY KEY);`); err != nil {
		panic(err)
	}

	var version int
	if err := tx.QueryRow(`SELECT id FROM ` + versionTable + ` ORDER BY id DESC;`).Scan(&version); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	names := assetNames(folder)
	numbers := sqlNumbers(folder, names...)
	highest := hn(numbers...)

	for i := version + 1; i <= highest; i++ {
		srp, err := migration.ReadFile(fmt.Sprintf("%s/%04d.sql", folder, i))
		if err != nil {
			panic(err)
		}

		if _, err := tx.Exec(`INSERT INTO `+versionTable+` (id) VALUES ($1);`, i); err != nil {
			panic(err)
		}

		if _, err := tx.Exec(string(srp)); err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func assetNames(folder string) (names []string) {
	dir, err := migration.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	for _, e := range dir {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}

	return
}

func sqlNumbers(folder string, assets ...string) []int {
	if len(assets) < 1 {
		return []int{}
	}

	ns := sqlNumbers(folder, assets[1:]...)
	reg := regexp.MustCompile(`^([0-9]{4})\.sql$`)

	name := assets[0]
	if !reg.MatchString(name) {
		return ns
	}

	number := reg.FindStringSubmatch(name)[1]
	i, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}
	return append(ns, i)
}

func hn(n ...int) int {
	if len(n) < 1 {
		return 0
	}

	if h := hn(n[1:]...); n[0] < h {
		return h
	}

	return n[0]
}
