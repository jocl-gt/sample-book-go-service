package models

type Book struct {
	ID            int     `db:"id" json:"id"`
	Title         string  `db:"title" json:"title"`
	YearPublished int     `db:"year_published" json:"yearPublished"`
	Rating        float64 `db:"rating" json:"rating"`
	Pages         int     `db:"pages" json:"pages"`
	GenreID       int     `db:"genre_id" json:"-"`
	AuthorID      int     `db:"author_id" json:"-"`
	Genre         Genre   `db:"-" json:"genre"`
	Author        Author  `db:"-" json:"author"`
}
