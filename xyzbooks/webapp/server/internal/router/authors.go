package router

import (
	"alvinlucillo/xyzbooks_webapp/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// @Summary Get all authors
// @Description Get all authors
// @Tags authors
// @Accept json
// @Produce json
// @Success 200 {array} Author
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /authors [get]
func (rt *Router) GetAuthors(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "GetAuthors").Logger()
	authors, err := rt.Svc.Repository.GetAuthors()
	if err != nil {
		l.Err(err).Msg("failed to get authors")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	authorResponse := []Author{}
	for _, author := range authors {
		authorResponse = append(authorResponse, Author{
			ID:         author.ID,
			FirstName:  author.FirstName,
			MiddleName: author.MiddleName.String,
			LastName:   author.LastName,
		})
	}

	rt.Svc.SendResponse(w, r, authorResponse)
}

// @Summary Get author by ID
// @Description Get author by ID
// @Tags authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} Author
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /authors/{id} [get]
func (rt *Router) GetAuthor(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "GetAuthor").Logger()

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		rt.Svc.SendError(w, http.StatusBadRequest, InvalidIDError)
		return
	}

	author, err := rt.Svc.Repository.GetAuthor(id)
	if err != nil {
		l.Err(err).Msg("failed to get author")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	authorResponse := Author{
		ID:         author.ID,
		FirstName:  author.FirstName,
		MiddleName: author.MiddleName.String,
		LastName:   author.LastName,
	}

	rt.Svc.SendResponse(w, r, authorResponse)
}

// @Summary Create author
// @Description Create author
// @Tags authors
// @Accept json
// @Produce json
// @Param author body Author true "Author"
// @Success 200 {object} Author
// @Failure 400 {object} ValidationErrors
// @Failure 404
// @Failure 500
// @Router /authors [post]
func (rt *Router) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "CreateAuthor").Logger()

	// Decode the JSON request to get the author details
	var author Author
	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		l.Err(err).Msg("failed to decode request body")
		rt.Svc.SendError(w, http.StatusBadRequest, err)
		return
	}

	newAuthor := models.Author{
		FirstName:  author.FirstName,
		MiddleName: sql.NullString{String: author.MiddleName, Valid: author.MiddleName != ""},
		LastName:   author.LastName,
	}

	// create author
	_, err = rt.Svc.Repository.CreateAuthor(newAuthor)
	if err != nil {
		l.Err(err).Msg("failed to create author")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, nil)
}

// @Summary Update author
// @Description Update author
// @Tags authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Param author body Author true "Author"
// @Success 200 {object} Author
// @Failure 400 {object} ValidationErrors
// @Failure 404
// @Failure 500
// @Router /authors/{id} [put]
func (rt *Router) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "UpdateAuthor").Logger()

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		rt.Svc.SendError(w, http.StatusBadRequest, InvalidIDError)
		return
	}

	// Decode the JSON request to get the author details
	var author Author
	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		l.Err(err).Msg("failed to decode request body")
		rt.Svc.SendError(w, http.StatusBadRequest, err)
		return
	}

	// check first if author exists
	existingAuthor, err := rt.Svc.Repository.GetAuthor(id)
	if err != nil {
		l.Err(err).Msg("failed to get author")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	existingAuthor.FirstName = author.FirstName
	existingAuthor.LastName = author.LastName
	existingAuthor.MiddleName = sql.NullString{String: author.MiddleName, Valid: author.MiddleName != ""}

	err = rt.Svc.Repository.UpdateAuthor(existingAuthor)
	if err != nil {
		l.Err(err).Msg("failed to update author")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, nil)
}

// @Summary Delete author
// @Description Delete author
// @Tags authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /authors/{id} [delete]
func (rt *Router) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "DeleteAuthor").Logger()

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		rt.Svc.SendError(w, http.StatusBadRequest, InvalidIDError)
		return
	}

	// check first if author exists
	_, err := rt.Svc.Repository.GetAuthor(id)
	if err != nil {
		l.Err(err).Msg("failed to get author")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// check if author has books
	authorBooks, err := rt.Svc.Repository.GetAuthorBookRelByAuthorID(id)
	if err != nil {
		l.Err(err).Msg("failed to get author book rel")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	if len(authorBooks) > 0 {
		rt.Svc.SendError(w, http.StatusBadRequest, "Author still has books. Cannot delete author.")
		return
	}

	err = rt.Svc.Repository.DeleteAuthor(id)
	if err != nil {
		l.Err(err).Msg("failed to delete author")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, nil)
}
