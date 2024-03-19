package repos

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	packageName    = "repos"
	RecordNotFound = "record not found"
)

type Employee struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"dob"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
	Department  string `json:"department"`
	Role        string `json:"role"`
}

type EmployeeRepo interface {
	GetEmployee(id string) (*Employee, error)
	GetEmployees() ([]*Employee, error)
	CreateEmployee(emp *Employee) (*Employee, error)
	DeleteEmployee(id string) error
	UpdateEmployee(emp *Employee) (*Employee, error)
}

type employeeRepo struct {
	logger  zerolog.Logger
	empData map[string]*Employee
	mu      sync.RWMutex
}

// NewEmployeeRepo creates a new employee repository
func NewEmployeeRepo(logger zerolog.Logger, empData *map[string]*Employee) EmployeeRepo {
	repo := &employeeRepo{
		logger: logger,
	}

	if empData != nil {
		repo.empData = *empData
	} else {
		repo.empData = make(map[string]*Employee)
	}

	return repo
}

// GetEmployee gets an employee by id
func (e *employeeRepo) GetEmployee(id string) (*Employee, error) {
	// Ensure we read once it's safe to do so
	e.mu.RLock()
	defer e.mu.RUnlock()

	if emp, ok := e.empData[id]; ok {
		return emp, nil
	}

	err := errors.New(RecordNotFound)

	e.logger.Error().Err(err).Msg("failed to get employee")

	return nil, err
}

// UpdateEmployee updates an employee
func (e *employeeRepo) UpdateEmployee(emp *Employee) (*Employee, error) {
	l := e.logger.With().Str("package", packageName).Str("func", "UpdateEmployee").Logger()

	// Ensure only one write at a time
	e.mu.Lock()
	defer e.mu.Unlock()

	// Update the employee once found
	for _, empData := range e.empData {
		if empData.ID == emp.ID {
			empData.FirstName = emp.FirstName
			empData.LastName = emp.LastName
			empData.DateOfBirth = emp.DateOfBirth
			empData.Email = emp.Email
			empData.IsActive = emp.IsActive
			empData.Department = emp.Department
			empData.Role = emp.Role

			return empData, nil
		}
	}

	// If we get here, the employee was not found
	err := errors.New(RecordNotFound)

	l.Error().Err(err).Msg("failed to update employee")

	return nil, err
}

func (e *employeeRepo) GetEmployees() ([]*Employee, error) {
	// Ensure we read once it's safe to do so
	e.mu.RLock()
	defer e.mu.RUnlock()

	employees := make([]*Employee, 0)

	for _, emp := range e.empData {
		employees = append(employees, emp)
	}

	return employees, nil
}

func (e *employeeRepo) CreateEmployee(emp *Employee) (*Employee, error) {
	// Ensure only one write at a time
	e.mu.Lock()
	defer e.mu.Unlock()

	emp.ID = uuid.New().String()
	e.empData[emp.ID] = emp

	return emp, nil
}

func (e *employeeRepo) DeleteEmployee(id string) error {
	l := e.logger.With().Str("package", packageName).Str("func", "DeleteEmployee").Logger()
	e.mu.Lock()
	defer e.mu.Unlock()

	// Delete the employee if found
	if _, ok := e.empData[id]; ok {
		delete(e.empData, id)
		return nil
	}

	// If we get here, the employee was not found
	err := errors.New(RecordNotFound)

	l.Error().Err(err).Msg("failed to delete employee")

	return err
}
