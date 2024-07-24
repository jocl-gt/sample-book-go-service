package models

type Era struct {
	ID      int    `db:"id" json:"id"`
	Title   string `db:"title" json:"title"`
	MinYear *int   `db:"min_year" json:"minYear,omitempty"`
	MaxYear *int   `db:"max_year" json:"maxYear,omitempty"`
}
