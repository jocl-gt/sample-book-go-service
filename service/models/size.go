package models

type Size struct {
	ID       int    `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	MinPages *int   `db:"min_pages" json:"minPages,omitempty"`
	MaxPages *int   `db:"max_pages" json:"maxPages,omitempty"`
}
