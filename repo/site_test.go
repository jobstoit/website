package repo

import (
	"context"
	"database/sql"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/jobstoit/website/dbc"
	"github.com/jobstoit/website/model"
)

func TestGetActiveSite(t *testing.T) {
	x, as := initTest(t)

	site := x.GetActiveSite(context.Background())
	as.Eq("testsite", site.Name)
	as.Eq(2, len(site.Pages))
	as.Eq(3, len(site.Pages[0].Rows))
}

func TestGetSiteByID(t *testing.T) {
	x, as := initTest(t)

	siteID := 101

	site := x.GetSiteByID(context.Background(), siteID)
	as.Eq("testsite", site.Name)
	as.Eq(2, len(site.Pages))
	as.Eq(3, len(site.Pages[0].Rows))
}

func TestListSites(t *testing.T) {
	x, as := initTest(t)

	sites := x.ListSites(context.Background())
	as.Eq(1, len(sites))
}

func TestCreateSite(t *testing.T) {
	x, as := initTest(t)

	id := x.CreateSite(context.Background(), "main1", "test@example.com")
	as.True(id > 0)

	q := `SELECT id FROM sites WHERE id = $1;`
	as.NoError(x.db.QueryRow(q, id).Scan(new(int)))
}

func TestCreatePage(t *testing.T) {
	x, as := initTest(t)

	id := x.CreatePage(context.Background(), 101, `/home`, `home`)
	as.True(id > 0)

	q := `SELECT id FROM pages WHERE id = $1;`
	as.NoError(x.db.QueryRow(q, id).Scan(new(int)))
}

func TestCreateRow(t *testing.T) {
	x, as := initTest(t)

	pageID := 102
	titles := []string{"uno", "dos"}
	texts := []string{"a fantastic application and know an read more about it here",
		"bla bla bla some lorum ipsom to demosntrate this "}
	media := []model.Medium{
		model.Medium{
			URI:  "https://example.com/assets/123.png",
			Type: "image",
		},
	}
	buttons := []model.Button{
		model.Button{
			URI:   "/about",
			Label: "About",
		},
		model.Button{
			URI:   "https://example.com/lorum/ipsum",
			Label: "Lorum Ipsum",
		},
	}

	id := x.CreateRow(context.Background(), pageID, titles, texts, media, buttons)

	q := `SELECT id FROM rows WHERE id = $1;`
	as.NoError(x.db.QueryRow(q, id).Scan(new(int)))

	q = `SELECT COUNT(*) FROM row_titles WHERE row_id = $1;`
	var titleCount int
	as.NoError(x.db.QueryRow(q, id).Scan(&titleCount))
	as.Eq(len(titles), titleCount)
}

func TestChangeRowSequence(t *testing.T) {
	x, as := initTest(t)

	pageID := 101
	newSequence := []int{103, 101, 102}
	x.ChangeRowSequence(context.Background(), pageID, newSequence)

	q := `SELECT id FROM rows WHERE page_id = $1 ORDER BY sequence ASC;`
	rows, err := x.db.Query(q, pageID)
	if err != nil {
		t.Log(err)
	}

	var seq []int
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			t.Log(err)
		}
		seq = append(seq, i)
	}

	if err := rows.Close(); err != nil {
		t.Log(err)
	}

	as.Eq(newSequence[0], seq[0])
	as.Eq(newSequence[1], seq[1])
	as.Eq(newSequence[2], seq[2])
}

func TestUpdateRow(t *testing.T) {
	x, as := initTest(t)

	rowID := 105
	titles := []string{"Websites"}
	texts := []string{
		"a fantastic application and know an read more about it here",
		"bla bla bla some lorum ipsom to demosntrate this",
		"and yet some more lorum ipsum stuff for demonstration purposes",
	}
	media := []model.Medium{}
	buttons := []model.Button{
		model.Button{
			URI:   "/info",
			Label: "Info",
		},
	}

	x.UpdateRow(context.Background(), rowID, titles, texts, media, buttons)

	q := `SELECT uri, label FROM row_buttons WHERE row_id = $1;`
	var uri, label string
	if err := x.db.QueryRow(q, rowID).Scan(&uri, &label); err != nil {
		t.Logf("error quering for the uri: %v", err)
	}
	as.Eq(buttons[0].URI, uri)
	as.Eq(buttons[0].Label, label)

	q = `SELECT COUNT(*) FROM row_texts WHERE row_id = $1;`
	var c int
	if err := x.db.QueryRow(q, rowID).Scan(&c); err != nil {
		t.Logf("error quering the count of the row texts: %v", err)
	}
	as.Eq(3, c)
}

func TestDeleteRow(t *testing.T) {
	x, as := initTest(t)

	rowID := 104
	x.DeleteRow(context.Background(), rowID)

	q := `SELECT id FROM row_titles WHERE row_id = $1;`
	err := x.db.QueryRow(q, rowID).Scan(new(int))
	if err != nil && err != sql.ErrNoRows {
		t.Log(err)
	}
	as.Eq(err, sql.ErrNoRows, "deleted row should not contain any more atributes")

	q = `SELECT id FROM rows WHERE id = $1;`
	err = x.db.QueryRow(q, rowID).Scan(new(int))
	if err != nil && err != sql.ErrNoRows {
		t.Log(err)
	}
	as.Eq(err, sql.ErrNoRows, "deleted row should not exist anymore")
}

func initTest(t *testing.T) (*Repo, assert.Assert) {
	repo := new(Repo)
	repo.db = dbc.OpenTest()
	ass := assert.New(t)
	return repo, ass
}
