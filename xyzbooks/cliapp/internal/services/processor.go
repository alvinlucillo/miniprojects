package services

import (
	"alvinlucillo/xyzbooks_cliapp/internal/utils"
	"encoding/json"

	"github.com/rs/zerolog"
)

const (
	packageName = "services"
)

type Book struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	ISBN13          string    `json:"isbn13"`
	ISBN10          string    `json:"isbn10"`
	PublicationYear int       `json:"publication_year"`
	Edition         string    `json:"edition"`
	Price           float64   `json:"price"`
	ImageURL        string    `json:"image_url"`
	Publisher       Publisher `json:"publisher"`
	Authors         []Author  `json:"authors"`
}

type Author struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type Publisher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProcessorService interface {
	Run() error
}

type ProcessorConfig struct {
	NumWorkers  int
	FileService FileService
	HttpService HttpService
	Logger      zerolog.Logger
}

func NewProcessorService(cfg ProcessorConfig) ProcessorService {
	return &processorService{
		cfg: cfg,
	}

}

type BookResult struct {
	Book        Book
	ISBNUpdated bool
	Error       ErrorDtl
}

type ErrorDtl struct {
	Err error
	Msg string
}

type processorService struct {
	cfg ProcessorConfig
}

func (s *processorService) Run() error {
	l := s.cfg.Logger.With().Str("package", packageName).Str("function", "Run").Logger()

	// 1 - Retrieve existing ISBN 13 from the file
	csvSvc := NewCSVProcessor(s.cfg.FileService, s.cfg.Logger)
	csvLines, err := csvSvc.RetrieveCSVData()
	if err != nil {
		l.Err(err).Str("function", "Run").Msg("Failed to retrieve CSV data")
		return err
	}

	// Create a map of ISBN 13 for faster lookup
	isbn13Map := make(map[string]CSVLine)

	for _, line := range csvLines {
		isbn13Map[line.ISBN13] = line
	}

	// 2 - Retrieve books from the API
	books, err := s.retrieveBooks()
	if err != nil {
		l.Err(err).Msg("Failed to retrieve books")
		return err
	}

	// 3 - Send books to workers
	//  a worker will receieve a book via bookChan channel and will return a result via resultChan channel

	bookChan := make(chan Book, len(books))
	resultChan := make(chan BookResult, len(books))

	// Start workers
	for i := 0; i < s.cfg.NumWorkers; i++ {
		go func() {
			for book := range bookChan {

				var result BookResult
				result.Book = book

				// Check if ISBN 10 is missing then convert it
				if book.ISBN10 == "" {
					// Convert the ISBN 13 to ISBN 10
					book.ISBN10, err = utils.ConvertISBN13ToISBN10(book.ISBN13)
					if err != nil {
						l.Err(err).Str("isbn13", book.ISBN13).Msg("Failed to convert ISBN 13 to ISBN 10")
						result.Error = ErrorDtl{Err: err, Msg: "Failed to convert ISBN 13 to ISBN 10"}

						resultChan <- result
						continue
					}

					// Update the book via API with the new ISBN 10
					err = s.updateBook(book)
					if err != nil {
						l.Err(err).Msg("Failed to update book")
						result.Error = ErrorDtl{Err: err, Msg: "Failed to update book"}

						resultChan <- result
						continue
					}

					result.ISBNUpdated = true
				}

				result.Book = book
				resultChan <- result
			}
		}()
	}

	// Send books from the API to the book channel so workers can start processing them
	for _, book := range books {
		bookChan <- book
	}

	// Close the book channel after all books have been sent
	close(bookChan)

	// 4 - Create the resultset to be stored in the file
	newCsvLines := []CSVLine{}
	summary := struct {
		TotalISBN    int
		TotalSkipped int
		TotalAdded   int
		TotalFailed  int
	}{
		TotalISBN:    len(books),
		TotalSkipped: 0,
		TotalAdded:   0,
		TotalFailed:  0,
	}
	for i := 0; i < len(books); i++ {
		res := <-resultChan

		if res.Error.Err != nil {
			summary.TotalFailed++
		}

		// Skip if ISBN 13 already exists in the file
		if _, ok := isbn13Map[res.Book.ISBN13]; ok {
			summary.TotalSkipped++
			continue
		}

		newCsvLines = append(newCsvLines, CSVLine{
			ISBN13: res.Book.ISBN13,
			ISBN10: res.Book.ISBN10,
		})
		summary.TotalAdded++
	}

	// 5 - Append new ISBN 13 to the file
	if len(newCsvLines) > 0 {
		err = csvSvc.AppendCSVData(newCsvLines)
	}
	if err != nil {
		l.Err(err).Str("function", "Run").Msg("Failed to append CSV data")
		return err
	}

	// 6 - Log the summary
	l.Info().Msgf("Total ISBN from the API: %d", summary.TotalISBN)
	l.Info().Msgf("Total Skipped: %d", summary.TotalSkipped)
	l.Info().Msgf("Total Added: %d", summary.TotalAdded)
	l.Info().Msgf("Total Failed: %d", summary.TotalFailed)

	l.Info().Msg("Done processing books!")

	return nil
}

func (s *processorService) retrieveBooks() ([]Book, error) {
	l := s.cfg.Logger.With().Str("package", packageName).Str("function", "retrieveBooks").Logger()

	data, err := s.cfg.HttpService.Get("/books")
	if err != nil {
		l.Err(err).Msg("Failed to retrieve books")
		return nil, err
	}

	var books []Book
	err = json.Unmarshal(data, &books)
	if err != nil {
		l.Err(err).Msg("Failed to unmarshal books")
		return nil, err
	}

	return books, nil
}

func (s *processorService) updateBook(book Book) error {
	l := s.cfg.Logger.With().Str("package", packageName).Str("function", "updateBook").Logger()

	// convert book into bytes
	bookBytes, err := json.Marshal(book)
	if err != nil {
		l.Err(err).Msg("Failed to marshal book")
		return err
	}

	_, err = s.cfg.HttpService.Put("/books/"+book.ID, bookBytes)
	if err != nil {
		l.Err(err).Msg("Failed to update book via API")
		return err
	}

	l.Info().Str("isbn10", book.ISBN10).Msg("Successfully updated book")

	return nil
}
