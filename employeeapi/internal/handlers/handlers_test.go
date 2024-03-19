package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"employeeapi/internal/repos"
	"employeeapi/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var logger = zerolog.New(os.Stdout)

const (
	employeeId1 = "97253edd-6fdf-476b-930d-ee8e3ccab0e3"
	employeeId2 = "104d0e21-81d1-4adc-888b-e0f75bbc7e6a"
)

// Test GetEmployees handler
func TestGetEmployees(t *testing.T) {
	testCases := []struct {
		name       string
		opts       string
		httpStatus int
		err        error
	}{
		{
			name:       "Successful - Get Employees - single",
			opts:       "single",
			httpStatus: http.StatusOK,
			err:        nil,
		},
		{
			name:       "Successful - Get Employees - multiple",
			opts:       "mult",
			httpStatus: http.StatusOK,
			err:        nil,
		},
		{
			name:       "Successful - Get Employees - empty list",
			httpStatus: http.StatusOK,
			err:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, data := mockRepo(tc.opts)

			svc := services.NewService(logger, false)
			svc.EmpRepo = repo

			h := NewHandler(logger, svc)

			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", "/employees", nil)
			require.NoError(t, err, "failed to create request")

			// Create a response recorder to record the response from the handler
			rr := httptest.NewRecorder()

			// Create a gin context with the request and response recorder
			c, r := gin.CreateTestContext(rr)
			c.Request = req

			// Set up routes
			h.SetupRoutes(r)

			// Call the handler
			// h.GetEmployees(c)
			r.ServeHTTP(rr, req)

			status := rr.Code

			require.Equal(t, tc.httpStatus, status)

			// Check the response body is what we expect
			if status == http.StatusOK {
				rrBody := rr.Body.String()
				employees := []Employee{}

				err = json.Unmarshal([]byte(rrBody), &employees)
				require.NoError(t, err, "failed to unmarshal response body")

				require.Equal(t, len(*data), len(employees), "number of employees returned does not match")

				for _, emp := range employees {
					e, ok := (*data)[emp.ID]
					require.True(t, ok, "employee not found in data")

					require.Equal(t, e.ID, emp.ID)
					require.Equal(t, e.FirstName, emp.FirstName)
					require.Equal(t, e.LastName, emp.LastName)
					require.Equal(t, e.DateOfBirth, emp.DateOfBirth)
					require.Equal(t, e.Email, emp.Email)
					require.Equal(t, e.IsActive, emp.IsActive)
					require.Equal(t, e.Department, emp.Department)
					require.Equal(t, e.Role, emp.Role)
				}
			}
		})
	}
}

// Test GetEmployee handler
func TestGetEmployee(t *testing.T) {
	testCases := []struct {
		name       string
		id         string
		opts       string
		httpStatus int
		err        error
	}{
		{
			name:       "Successful - Get Employee",
			id:         employeeId1,
			opts:       "single",
			httpStatus: http.StatusOK,
			err:        nil,
		},
		{
			name:       "Failed - Get Employee - not found",
			id:         uuid.New().String(),
			httpStatus: http.StatusNotFound,
			err:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, data := mockRepo(tc.opts)

			svc := services.NewService(logger, false)
			svc.EmpRepo = repo

			h := NewHandler(logger, svc)

			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", "/employees/"+tc.id, nil)
			require.NoError(t, err, "failed to create request")

			// Create a response recorder to record the response from the handler
			rr := httptest.NewRecorder()

			// Create a gin context with the request and response recorder
			c, r := gin.CreateTestContext(rr)
			c.Request = req

			// Set up routes
			h.SetupRoutes(r)

			// Call the handler
			r.ServeHTTP(rr, req)

			status := rr.Code

			require.Equal(t, tc.httpStatus, status)

			// Check the response body is what we expect
			if status == http.StatusOK {
				rrBody := rr.Body.String()
				emp := Employee{}

				err = json.Unmarshal([]byte(rrBody), &emp)
				require.NoError(t, err, "failed to unmarshal response body")

				require.NotNil(t, emp, "employee is nil")

				e, ok := (*data)[employeeId1]
				require.True(t, ok, "employee not found in data")

				require.Equal(t, e.ID, emp.ID)
				require.Equal(t, e.FirstName, emp.FirstName)
				require.Equal(t, e.LastName, emp.LastName)
				require.Equal(t, e.DateOfBirth, emp.DateOfBirth)
				require.Equal(t, e.Email, emp.Email)
				require.Equal(t, e.IsActive, emp.IsActive)
				require.Equal(t, e.Department, emp.Department)
				require.Equal(t, e.Role, emp.Role)
			}
		})
	}
}

// Test DeleteEmployee handler
func TestDeleteEmployee(t *testing.T) {
	testCases := []struct {
		name       string
		id         string
		opts       string
		httpStatus int
		err        error
	}{
		{
			name:       "Successful - Delete Employee",
			id:         employeeId1,
			opts:       "single",
			httpStatus: http.StatusOK,
			err:        nil,
		},
		{
			name:       "Failed - Delete Employee - not found",
			id:         uuid.New().String(),
			httpStatus: http.StatusNotFound,
			err:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, data := mockRepo(tc.opts)

			svc := services.NewService(logger, false)
			svc.EmpRepo = repo

			h := NewHandler(logger, svc)

			// Create a request to pass to our handler
			req, err := http.NewRequest("DELETE", "/employees/"+tc.id, nil)
			require.NoError(t, err, "failed to create request")

			// Create a response recorder to record the response from the handler
			rr := httptest.NewRecorder()

			// Create a gin context with the request and response recorder
			c, r := gin.CreateTestContext(rr)
			c.Request = req

			// Set up routes
			h.SetupRoutes(r)

			// Call the handler
			r.ServeHTTP(rr, req)

			status := rr.Code

			require.Equal(t, tc.httpStatus, status)

			if status == http.StatusOK {
				_, ok := (*data)[employeeId1]
				require.False(t, ok, "employee not deleted")
			}
		})
	}
}

// Test CreateEmployee handler
func TestCreateEmployee(t *testing.T) {

	_, data := mockRepo("single")
	employee := (*data)[employeeId1]

	testCases := []struct {
		name       string
		employee   Employee
		httpStatus int
		errorTitle string
		error      ValidationError
	}{
		{
			name: "Successful - Create Employee",
			employee: Employee{
				FirstName:   employee.FirstName,
				LastName:    employee.LastName,
				DateOfBirth: employee.DateOfBirth,
				Email:       employee.Email,
				IsActive:    employee.IsActive,
				Department:  employee.Department,
				Role:        employee.Role,
			},
			httpStatus: http.StatusCreated,
		},
		{
			name: "Successful - Create Employee - required fields only",
			employee: Employee{
				FirstName:   employee.FirstName,
				LastName:    employee.LastName,
				DateOfBirth: employee.DateOfBirth,
				Email:       employee.Email,
			},
			httpStatus: http.StatusCreated,
		},
		{
			name:       "Failed - Create Employee - empty body",
			httpStatus: http.StatusBadRequest,
			employee:   Employee{},
			error: ValidationError{
				Code:  http.StatusBadRequest,
				Title: http.StatusText(http.StatusBadRequest),
				Errors: []ValidationErrorDtl{
					{
						Field:   "FirstName",
						Message: "Key: 'Employee.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag",
					},
					{
						Field:   "LastName",
						Message: "Key: 'Employee.LastName' Error:Field validation for 'LastName' failed on the 'required' tag",
					},
					{
						Field:   "DateOfBirth",
						Message: "Key: 'Employee.DateOfBirth' Error:Field validation for 'DateOfBirth' failed on the 'required' tag",
					},
					{
						Field:   "Email",
						Message: "Key: 'Employee.Email' Error:Field validation for 'Email' failed on the 'required' tag",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, _ := mockRepo()

			svc := services.NewService(logger, false)
			svc.EmpRepo = repo

			h := NewHandler(logger, svc)

			body, err := json.Marshal(tc.employee)
			require.NoError(t, err, "failed to marshal JSON")

			// Create a request to pass to our handler
			req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(body))
			require.NoError(t, err, "failed to create request")

			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to record the response from the handler
			rr := httptest.NewRecorder()

			// Create a gin context with the request and response recorder
			c, r := gin.CreateTestContext(rr)
			c.Request = req

			// Set up routes
			h.SetupRoutes(r)

			// Call the handler
			r.ServeHTTP(rr, req)

			status := rr.Code

			require.Equal(t, tc.httpStatus, status)

			rrBody := rr.Body.String()

			// Check the response body is what we expect
			if status == http.StatusOK {
				emp := Employee{}

				err = json.Unmarshal([]byte(rrBody), &emp)
				require.NoError(t, err, "failed to unmarshal response body")

				require.NotNil(t, emp, "employee is nil")

				e, ok := (*data)[emp.ID]
				require.True(t, ok, "employee not found in data")

				require.Equal(t, e.ID, emp.ID)
				require.Equal(t, e.FirstName, emp.FirstName)
				require.Equal(t, e.LastName, emp.LastName)
				require.Equal(t, e.DateOfBirth, emp.DateOfBirth)
				require.Equal(t, e.Email, emp.Email)
				require.Equal(t, e.IsActive, emp.IsActive)
				require.Equal(t, e.Department, emp.Department)
				require.Equal(t, e.Role, emp.Role)
			}

			if status == http.StatusBadRequest {

				errResp := ValidationError{}
				err = json.Unmarshal([]byte(rrBody), &errResp)
				require.NoError(t, err, "failed to unmarshal response body")

				require.NotNil(t, errResp, "error response is nil")

				require.Equal(t, tc.error.Code, errResp.Code)
				require.Equal(t, tc.error.Title, errResp.Title)

				for i, errDtl := range errResp.Errors {
					require.Equal(t, tc.error.Errors[i].Field, errDtl.Field)
					require.Equal(t, tc.error.Errors[i].Message, errDtl.Message)
				}
			}
		})
	}
}

// Test UpdateEmployee handler
func TestUpdateEmployee(t *testing.T) {

	_, data := mockRepo("mult")
	emp2 := (*data)[employeeId2]

	testCases := []struct {
		name       string
		id         string
		employee   Employee
		opts       string
		httpStatus int
		errorTitle string
		error      ValidationError
	}{
		{
			name: "Successful - Update Employee",
			id:   employeeId1,
			employee: Employee{
				FirstName:   emp2.FirstName,
				LastName:    emp2.LastName,
				DateOfBirth: emp2.DateOfBirth,
				Email:       emp2.Email,
				IsActive:    emp2.IsActive,
				Department:  emp2.Department,
				Role:        emp2.Role,
			},
			opts:       "single",
			httpStatus: http.StatusOK,
		},
		{
			name: "Successful - Update Employee - required fields only",
			id:   employeeId1,
			employee: Employee{
				FirstName:   emp2.FirstName,
				LastName:    emp2.LastName,
				DateOfBirth: emp2.DateOfBirth,
				Email:       emp2.Email,
			},
			opts:       "single",
			httpStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, data := mockRepo(tc.opts)

			emp1 := (*data)[employeeId1]

			svc := services.NewService(logger, false)
			svc.EmpRepo = repo

			h := NewHandler(logger, svc)

			body, err := json.Marshal(tc.employee)
			require.NoError(t, err, "failed to marshal JSON")

			// Create a request to pass to our handler
			req, err := http.NewRequest("PUT", "/employees/"+tc.id, bytes.NewBuffer(body))
			require.NoError(t, err, "failed to create request")

			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to record the response from the handler
			rr := httptest.NewRecorder()

			// Create a gin context with the request and response recorder
			c, r := gin.CreateTestContext(rr)
			c.Request = req

			// Set up routes
			h.SetupRoutes(r)

			// Call the handler
			r.ServeHTTP(rr, req)

			status := rr.Code

			require.Equal(t, tc.httpStatus, status)

			rrBody := rr.Body.String()

			// Check the response body is what we expect
			if status == http.StatusOK {
				empResp := Employee{}

				err = json.Unmarshal([]byte(rrBody), &empResp)
				require.NoError(t, err, "failed to unmarshal response body")

				require.NotNil(t, empResp, "employee is nil")

				e, ok := (*data)[empResp.ID]
				require.True(t, ok, "employee not found in data")

				// Check repo data vs response
				require.Equal(t, e.ID, empResp.ID)
				require.Equal(t, e.FirstName, empResp.FirstName)
				require.Equal(t, e.LastName, empResp.LastName)
				require.Equal(t, e.DateOfBirth, empResp.DateOfBirth)
				require.Equal(t, e.Email, empResp.Email)
				require.Equal(t, e.IsActive, empResp.IsActive)
				require.Equal(t, e.Department, empResp.Department)
				require.Equal(t, e.Role, empResp.Role)

				// Check request body vs response
				require.Equal(t, tc.employee.FirstName, empResp.FirstName)
				require.Equal(t, tc.employee.LastName, empResp.LastName)
				require.Equal(t, tc.employee.DateOfBirth, empResp.DateOfBirth)
				require.Equal(t, tc.employee.Email, empResp.Email)

				// Check that optional fields retain their original values
				// if not passed in the request body
				if tc.employee.Role == "" {
					require.Equal(t, emp1.Department, empResp.Department)
					require.Equal(t, emp1.Role, empResp.Role)
				} else {
					require.Equal(t, tc.employee.IsActive, empResp.IsActive)
					require.Equal(t, tc.employee.Department, empResp.Department)
					require.Equal(t, tc.employee.Role, empResp.Role)
				}

			}

			if status == http.StatusBadRequest {

				errResp := ValidationError{}
				err = json.Unmarshal([]byte(rrBody), &errResp)
				require.NoError(t, err, "failed to unmarshal response body")

				require.NotNil(t, errResp, "error response is nil")

				require.Equal(t, tc.error.Code, errResp.Code)
				require.Equal(t, tc.error.Title, errResp.Title)

				for i, errDtl := range errResp.Errors {
					require.Equal(t, tc.error.Errors[i].Field, errDtl.Field)
					require.Equal(t, tc.error.Errors[i].Message, errDtl.Message)
				}
			}
		})
	}
}

func mockRepo(opts ...string) (repos.EmployeeRepo, *map[string]*repos.Employee) {

	data := make(map[string]*repos.Employee, 0)
	repo := repos.NewEmployeeRepo(logger, &data)

	emp1 := &repos.Employee{
		ID:          employeeId1,
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1985-05-15",
		Email:       "johndoe@example.com",
		IsActive:    false,
		Department:  "Engineering",
		Role:        "Software Developer",
	}

	emp2 := &repos.Employee{
		ID:          employeeId2,
		FirstName:   "Jane",
		LastName:    "Smith",
		DateOfBirth: "1990-09-22",
		Email:       "jane.smith@example.com",
		IsActive:    true,
		Department:  "Marketing",
		Role:        "Marketing Specialist",
	}

	if len(opts) > 0 {
		switch opts[0] {
		case "single":
			data[emp1.ID] = emp1
		case "mult":
			data[emp1.ID] = emp1
			data[emp2.ID] = emp2
		}
	}

	return repo, &data
}
