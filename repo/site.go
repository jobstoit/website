// Copyright 2021 Job Stoit. All rights reseved

package repo

import (
	"context"
	"database/sql"

	"github.com/jobstoit/website/model"
)

// CreateSite adds a new site to the repository
func (x Repo) CreateSite(ctx context.Context, name string, createdBy int) (id int) {
	q := `INSERT INTO sites (name, created_by)
		VALUES ($1, $2)
		RETURNING id;`

	if err := x.db.QueryRowContext(ctx, q, name, createdBy).Scan(&id); err != nil {
		panic(err)
	}

	return
}

// CreatePage adds a page to the site into the repostirory
func (x Repo) CreatePage(ctx context.Context, siteID int, uri, label string) (id int) {
	q := `INSERT INTO pages (site_id, uri_path, label)
		VALUES ($1, $2, $3)
		RETURNING id;`

	if err := x.db.QueryRowContext(ctx, q, siteID, uri, label).Scan(&id); err != nil {
		panic(err)
	}

	return
}

// CreateRow adds a new row to the given page
func (x Repo) CreateRow(ctx context.Context, pageID int, titles, texts []string, media []model.Medium, btns []model.Button) (id int) {
	tx, err := x.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() // nolint: errcheck

	qs := `SELECT sequence
		FROM rows
		WHERE page_id = $1
		ORDER BY sequence ASC;`

	var sequence int
	if err := tx.QueryRowContext(ctx, qs, pageID).Scan(&sequence); err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	sequence += 1

	q := `INSERT INTO rows (page_id, sequence)
		VALUES ($1, $2)
		RETURNING id;`

	if err := tx.QueryRowContext(ctx, q, pageID, sequence).Scan(&id); err != nil {
		panic(err)
	}

	insertRowTitles(tx, id, titles)
	insertRowTexts(tx, id, texts)
	insertRowMedia(tx, id, media)
	insertRowButtons(tx, id, btns)

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return
}

// insertRowTitles adds titles to the given row
func insertRowTitles(tx querier, rowID int, titles []string) {
	q := `INSERT INTO row_titles (row_id, sequence, context)
		VALUES($1, $2, $3);`

	stmt, err := tx.Prepare(q)
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // nolint: errcheck

	for i, tl := range titles {
		if _, err := stmt.Exec(rowID, i, tl); err != nil {
			panic(err)
		}
	}
}

// insertRowTexts add texts to the given row
func insertRowTexts(tx querier, rowID int, texts []string) {
	q := `INSERT INTO row_texts (row_id, sequence, context)
		VALUES($1, $2, $3);`

	stmt, err := tx.Prepare(q)
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // nolint: errcheck

	for i, tx := range texts {
		if _, err := stmt.Exec(rowID, i, tx); err != nil {
			panic(err)
		}
	}
}

// insertRowMedia adds media to the given row
func insertRowMedia(tx querier, rowID int, media []model.Medium) {
	q := `INSERT INTO row_media (row_id, sequence, uri, type)
		VALUES($1, $2, $3, $4);`

	stmt, err := tx.Prepare(q)
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // nolint: errcheck

	for i, md := range media {
		if _, err := stmt.Exec(rowID, i, md.URI, md.Type); err != nil {
			panic(err)
		}
	}
}

// insertRowmodel.Buttons adds buttons to the given row
func insertRowButtons(tx querier, rowID int, btns []model.Button) {
	q := `INSERT INTO row_buttons (row_id, sequence, uri, label)
		VALUES($1, $2, $3, $4);`

	stmt, err := tx.Prepare(q)
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // nolint: errcheck

	for i, btn := range btns {
		if _, err := stmt.Exec(rowID, i, btn.URI, btn.Label); err != nil {
			panic(err)
		}
	}
}

// ChangeRowSequence alters the sequence of a pages rows
func (x Repo) ChangeRowSequence(ctx context.Context, pageID int, rowIDs []int) {
	q := `UPDATE rows
		SET sequence = $3
		WHERE page_id = $1 AND id = $2;`

	stmt, err := x.db.PrepareContext(ctx, q)
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // nolint: errcheck

	for seq, id := range rowIDs {
		if _, err := stmt.Exec(pageID, id, seq); err != nil {
			panic(err)
		}
	}
}

// UpdateRow deletes previous attributes of the row and inserts the new given ones
func (x Repo) UpdateRow(ctx context.Context, rowID int, titles, texts []string, media []model.Medium, btns []model.Button) {
	tx, err := x.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() // nolint: errcheck

	deleteRowAttributes(tx, rowID)
	insertRowTitles(tx, rowID, titles)
	insertRowTexts(tx, rowID, texts)
	insertRowMedia(tx, rowID, media)
	insertRowButtons(tx, rowID, btns)

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

// DeleteRow deletes a given row
func (x Repo) DeleteRow(ctx context.Context, id int) {
	tx, err := x.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() // nolint: errcheck

	deleteRowAttributes(tx, id)

	q := `DELETE FROM rows WHERE id = $1;`
	if _, err := tx.Exec(q, id); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func deleteRowAttributes(tx querier, rowID int) {
	q := `DELETE FROM row_titles WHERE row_id = $1;`
	if _, err := tx.Exec(q, rowID); err != nil {
		panic(err)
	}

	q = `DELETE FROM row_texts WHERE row_id = $1;`
	if _, err := tx.Exec(q, rowID); err != nil {
		panic(err)
	}

	q = `DELETE FROM row_media WHERE row_id = $1;`
	if _, err := tx.Exec(q, rowID); err != nil {
		panic(err)
	}

	q = `DELETE FROM row_buttons WHERE row_id = $1;`
	if _, err := tx.Exec(q, rowID); err != nil {
		panic(err)
	}
}

// ListSites returns the sites with name and id
func (x Repo) ListSites(ctx context.Context) (sts []model.SiteListItem) {
	q := `SELECT id, name
		FROM sites;`

	rows, err := x.db.QueryContext(ctx, q)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var st model.SiteListItem
		if err := rows.Scan(&st.ID, &st.Name); err != nil {
			panic(err)
		}

		sts = append(sts, st)
	}

	return
}

// GetActiveSite returns the site that has the highest activated_since value
func (x Repo) GetActiveSite(ctx context.Context) model.Site {
	q := `SELECT id, name, css_uri, title
		FROM sites
		ORDER BY active_since
		LIMIT 1;`

	return x.getSiteQuery(ctx, q)
}

// GetSiteByID returns the site by id
func (x Repo) GetSiteByID(ctx context.Context, id int) (s model.Site) {
	q := `SELECT id, name, css_uri, title
		FROM sites
		WHERE id = $1;`
	return x.getSiteQuery(ctx, q, id)
}

func (x Repo) getSiteQuery(ctx context.Context, q string, args ...interface{}) (s model.Site) {
	tx, err := x.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() // nolint: errcheck

	var nullCSSURI, nullTitle sql.NullString
	if err := tx.QueryRow(q, args...).Scan(&s.ID, &s.Name, &nullCSSURI, &nullTitle); err != nil {
		if err == sql.ErrNoRows {
			return
		}
		panic(err)
	}
	nullCSSURI.Scan(&s.Title) // nolint: errcheck
	nullTitle.Scan(&s.Title)  // nolint: errcheck

	s.Pages = getPagesBySiteID(tx, s.ID)

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return
}

func getPagesBySiteID(tx querier, siteID int) (pgs []model.Page) {
	q := `DECLARE pages_cursor CURSOR FOR
		SELECT id, uri_path, label
		FROM pages
		WHERE site_id = $1;`

	if _, err := tx.Exec(q, siteID); err != nil {
		panic(err)
	}

	q = `CLOSE pages_cursor;`
	defer tx.Exec(q) // nolint: errcheck

	q = `FETCH NEXT FROM pages_cursor;`
	var pg model.Page
	for {
		if err := tx.QueryRow(q).Scan(&pg.ID, &pg.URI, &pg.Label); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			panic(err)
		}

		pg.Rows = getRowsByPageID(tx, pg.ID)

		pgs = append(pgs, pg)
		pg = model.Page{}
	}

	return
}

func getRowsByPageID(tx querier, pageID int) (rws []model.Row) {
	q := `DECLARE row_cursor CURSOR FOR
		SELECT id
		FROM rows
		WHERE page_id = $1
		ORDER BY sequence ASC;`

	if _, err := tx.Exec(q, pageID); err != nil {
		panic(err)
	}

	q = `CLOSE row_cursor;`
	defer tx.Exec(q) // nolint: errcheck

	q = `FETCH NEXT FROM row_cursor;`
	var rw model.Row
	for {
		if err := tx.QueryRow(q).Scan(&rw.ID); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			panic(err)
		}

		rw.Titles = getTitlesByRowID(tx, rw.ID)
		rw.Texts = getTextsByRowID(tx, rw.ID)
		rw.Media = getMediaByRowID(tx, rw.ID)
		rw.Buttons = getButtonsByRowID(tx, rw.ID)

		rws = append(rws, rw)
		rw = model.Row{}
	}

	return
}

func getTitlesByRowID(tx querier, rowID int) (tls []string) {
	q := `SELECT context
		FROM row_titles
		WHERE row_id = $1
		ORDER BY sequence ASC;`

	rows, err := tx.Query(q, rowID)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // nolint: errcheck

	for rows.Next() {
		var tl string
		if err := rows.Scan(&tl); err != nil {
			panic(err)
		}
		tls = append(tls, tl)
	}

	return
}

func getTextsByRowID(tx querier, rowID int) (txs []string) {
	q := `SELECT context
		FROM row_texts
		WHERE row_id = $1
		ORDER BY sequence ASC;`

	rows, err := tx.Query(q, rowID)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // nolint: errcheck

	for rows.Next() {
		var tex string
		if err := rows.Scan(&tex); err != nil {
			panic(err)
		}
		txs = append(txs, tex)
	}

	return
}

func getMediaByRowID(tx querier, rowID int) (mds []model.Medium) {
	q := `SELECT uri, type, alt
		FROM row_media
		WHERE row_id = $1
		ORDER BY sequence ASC;`

	rows, err := tx.Query(q, rowID)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // nolint: errcheck

	for rows.Next() {
		var md model.Medium
		var nullAlt sql.NullString
		if err := rows.Scan(&md.URI, &md.Type, &nullAlt); err != nil {
			panic(err)
		}
		nullAlt.Scan(&md.Alt) // nolint: errcheck

		mds = append(mds, md)
	}

	return
}

func getButtonsByRowID(tx querier, rowID int) (btns []model.Button) {
	q := `SELECT label, uri
		FROM row_buttons
		WHERE row_id = $1
		ORDER BY sequence ASC;`

	rows, err := tx.Query(q, rowID)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // nolint: errcheck

	for rows.Next() {
		var btn model.Button
		if err := rows.Scan(&btn.Label, &btn.URI); err != nil {
			panic(err)
		}

		btns = append(btns, btn)
	}

	return
}
