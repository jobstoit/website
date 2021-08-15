package dbc

import (
	"database/sql"
	"testing"
)

func TestOpen(t *testing.T) {
	defer func() {
		if v := recover(); v != nil {
			t.Error(v)
		}
	}()

	db := OpenTest()

	q := `SELECT id FROM migrate_version ORDER BY id DESC LIMIT 1;`

	var migrationVersion int
	if err := db.QueryRow(q).Scan(&migrationVersion); err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}

	t.Logf("migration version: %d", migrationVersion)

}
