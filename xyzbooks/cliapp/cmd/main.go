package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	services "alvinlucillo/xyzbooks_cliapp/internal/services"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

const (
	packagename = "main"
)

type Environment struct {
	POLL_INTERVAL_SECS int    `envconfig:"POLL_INTERVAL_SECS" default:"5"`
	API_SERVER_URL     string `envconfig:"API_SERVER_URL" default:"http://localhost:9001/api"`
	// The number of times to retry the process before terminating
	ERROR_LIMIT int `envconfig:"ERROR_LIMIT" default:"5"`
}

type Flags struct {
	FilePath *string
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	l := logger.With().Str("package", packagename).Logger()

	var cfg Environment
	err := envconfig.Process("", &cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to process env vars")
	}

	flags := checkFlags(logger)

	httpService := services.NewHttpService(cfg.API_SERVER_URL, logger)
	fileService, err := services.NewFileService(services.FileServiceConfig{
		FilePath: *flags.FilePath,
		Logger:   logger,
		FileType: ".csv",
	})
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create file service")
	}

	defer fileService.Close()

	numWorkers := runtime.NumCPU() * 2

	processorService := services.NewProcessorService(services.ProcessorConfig{
		NumWorkers:  numWorkers,
		FileService: fileService,
		HttpService: httpService,
		Logger:      logger,
	})

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range c {
			l.Info().Msgf("Received %s, starting graceful shutdown", sig)
			os.Exit(0)
		}
	}()

	l.Info().Msgf("Starting processor service with poll interval (sec): %d", cfg.POLL_INTERVAL_SECS)
	pollCount := 0
	errorCount := 0

	for {
		pollCount++

		l.Info().Msgf("Poll count: %d", pollCount)
		err = processorService.Run()
		if err != nil {
			l.Err(err).Err(err).Msg("Failed to run processor service")
			errorCount++
		}

		if errorCount > cfg.ERROR_LIMIT {
			l.Fatal().Msgf("Reached error limit of %d, terminating...", cfg.ERROR_LIMIT)
		}

		time.Sleep(time.Duration(cfg.POLL_INTERVAL_SECS) * time.Second)
	}

}

func checkFlags(logger zerolog.Logger) Flags {
	l := logger.With().Str("package", packagename).Str("function", "checkFlags").Logger()

	flags := Flags{}
	// Flags definition
	flags.FilePath = flag.String("file", "test.csv", "Full path of the CSV file")

	// Parse the flags
	flag.Parse()

	if *flags.FilePath == "" {
		l.Fatal().Msg("Please provide a CSV file with the -file flag")
	}

	return flags
}
