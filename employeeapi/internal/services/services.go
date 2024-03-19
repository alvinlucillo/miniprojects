package services

import (
	"employeeapi/internal/repos"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	EmpRepo repos.EmployeeRepo
}

// NewService creates a new service
// Service layer wraps the repository layer and performs any additional process
func NewService(logger zerolog.Logger, initData bool) Service {
	empRepo := repos.NewEmployeeRepo(logger, nil)

	if initData {
		_, err := empRepo.CreateEmployee(&repos.Employee{
			ID:          uuid.New().String(),
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: "1985-05-15",
			Email:       "john.doe@example.com",
			IsActive:    true,
			Department:  "Engineering",
			Role:        "Software Developer",
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to create employee")
		}

		_, err = empRepo.CreateEmployee(&repos.Employee{
			ID:          uuid.New().String(),
			FirstName:   "Jane",
			LastName:    "Smith",
			DateOfBirth: "1990-09-22",
			Email:       "jane.smith@example.com",
			IsActive:    true,
			Department:  "Marketing",
			Role:        "Marketing Specialist",
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to create employee")
		}

		_, err = empRepo.CreateEmployee(&repos.Employee{
			ID:          uuid.New().String(),
			FirstName:   "Robert",
			LastName:    "Johnson",
			DateOfBirth: "1988-03-10",
			Email:       "robert.johnson@example.com",
			IsActive:    false,
			Department:  "Finance",
			Role:        "Financial Analyst",
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to create employee")
		}

		_, err = empRepo.CreateEmployee(&repos.Employee{
			ID:          uuid.New().String(),
			FirstName:   "Emily",
			LastName:    "Williams",
			DateOfBirth: "1995-12-08",
			Email:       "emily.williams@example.com",
			IsActive:    true,
			Department:  "Human Resources",
			Role:        "HR Specialist",
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to create employee")
		}
	}

	return Service{
		EmpRepo: empRepo,
	}
}
