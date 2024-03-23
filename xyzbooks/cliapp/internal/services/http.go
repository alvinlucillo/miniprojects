package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type httpSvc struct {
	baseURL string
	client  *http.Client
	logger  zerolog.Logger
}

type HttpService interface {
	Get(endpoint string) ([]byte, error)
	Put(endpoint string, payload []byte) ([]byte, error)
}

func NewHttpService(baseURL string, logger zerolog.Logger) HttpService {
	return &httpSvc{
		baseURL: baseURL,
		client:  &http.Client{},
		logger:  logger,
	}
}

func (hs *httpSvc) Get(endpoint string) ([]byte, error) {
	l := hs.logger.With().Str("package", packageName).Str("function", "Get").Logger()

	resp, err := hs.client.Get(hs.baseURL + endpoint)
	if err != nil {
		l.Err(err).Msg("Failed to send GET request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Err(err).Msg("Failed to read response body")
		return nil, err
	}

	return body, nil
}

func (hs *httpSvc) Put(endpoint string, payload []byte) ([]byte, error) {
	l := hs.logger.With().Str("package", packageName).Str("function", "Put").Logger()

	req, err := http.NewRequest(http.MethodPut, hs.baseURL+endpoint, bytes.NewBuffer(payload))
	if err != nil {
		l.Err(err).Msg("Failed to create PUT request")
		return nil, err
	}

	resp, err := hs.client.Do(req)
	if err != nil {
		l.Err(err).Msg("Failed to send PUT request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Err(err).Msg("Failed to read response body")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		l.Err(errors.New(string(body))).Msg("Received non-200 status code")
		return nil, fmt.Errorf("received http error: %d %s - %s", resp.StatusCode, http.StatusText(resp.StatusCode), body)
	}

	return body, nil
}
