package repos

import "alvinlucillo/xyzbooks_webapp/internal/models"

func (r *repo) GetPublishers() ([]models.Publisher, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetPublishers").Logger()
	publishers := []models.Publisher{}

	err := r.tx.Select(&publishers, "SELECT * FROM publisher")

	if err != nil {
		l.Err(err).Msg("failed to get publishers")
		return nil, err
	}

	return publishers, nil
}

func (r *repo) GetPublisher(id string) (models.Publisher, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "GetPublisher").Logger()
	publisher := models.Publisher{}

	err := r.tx.Get(&publisher, "SELECT * FROM publisher WHERE id = $1", id)

	if err != nil {
		l.Err(err).Str("id", id).Msg("failed to get publisher")
		return publisher, err
	}

	return publisher, nil
}

func (r *repo) CreatePublisher(publisher models.Publisher) (string, error) {
	l := r.logger.With().Str("package", packageName).Str("function", "CreatePublisher").Logger()

	var id string
	err := r.tx.QueryRow("INSERT INTO publisher (name) VALUES ($1) RETURNING id", publisher.Name).Scan(&id)

	if err != nil {
		l.Err(err).Msg("failed to create publisher")
		return "", err
	}

	return id, nil
}

func (r *repo) UpdatePublisher(publisher models.Publisher) error {
	l := r.logger.With().Str("package", packageName).Str("function", "UpdatePublisher").Logger()

	_, err := r.tx.NamedExec("UPDATE publisher SET name = :name WHERE id = :id", publisher)

	if err != nil {
		l.Err(err).Msg("failed to update publisher")
		return err
	}

	return nil
}

func (r *repo) DeletePublisher(id string) error {
	l := r.logger.With().Str("package", packageName).Str("function", "DeletePublisher").Logger()

	_, err := r.tx.Exec("DELETE FROM publisher WHERE id = $1", id)

	if err != nil {
		l.Err(err).Msg("failed to delete publisher")
		return err
	}

	return nil
}
