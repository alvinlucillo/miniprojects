package repos

import (
	"alvinlucillo/xyzbooks_webapp/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	packageName = "repos"
)

type Repository interface {
	GetBook(id string) (models.Book, error)
	GetBooks() ([]models.Book, error)
	GetBooksByPublisherID(id string) ([]models.Book, error)
	GetBookByISBN13(isbn13 string) (*models.Book, error)
	GetBookByISBN10(isbn10 string) (*models.Book, error)
	CreateBook(book models.Book) (string, error)
	UpdateBook(book models.Book) error
	DeleteBook(id string) error
	GetAuthors() ([]models.Author, error)
	GetAuthor(id string) (models.Author, error)
	CreateAuthor(author models.Author) (string, error)
	UpdateAuthor(author models.Author) error
	DeleteAuthor(id string) error
	GetPublishers() ([]models.Publisher, error)
	GetPublisher(id string) (models.Publisher, error)
	CreatePublisher(publisher models.Publisher) (string, error)
	UpdatePublisher(publisher models.Publisher) error
	DeletePublisher(id string) error
	GetAuthorBookRelByBookID(id string) ([]models.AuthorBook, error)
	GetAuthorBookRelByAuthorID(id string) ([]models.AuthorBook, error)
	AddAuthorBookRel(id string, authorId string) error
	DeleteAuthorBookRelByBookID(id string) error
	DeleteAuthorBookRelByAuthorID(id string) error
}

type repo struct {
	tx     *sqlx.Tx
	logger zerolog.Logger
}

func NewRepository(tx *sqlx.Tx, logger zerolog.Logger) Repository {
	return &repo{
		tx:     tx,
		logger: logger,
	}
}
