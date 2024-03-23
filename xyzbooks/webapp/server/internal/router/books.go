package router

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"alvinlucillo/xyzbooks_webapp/internal/models"
	"alvinlucillo/xyzbooks_webapp/internal/utils"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// @Summary Get books
// @Description Get all books
// @Tags books
// @ID get-books
// @Produce  json
// @Success 200 {array} Book
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /books [get]
func (rt *Router) GetBooks(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "GetBooks").Logger()

	books, err := rt.Svc.Repository.GetBooks()
	if err != nil {
		l.Err(err).Msg("failed to get books")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	bookResponse := []Book{}
	for _, book := range books {
		price, publisher, authors, err := rt.getBookDetails(book)
		if err != nil {
			l.Err(err).Msg("failed to get book details")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}

		b := Book{
			ID:              book.ID,
			Title:           book.Title,
			ISBN13:          book.ISBN13,
			ISBN10:          book.ISBN10,
			PublicationYear: book.PublicationYear,
			Edition:         book.Edition.String,
			Price:           price,
			Publisher:       publisher,
			Authors:         authors,
		}

		bookResponse = append(bookResponse, b)
	}

	rt.Svc.SendResponse(w, r, bookResponse)
}

// @Summary Get book by id
// @Description Get book by id
// @Tags books
// @ID get-book-by-id
// @Produce  json
// @Success 200 {object} Book
// @Failure 404
// @Failure 500
// @Router /books/{isbn13} [get]
func (rt *Router) GetBookByISBN13(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "GetBookByISBN13").Logger()

	vars := mux.Vars(r)
	isbn13 := vars["isbn13"]

	book, err := rt.Svc.Repository.GetBookByISBN13(isbn13)
	if err != nil {
		l.Err(err).Str("isbn13", isbn13).Msg("failed to get book")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	price, publisher, authors, err := rt.getBookDetails(*book)
	if err != nil {
		l.Err(err).Msg("failed to get book details")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	bookResponse := Book{
		ID:              book.ID,
		Title:           book.Title,
		ISBN13:          book.ISBN13,
		ISBN10:          book.ISBN10,
		PublicationYear: book.PublicationYear,
		Edition:         book.Edition.String,
		Price:           price,
		Publisher:       publisher,
		Authors:         authors,
		ImageURL:        book.ImageURL.String,
	}

	rt.Svc.SendResponse(w, r, bookResponse)
}

// @Summary Update book
// @Description Update book
// @Tags books
// @ID update-book
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Param book body Book true "Book"
// @Success 200 {object} Book
// @Failure 400 {object} ValidationErrors
// @Failure 404
// @Failure 500
// @Router /books/{id} [put]
func (rt *Router) UpdateBook(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "UpdateBook").Logger()

	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve existing book by id
	existingBook, err := rt.Svc.Repository.GetBook(id)
	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get book")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// Decode the JSON request to get the book details
	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		l.Err(err).Msg("failed to decode request body")
		rt.Svc.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = rt.jsonValidator.Struct(book)
	if err != nil {
		rt.Svc.SendError(w, http.StatusBadRequest, rt.getValidationErrors(err))
		return
	}

	// Retrieve existing book by ISBN13
	bookByISBN13, err := rt.Svc.Repository.GetBookByISBN13(book.ISBN13)
	if err != nil {
		// Skip the error if it is a SQL no rows error
		if !utils.IsSQLNoRowsErr(err) {
			l.Err(err).Str("id", id).Msg("failed to get book")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}
	}

	// If the existing book is not the same as the book retrieved by ISBN13, then the ISBN13 already exists
	if bookByISBN13 != nil && existingBook.ID != bookByISBN13.ID {
		rt.Svc.SendError(w, http.StatusBadRequest, "ISBN 13 already exists")
		return
	}

	bookByISBN10, err := rt.Svc.Repository.GetBookByISBN10(book.ISBN10)
	if err != nil {
		// Skip the error if it is an SQL no rows error
		if !utils.IsSQLNoRowsErr(err) {
			l.Err(err).Str("id", id).Msg("failed to get book")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}
	}

	// If the existing book is not the same as the book retrieved by ISBN10, then the ISBN10 already exists
	if bookByISBN10 != nil && existingBook.ID != bookByISBN10.ID {
		rt.Svc.SendError(w, http.StatusBadRequest, "ISBN 10 already exists")
		return
	}

	// Check if ISBN10 and ISBN13 are valid
	if bookByISBN10 != nil && !utils.ValidateISBN10(book.ISBN10) {
		l.Error().Msg("Invalid ISBN 10")
		rt.Svc.SendError(w, http.StatusBadRequest, "Invalid ISBN 10")
		return
	}

	if !utils.ValidateISBN13(book.ISBN13) {
		l.Error().Msg("Invalid ISBN 13")
		rt.Svc.SendError(w, http.StatusBadRequest, "Invalid ISBN 13")
		return
	}

	// Check if publisher exists
	_, err = rt.Svc.Repository.GetPublisher(book.Publisher.ID)
	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// Update the book
	existingBook.Title = book.Title
	existingBook.ISBN13 = book.ISBN13
	existingBook.ISBN10 = book.ISBN10
	existingBook.PublicationYear = book.PublicationYear
	existingBook.Edition = sql.NullString{String: book.Edition, Valid: book.Edition != ""}
	existingBook.Price = decimal.NewFromFloat(book.Price)
	existingBook.PublisherID = book.Publisher.ID
	existingBook.ImageURL = sql.NullString{String: book.ImageURL, Valid: book.ImageURL != ""}

	err = rt.Svc.Repository.UpdateBook(existingBook)
	if err != nil {
		l.Err(err).Msg("failed to update book")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// Update author and book relationship
	// Remove first all authors from the book
	err = rt.Svc.Repository.DeleteAuthorBookRelByBookID(existingBook.ID)
	if err != nil {
		l.Err(err).Msg("failed to remove authors")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// Add the new authors to the book
	for _, author := range book.Authors {
		// Check if author exists
		_, err = rt.Svc.Repository.GetAuthor(author.ID)
		if err != nil {
			l.Err(err).Str("id", id).Msg("failed to get author")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}

		err = rt.Svc.Repository.AddAuthorBookRel(existingBook.ID, author.ID)
		if err != nil {
			l.Err(err).Msg("failed to add author")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return

		}
	}

	rt.Svc.SendResponse(w, r, book)
}

// @Summary Create book
// @Description Create book
// @Tags books
// @ID create-book
// @Accept  json
// @Produce  json
// @Produce  json
// @Param id path string true "Book ID"
// @Param book body Book true "Book"
// @Success 200 {object} Book
// @Failure 400 {object} ValidationErrors
// @Failure 404
// @Failure 500
// @Router /books/{id} [post]
func (rt *Router) CreateBook(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "CreateBook").Logger()

	// Decode the JSON request to get the book details
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		l.Err(err).Msg("failed to decode request body")
		rt.Svc.SendError(w, http.StatusBadRequest, err)
		return
	}

	// Check if ISBN10 and ISBN13 are valid
	if book.ISBN10 != "" && !utils.ValidateISBN10(book.ISBN10) {
		l.Error().Msg("Invalid ISBN 10")
		rt.Svc.SendError(w, http.StatusBadRequest, "Invalid ISBN 10")
		return
	}

	if !utils.ValidateISBN13(book.ISBN13) {
		l.Error().Msg("Invalid ISBN 13")
		rt.Svc.SendError(w, http.StatusBadRequest, "Invalid ISBN 13")
		return
	}

	// Check if publisher exists
	_, err = rt.Svc.Repository.GetPublisher(book.Publisher.ID)
	if err != nil {
		l.Err(err).Msg("failed to get publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// Check if book exists
	bookISBN13, err := rt.Svc.Repository.GetBookByISBN13(book.ISBN13)
	if err != nil {
		// Skip the error if it is a SQL no rows error
		if !utils.IsSQLNoRowsErr(err) {
			l.Err(err).Msg("failed to get book")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}
	}

	// Existing book means the ISBN13 already exists
	if bookISBN13 != nil {
		rt.Svc.SendError(w, http.StatusBadRequest, "ISBN 13 already exists")
		return
	}

	// Check if book exists
	bookISBN10, err := rt.Svc.Repository.GetBookByISBN10(book.ISBN10)
	if err != nil {
		// Skip the error if it is a SQL no rows error
		if !utils.IsSQLNoRowsErr(err) {
			l.Err(err).Msg("failed to get book")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}
	}

	// Existing book means the ISBN10 already exists
	if bookISBN10 != nil {
		rt.Svc.SendError(w, http.StatusBadRequest, "ISBN 10 already exists")
		return
	}

	// Create the book
	newBook := models.Book{
		Title:           book.Title,
		ISBN13:          book.ISBN13,
		ISBN10:          book.ISBN10,
		PublicationYear: book.PublicationYear,
		Edition:         sql.NullString{String: book.Edition, Valid: book.Edition != ""},
		Price:           decimal.NewFromFloat(book.Price),
		PublisherID:     book.Publisher.ID,
		ImageURL:        sql.NullString{String: book.ImageURL, Valid: book.ImageURL != ""},
	}

	newBookID, err := rt.Svc.Repository.CreateBook(newBook)
	if err != nil {
		l.Err(err).Msg("failed to create book")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// Add the authors to the book
	for _, author := range book.Authors {
		// Check if author exists
		_, err = rt.Svc.Repository.GetAuthor(author.ID)
		if err != nil {
			l.Err(err).Msg("failed to get author")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}

		err = rt.Svc.Repository.AddAuthorBookRel(newBookID, author.ID)
		if err != nil {
			l.Err(err).Msg("failed to add author")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return

		}
	}

	rt.Svc.SendResponse(w, r, book)
}

// @Summary Delete book
// @Description Delete book
// @Tags books
// @ID delete-book
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /books/{id} [delete]
func (rt *Router) DeleteBook(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "DeleteBook").Logger()

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		rt.Svc.SendError(w, http.StatusBadRequest, InvalidIDError)
		return
	}

	// check first if book exists
	_, err := rt.Svc.Repository.GetBook(id)
	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get book")
		if utils.IsSQLNoRowsErr(err) {
			rt.Svc.SendError(w, http.StatusNotFound, "Book not found")
			return
		}
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// delete authors and book relationship
	err = rt.Svc.Repository.DeleteAuthorBookRelByBookID(id)
	if err != nil {
		// skip if no rows error
		if !utils.IsSQLNoRowsErr(err) {
			l.Err(err).Msg("failed to remove authors")
			rt.Svc.SendError(w, rt.getHttpCode(err), err)
			return
		}
	}

	// delete book
	err = rt.Svc.Repository.DeleteBook(id)
	if err != nil {
		l.Err(err).Msg("failed to delete book")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, nil)
}
