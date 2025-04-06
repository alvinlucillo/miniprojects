package job

import (
	"time"

	"github.com/rs/zerolog"
)

type JobService struct {
	logger zerolog.Logger
}

func NewJobService(logger zerolog.Logger) JobService {
	return JobService{
		logger: logger,
	}

}

func (j JobService) Run() {
	// Implement the logic for running the job
	j.poll()
}

func (j JobService) poll() {
	for {
		time.Sleep(time.Second * 10)
		j.logger.Println("Polling for new jobs...")
	}
}
