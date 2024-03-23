package repos

import "alvinlucillo/xyzbooks_webapp/internal/models"

func (r *repo) GetAuthors() ([]models.Author, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetAuthors").Logger()
	authors := []models.Author{}

	err := r.tx.Select(&authors, "SELECT * FROM author")

	if err != nil {
		l.Err(err).Msg("failed to get authors")
		return nil, err
	}

	return authors, nil
}

func (r *repo) GetAuthor(id string) (models.Author, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetAuthor").Logger()
	author := models.Author{}

	err := r.tx.Get(&author, "SELECT * FROM author WHERE id = $1", id)

	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get author")
		return author, err
	}

	return author, nil
}

func (r *repo) CreateAuthor(author models.Author) (string, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "CreateAuthor").Logger()

	var id string
	err := r.tx.QueryRow("INSERT INTO author (first_name, last_name, middle_name) VALUES ($1, $2, $3) RETURNING id", author.FirstName, author.LastName, author.MiddleName).Scan(&id)

	if err != nil {
		l.Err(err).Msg("failed to create author")
		return "", err
	}

	return id, nil
}

func (r *repo) UpdateAuthor(author models.Author) error {
	l := r.logger.With().Str("package", packageName).Str("function", "UpdateAuthor").Logger()

	_, err := r.tx.NamedExec("UPDATE author SET first_name = :first_name, last_name = :last_name, middle_name = :middle_name WHERE id = :id", author)

	if err != nil {
		l.Err(err).Msg("failed to update author")
		return err
	}

	return nil
}

func (r *repo) DeleteAuthor(id string) error {
	l := r.logger.With().Str("package", packageName).Str("function", "DeleteAuthor").Logger()

	_, err := r.tx.Exec("DELETE FROM author WHERE id = $1", id)

	if err != nil {
		l.Err(err).Msg("failed to delete author")
		return err
	}

	return nil
}

func (r *repo) GetAuthorBookRelByBookID(id string) ([]models.AuthorBook, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetAuthorBookRelByBookID").Logger()
	authorBooks := []models.AuthorBook{}

	err := r.tx.Select(&authorBooks, "SELECT * FROM authors_books where book_id = $1", id)

	if err != nil {
		l.Err(err).Msg("failed to get author_books")
		return nil, err
	}

	return authorBooks, nil
}

func (r *repo) GetAuthorBookRelByAuthorID(id string) ([]models.AuthorBook, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetAuthorBookRelByAuthorID").Logger()
	authorBooks := []models.AuthorBook{}

	err := r.tx.Select(&authorBooks, "SELECT * FROM authors_books where author_id = $1", id)

	if err != nil {
		l.Err(err).Msg("failed to get author_books")
		return nil, err
	}

	return authorBooks, nil
}

func (r *repo) AddAuthorBookRel(id string, authorId string) error {
	l := r.logger.With().Str("package", packageName).Str("function", "AddAuthorBookRel").Logger()

	_, err := r.tx.Exec("INSERT INTO authors_books (book_id, author_id) VALUES ($1, $2)", id, authorId)

	if err != nil {
		l.Err(err).Msg("failed to add authors_books")
		return err
	}

	return nil
}

func (r *repo) DeleteAuthorBookRelByBookID(id string) error {
	l := r.logger.With().Str("package", packageName).Str("function", "DeleteAuthorBookRelByBookID").Logger()

	_, err := r.tx.Exec("DELETE FROM authors_books WHERE book_id = $1", id)

	if err != nil {
		l.Err(err).Msg("failed to remove authors_books")
		return err
	}

	return nil
}

func (r *repo) DeleteAuthorBookRelByAuthorID(id string) error {
	l := r.logger.With().Str("package", packageName).Str("function", "DeleteAuthorBookRelByAuthorID").Logger()

	_, err := r.tx.Exec("DELETE FROM authors_books WHERE author_id = $1", id)

	if err != nil {
		l.Err(err).Msg("failed to remove authors_books")
		return err
	}

	return nil
}
