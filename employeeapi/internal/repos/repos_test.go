package repos

import (
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// Tests the GetEmployee method
func TestGetEmployee(t *testing.T) {
	logger := zerolog.New(os.Stdout)

	data := make(map[string]*Employee)

	employee := &Employee{
		ID:          uuid.New().String(),
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1985-05-15",
		Email:       "johndoe@example.com",
		IsActive:    true,
		Department:  "Engineering",
		Role:        "Software Developer",
	}

	data[employee.ID] = employee

	repo := NewEmployeeRepo(logger, &data)

	testCases := []struct {
		name     string
		id       string
		employee *Employee
		err      error
	}{
		{
			name:     "Successful - Get Employee",
			id:       employee.ID,
			employee: employee,
			err:      nil,
		},
		{
			name: "Failed - Get Employee",
			err:  errors.New(RecordNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			employee, err := repo.GetEmployee(tc.id)

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.employee, employee)
		})
	}
}

// Tests the GetEmployees method
func TestGetEmployees(t *testing.T) {

	logger := zerolog.New(os.Stdout)

	data := make(map[string]*Employee)

	employee1 := &Employee{
		ID:          uuid.New().String(),
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1985-05-15",
		Email:       "johndoe@example.com",
		IsActive:    true,
		Department:  "Engineering",
		Role:        "Software Developer",
	}

	employee2 := &Employee{
		ID:          uuid.New().String(),
		FirstName:   "Jane",
		LastName:    "Smith",
		DateOfBirth: "1990-09-22",
		Email:       "janesmith@example.com",
		IsActive:    true,
		Department:  "Marketing",
		Role:        "Marketing Specialist",
	}

	data[employee1.ID] = employee1
	data[employee2.ID] = employee2

	repo := NewEmployeeRepo(logger, &data)

	testCases := []struct {
		name      string
		employees []*Employee
		err       error
	}{
		{
			name:      "Successful - Get Employees",
			employees: []*Employee{employee1, employee2},
			err:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			employees, err := repo.GetEmployees()

			require.Equal(t, tc.employees, employees)
			require.Nil(t, err)
		})
	}
}

// Tests the CreateEmployee method
func TestCreateEmployee(t *testing.T) {

	logger := zerolog.New(os.Stdout)

	data := make(map[string]*Employee)

	employee := &Employee{
		ID:          uuid.New().String(),
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1985-05-15",
		Email:       "johndoe@example.com",
		IsActive:    true,
		Department:  "Engineering",
		Role:        "Software Developer",
	}

	repo := NewEmployeeRepo(logger, &data)

	testCases := []struct {
		name     string
		employee *Employee
		err      error
	}{
		{
			name:     "Successful - Create Employee",
			employee: employee,
			err:      nil,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			employee, err := repo.CreateEmployee(tc.employee)

			require.Equal(t, tc.employee, employee)
			require.Nil(t, err)
		})
	}
}

// Tests the DeleteEmployee method
func TestDeleteEmployee(t *testing.T) {

	logger := zerolog.New(os.Stdout)

	data := make(map[string]*Employee)

	employee := &Employee{
		ID:          uuid.New().String(),
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1985-05-15",
		Email:       "johndoe@example.com",
		IsActive:    true,
		Department:  "Engineering",
		Role:        "Software Developer",
	}

	data[employee.ID] = employee

	repo := NewEmployeeRepo(logger, &data)

	testCases := []struct {
		name string
		id   string
		err  error
	}{
		{
			name: "Successful - Delete Employee",
			id:   employee.ID,
			err:  nil,
		},
		{
			name: "Failed - Delete Employee",
			id:   "123",
			err:  errors.New(RecordNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.DeleteEmployee(tc.id)

			require.Equal(t, tc.err, err)

			if tc.err == nil {
				_, ok := data[tc.id]
				require.False(t, ok)
			}
		})
	}
}

// Tests the UpdateEmployee method
func TestUpdateEmployee(t *testing.T) {

	logger := zerolog.New(os.Stdout)

	data := make(map[string]*Employee)

	employee := &Employee{
		ID:          uuid.New().String(),
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1985-05-15",
		Email:       "johndoe@exmaple.com",
		IsActive:    true,
		Department:  "Engineering",
		Role:        "Software Developer",
	}

	data[employee.ID] = employee

	repo := NewEmployeeRepo(logger, &data)

	testCases := []struct {
		name     string
		employee *Employee
		err      error
	}{
		{
			name: "Successful - Update Employee",
			employee: &Employee{
				ID:          employee.ID,
				FirstName:   "Jane",
				LastName:    "Smith",
				DateOfBirth: "1990-09-22",
				Email:       "janesmith@exmaple.com",
				IsActive:    true,
				Department:  "Marketing",
				Role:        "Marketing Specialist",
			},
			err: nil,
		},
		{
			name: "Failed - Update Employee",
			employee: &Employee{
				ID:          "123",
				FirstName:   "Jane",
				LastName:    "Smith",
				DateOfBirth: "1990-09-22",
				Email:       "janesmith@example.com",
				IsActive:    true,
				Department:  "Marketing",
				Role:        "Marketing Specialist",
			},
			err: errors.New(RecordNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			employee, err := repo.UpdateEmployee(tc.employee)

			require.Equal(t, tc.err, err)

			if tc.err == nil {
				require.Equal(t, tc.employee, employee)
				require.Equal(t, data[tc.employee.ID], employee)
			}
		})
	}
}
