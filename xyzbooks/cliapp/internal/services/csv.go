package services

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

type CSVLine struct {
	ISBN13 string
	ISBN10 string
}

type CSVProcessor interface {
	RetrieveCSVData() ([]CSVLine, error)
	AppendCSVData([]CSVLine) error
}

type csvProcessor struct {
	logger  zerolog.Logger
	fileSvc FileService
}

func NewCSVProcessor(fileSvc FileService, logger zerolog.Logger) CSVProcessor {
	return &csvProcessor{
		fileSvc: fileSvc,
		logger:  logger,
	}
}

func (c *csvProcessor) RetrieveCSVData() ([]CSVLine, error) {
	l := c.logger.With().Str("package", packageName).Str("function", "RetrieveCSVData").Logger()

	data, err := c.fileSvc.ReadData()
	if err != nil {
		l.Err(err).Str("function", "RetrieveCSVData").Msg("Failed to read data from file")
		return nil, err
	}

	csvLines := []CSVLine{}
	for i, d := range data {
		line := strings.Split(d, ",")

		if len(line) == 0 {
			l.Info().Msgf("Skipping invalid line no. %d", i+1)
			continue
		}

		csvLines = append(csvLines, CSVLine{
			ISBN13: line[0],
		})
	}

	return csvLines, nil
}

func (c *csvProcessor) AppendCSVData(lines []CSVLine) error {
	data := []string{}
	for _, line := range lines {
		data = append(data, fmt.Sprintf("%s,%s", line.ISBN13, line.ISBN10))
	}

	return c.fileSvc.AppendData(strings.Join(data, "\n"))
}
