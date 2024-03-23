package router

import (
	"alvinlucillo/xyzbooks_webapp/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// @Summary Get all publishers
// @Description Get all publishers
// @Tags publishers
// @Accept json
// @Produce json
// @Success 200 {array} Publisher
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /publishers [get]
func (rt *Router) GetPublishers(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "GetPublishers").Logger()
	publishers, err := rt.Svc.Repository.GetPublishers()
	if err != nil {
		l.Err(err).Msg("failed to get publishers")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	publisherResponse := []Publisher{}
	for _, publisher := range publishers {
		publisherResponse = append(publisherResponse, Publisher{
			ID:   publisher.ID,
			Name: publisher.Name,
		})
	}

	rt.Svc.SendResponse(w, r, publisherResponse)
}

// @Summary Get publisher by ID
// @Description Get publisher by ID
// @Tags publishers
// @Accept json
// @Produce json
// @Param id path string true "Publisher ID"
// @Success 200 {object} Publisher
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /publishers/{id} [get]
func (rt *Router) GetPublisher(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "GetPublisher").Logger()

	vars := mux.Vars(r)

	id := vars["id"]

	publisher, err := rt.Svc.Repository.GetPublisher(id)
	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, Publisher{
		ID:   publisher.ID,
		Name: publisher.Name,
	})
}

// @Summary Create publisher
// @Description Create publisher
// @Tags publishers
// @Accept json
// @Produce json
// @Param publisher body Publisher true "Publisher"
// @Success 200 {object} Publisher
// @Failure 400 {object} ValidationErrors
// @Failure 500
// @Router /publishers [post]
func (rt *Router) CreatePublisher(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "CreatePublisher").Logger()

	var publisher Publisher
	err := json.NewDecoder(r.Body).Decode(&publisher)
	if err != nil {
		l.Err(err).Msg("failed to decode publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	id, err := rt.Svc.Repository.CreatePublisher(models.Publisher{
		Name: publisher.Name,
	})
	if err != nil {
		l.Err(err).Msg("failed to create publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, Publisher{
		ID:   id,
		Name: publisher.Name,
	})
}

// @Summary Update publisher
// @Description Update publisher
// @Tags publishers
// @Accept json
// @Produce json
// @Param id path string true "Publisher ID"
// @Param publisher body Publisher true "Publisher"
// @Success 200 {object} Publisher
// @Failure 400 {object} ValidationErrors
// @Failure 404
// @Failure 500
// @Router /publishers/{id} [put]
func (rt *Router) UpdatePublisher(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "UpdatePublisher").Logger()

	vars := mux.Vars(r)
	id := vars["id"]

	var publisher Publisher
	err := json.NewDecoder(r.Body).Decode(&publisher)
	if err != nil {
		l.Err(err).Msg("failed to decode publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// check if publisher exists
	existingPublisher, err := rt.Svc.Repository.GetPublisher(id)
	if err != nil {
		l.Err(err).Msg("failed to get publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	existingPublisher.Name = publisher.Name

	err = rt.Svc.Repository.UpdatePublisher(existingPublisher)
	if err != nil {
		l.Err(err).Msg("failed to update publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, Publisher{
		ID:   id,
		Name: publisher.Name,
	})
}

// @Summary Delete publisher
// @Description Delete publisher
// @Tags publishers
// @Accept json
// @Produce json
// @Param id path string true "Publisher ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /publishers/{id} [delete]
func (rt *Router) DeletePublisher(w http.ResponseWriter, r *http.Request) {
	l := rt.Logger.With().Str("package", packageName).Str("function", "DeletePublisher").Logger()

	vars := mux.Vars(r)
	id := vars["id"]

	// check if publisher exists
	_, err := rt.Svc.Repository.GetPublisher(id)
	if err != nil {
		l.Err(err).Msg("failed to get publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	// check if publisher has books
	books, err := rt.Svc.Repository.GetBooksByPublisherID(id)
	if err != nil {
		l.Err(err).Msg("failed to get books")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	if len(books) > 0 {
		rt.Svc.SendError(w, http.StatusBadRequest, "Publisher still has books. Cannot delete publisher.")
		return
	}

	err = rt.Svc.Repository.DeletePublisher(id)
	if err != nil {
		l.Err(err).Msg("failed to delete publisher")
		rt.Svc.SendError(w, rt.getHttpCode(err), err)
		return
	}

	rt.Svc.SendResponse(w, r, nil)
}
