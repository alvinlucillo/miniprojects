package services

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type FileService interface {
	ReadData() ([]string, error)
	AppendData(data string) error
	Close() error
}

type FileServiceConfig struct {
	FilePath string
	Logger   zerolog.Logger
	FileType string
}

type fileSvc struct {
	file   *os.File
	logger zerolog.Logger
}

func NewFileService(cfg FileServiceConfig) (FileService, error) {
	l := cfg.Logger.With().Str("package", packageName).Str("function", "NewFileService").Logger()

	// Open file for appending
	file, err := os.OpenFile(cfg.FilePath, os.O_RDWR, 0644)
	if err != nil {
		l.Err(err).Str("function", "NewFileService").Msg("Failed to open file")
		return nil, err
	}

	// Check if it's a csv file
	if !strings.HasSuffix(cfg.FilePath, cfg.FileType) {
		l.Error().Str("function", "NewFileService").Msg("File is not a csv file")
		return nil, errors.New("file is not a csv file")
	}

	return &fileSvc{logger: cfg.Logger, file: file}, nil
}

func (fs *fileSvc) ReadData() ([]string, error) {
	l := fs.logger.With().Str("package", packageName).Str("function", "ReadData").Logger()

	// Seek to the start of the file
	_, err := fs.file.Seek(0, 0)
	if err != nil {
		l.Err(err).Str("function", "ReadData").Msg("Failed to seek to the start of the file")
		return nil, err
	}

	// Read the whole file
	data, err := io.ReadAll(fs.file)
	if err != nil {
		l.Err(err).Str("function", "ReadData").Msg("Failed to read file")
		return nil, err
	}

	// Split the file content into lines
	lines := strings.Split(string(data), "\n")

	return lines, nil
}

func (fs *fileSvc) AppendData(data string) error {
	l := fs.logger.With().Str("package", packageName).Str("function", "AppendData").Logger()

	// Check if the file is empty
	stat, err := fs.file.Stat()
	if err != nil {
		l.Err(err).Str("function", "AppendData").Msg("Failed to get file stats")
		return err
	}

	// If the file is not empty, check if it ends with a newline
	if stat.Size() > 0 {
		buf := make([]byte, 1)
		ret, _ := fs.file.Seek(-1, io.SeekCurrent)
		_, err := fs.file.ReadAt(buf, ret)
		if err != nil && err != io.EOF {
			l.Err(err).Str("function", "AppendData").Msg("Failed to read the last byte of the file")
			return err
		}

		// If the file doesn't end with a newline, add one
		if buf[0] != '\n' {
			_, err = fs.file.WriteString("\n")
			if err != nil {
				l.Err(err).Str("function", "AppendData").Msg("Failed to append newline to file")
				return err
			}
		}

		// Move the file pointer back to the end of the file
		_, err = fs.file.Seek(0, io.SeekEnd)
		if err != nil {
			l.Err(err).Str("function", "AppendData").Msg("Failed to seek to the end of the file")
			return err
		}
	}

	// Append the new data
	_, err = fs.file.WriteString(data + "\n")
	if err != nil {
		l.Err(err).Str("function", "AppendData").Msg("Failed to append data to file")
		return err
	}

	return nil
}

func (fs *fileSvc) Close() error {
	return fs.file.Close()
}
