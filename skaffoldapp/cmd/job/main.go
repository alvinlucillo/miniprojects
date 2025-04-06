package main

import (
	"os"
	"skaffoldapp/internal/job"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	jobService := job.NewJobService(logger)

	jobService.Run()

}
