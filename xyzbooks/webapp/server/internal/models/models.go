package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type Book struct {
	ID              string          `db:"id"`
	Title           string          `db:"title"`
	ISBN13          string          `db:"isbn13"`
	ISBN10          string          `db:"isbn10"`
	PublisherID     string          `db:"publisher_id"`
	PublicationYear int             `db:"publication_year"`
	Edition         sql.NullString  `db:"edition"`
	ImageURL        sql.NullString  `db:"image_url"`
	Price           decimal.Decimal `db:"price"`
	Dates
}

type Author struct {
	ID         string         `db:"id"`
	FirstName  string         `db:"first_name"`
	LastName   string         `db:"last_name"`
	MiddleName sql.NullString `db:"middle_name"`
	Dates
}

type Publisher struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Dates
}

type AuthorBook struct {
	AuthorID string `db:"author_id"`
	BookID   string `db:"book_id"`
	ID       string `db:"id"`
	Dates
}

type Dates struct {
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}
