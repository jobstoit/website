package model

// Site is the model for a site construction
type Site struct {
	ID        int      `json:"id"`
	CreatedBy int      `json:"created_by"`
	Name      string   `json:"name"`
	CssURI    string   `json:"css_uri"`
	Title     string   `json:"title"`
	Pages     []Page   `json:"pages"`
	NavHeader []Button `json:"nav_header"`
	NavFooter []Button `json:"nav_footer"`
	NavPanel  []Button `json:"nav_panel"`
}

// Page is a site's page
type Page struct {
	ID    int    `json:"id"`
	URI   string `json:"uri"`
	Label string `json:"label"`
	Rows  []Row  `json:"rows"`
}

// Row contains the row information
type Row struct {
	ID      int      `json:"id"`
	Style   int      `json:"style"`
	Titles  []string `json:"titles"`
	Texts   []string `json:"texts"`
	Media   []Medium `json:"media"`
	Buttons []Button `json:"buttons"`
}

// Medium contains media information
type Medium struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
	Alt  string `json:"alt"`
}

// Button contains the row button information
type Button struct {
	Label string `json:"label"`
	URI   string `json:"uri"`
}

// SitelistItem is a small reprensitation of the sites for listing the sites
type SiteListItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
