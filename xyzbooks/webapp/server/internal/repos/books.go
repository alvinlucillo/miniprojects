package repos

import "alvinlucillo/xyzbooks_webapp/internal/models"

func (r *repo) GetBooks() ([]models.Book, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetBooks").Logger()
	books := []models.Book{}

	err := r.tx.Select(&books, "SELECT * FROM book")

	if err != nil {
		l.Err(err).Msg("failed to get books")
		return nil, err
	}

	return books, nil
}

func (r *repo) GetBooksByPublisherID(id string) ([]models.Book, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetBooksByPublisherID").Logger()
	books := []models.Book{}

	err := r.tx.Select(&books, "SELECT * FROM book WHERE publisher_id = $1", id)

	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get books")
		return nil, err
	}

	return books, nil
}

func (r *repo) GetBook(id string) (models.Book, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetBook").Logger()
	book := models.Book{}

	err := r.tx.Get(&book, "SELECT * FROM book WHERE id = $1", id)

	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get book")
		return book, err
	}

	return book, nil
}

func (r *repo) GetBookByISBN13(isbn13 string) (*models.Book, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetBookByISBN13").Logger()
	book := &models.Book{}

	err := r.tx.Get(book, "SELECT * FROM book WHERE isbn13 = $1", isbn13)

	if err != nil {
		l.Err(err).Str("isbn13", isbn13).Msg("failed to get book")
		return nil, err
	}

	return book, nil
}

func (r *repo) GetBookByISBN10(isbn10 string) (*models.Book, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetBookByISBN10").Logger()
	book := &models.Book{}

	err := r.tx.Get(book, "SELECT * FROM book WHERE isbn10 = $1", isbn10)

	if err != nil {
		l.Err(err).Str("isbn10", isbn10).Msg("failed to get book")
		return nil, err
	}

	return book, nil
}

func (r *repo) CreateBook(book models.Book) (string, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "CreateBook").Logger()

	var id string

	err := r.tx.QueryRow("INSERT INTO book (title, isbn13, isbn10, publication_year, edition, publisher_id, price, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", book.Title, book.ISBN13, book.ISBN10, book.PublicationYear, book.Edition, book.PublisherID, book.Price, book.ImageURL).Scan(&id)

	if err != nil {
		l.Err(err).Msg("failed to create book")
		return "", err
	}

	return id, nil
}

func (r *repo) UpdateBook(book models.Book) error {
	l := r.logger.With().Str("package", packageName).Str("function", "UpdateBook").Logger()

	_, err := r.tx.NamedExec("UPDATE book SET title = :title, isbn13 = :isbn13, isbn10 = :isbn10, publication_year = :publication_year, edition = :edition, publisher_id = :publisher_id, price = :price, image_url = :image_url WHERE id = :id", book)

	if err != nil {
		l.Err(err).Msg("failed to update book")
		return err
	}

	return nil
}

func (r *repo) DeleteBook(id string) error {
	l := r.logger.With().Str("package", packageName).Str("function", "DeleteBook").Logger()

	_, err := r.tx.Exec("DELETE FROM book WHERE id = $1", id)

	if err != nil {
		l.Err(err).Msg("failed to delete book")
		return err
	}

	return nil
}
