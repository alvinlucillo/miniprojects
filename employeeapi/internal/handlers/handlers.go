package handlers

import (
	"net/http"
	"time"

	"employeeapi/internal/repos"
	"employeeapi/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

const (
	packageName = "handlers"

	ErrorIDRequired  = "id is required"
	ErrorEmpNotFound = "employee not found"

	RoleUnassigned       = "Unassigned"
	DepartmentUnassigned = "Unassigned"
)

type Handler struct {
	logger        zerolog.Logger
	svc           services.Service
	jsonValidator *validator.Validate
}

type Employee struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name" validate:"required,max=100"`
	LastName    string `json:"last_name" validate:"required,max=100"`
	DateOfBirth string `json:"dob" validate:"required,dob"`
	Email       string `json:"email" validate:"required,email"`
	IsActive    bool   `json:"is_active" validate:"omitempty"`
	Department  string `json:"department" validate:"omitempty,oneof=Engineering Marketing Finance 'Human Resources' Unassigned"`
	Role        string `json:"role" validate:"omitempty,oneof='Software Developer' 'Marketing Specialist' 'Financial Analyst' 'HR Specialist' Unassigned"`
}

type ValidationError struct {
	Code   int                  `json:"code"`
	Title  string               `json:"title"`
	Errors []ValidationErrorDtl `json:"errors,omitempty"`
}

type ValidationErrorDtl struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewHandler(logger zerolog.Logger, svc services.Service) Handler {
	// Set up custom validator and register custom validation
	validator := validator.New()
	validator.RegisterValidation("dob", dateOfBirthFormat)

	return Handler{
		logger:        logger,
		svc:           svc,
		jsonValidator: validator,
	}
}

// SetupRoutes sets up the routes for the handler
func (h Handler) SetupRoutes(gin *gin.Engine) {
	gin.GET("/employees", h.GetEmployees)
	gin.GET("/employees/:id", h.GetEmployee)
	gin.POST("/employees", h.CreateEmployee)
	gin.PUT("/employees/:id", h.UpdateEmployee)
	gin.DELETE("/employees/:id", h.DeleteEmployee)
}

// GetEmployee gets an employee by id
func (h Handler) GetEmployee(c *gin.Context) {
	l := h.logger.With().Str("package", packageName).Str("func", "GetEmployee").Logger()

	id := c.Param("id")

	if id == "" {
		l.Error().Msg("id is empty")
		h.sendErrorResponse(c, http.StatusBadRequest, ErrorIDRequired)
		return
	}

	emp, err := h.svc.EmpRepo.GetEmployee(id)
	if err != nil {
		l.Error().Err(err).Msg("failed to get employee")

		if err.Error() == repos.RecordNotFound {
			h.sendErrorResponse(c, http.StatusNotFound, ErrorEmpNotFound)

		} else {
			h.sendErrorResponse(c, http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, Employee{
		ID:          emp.ID,
		FirstName:   emp.FirstName,
		LastName:    emp.LastName,
		DateOfBirth: emp.DateOfBirth,
		Email:       emp.Email,
		IsActive:    emp.IsActive,
		Department:  emp.Department,
		Role:        emp.Role,
	})

}

// GetEmployees returns all employees
func (h Handler) GetEmployees(c *gin.Context) {
	l := h.logger.With().Str("package", packageName).Str("func", "GetEmployees").Logger()

	employees, err := h.svc.EmpRepo.GetEmployees()
	if err != nil {
		l.Error().Err(err).Msg("failed to get employees")
		h.sendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	resp := make([]Employee, 0)
	for _, emp := range employees {
		resp = append(resp, Employee{
			ID:          emp.ID,
			FirstName:   emp.FirstName,
			LastName:    emp.LastName,
			DateOfBirth: emp.DateOfBirth,
			Email:       emp.Email,
			IsActive:    emp.IsActive,
			Department:  emp.Department,
			Role:        emp.Role,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// CreateEmployee creates an employee
func (h Handler) CreateEmployee(c *gin.Context) {
	l := h.logger.With().Str("package", packageName).Str("func", "CreateEmployee").Logger()

	var emp Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		l.Error().Err(err).Msg("failed to bind json")
		h.sendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	err := h.jsonValidator.Struct(emp)
	if err != nil {
		l.Error().Err(err).Msg("failed to validate json")

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			h.handleValidationErrors(c, validationErrors)
		} else {
			h.sendErrorResponse(c, http.StatusInternalServerError)
		}
		return
	}

	empToCreate := &repos.Employee{
		FirstName:   emp.FirstName,
		LastName:    emp.LastName,
		DateOfBirth: emp.DateOfBirth,
		Email:       emp.Email,
		IsActive:    emp.IsActive,
	}

	// Assign default values if not provided
	if emp.Department == "" {
		empToCreate.Department = DepartmentUnassigned
	}

	if emp.Role == "" {
		empToCreate.Role = RoleUnassigned
	}

	empRec, err := h.svc.EmpRepo.CreateEmployee(empToCreate)
	if err != nil {
		l.Error().Err(err).Msg("failed to create employee")
		h.sendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	// Assign with actual values
	emp.ID = empRec.ID
	emp.Role = empRec.Role
	emp.Department = empRec.Department

	c.JSON(http.StatusCreated, emp)
}

// UpdateEmployee updates an employee
func (h Handler) UpdateEmployee(c *gin.Context) {
	l := h.logger.With().Str("package", packageName).Str("func", "UpdateEmployee").Logger()

	id := c.Param("id")
	if id == "" {
		l.Error().Msg("id is empty")
		h.sendErrorResponse(c, http.StatusBadRequest, ErrorIDRequired)
		return
	}

	existingEmpRec, err := h.svc.EmpRepo.GetEmployee(id)
	if err != nil {
		l.Error().Err(err).Str("id", id).Msg("failed to get employee")

		if err.Error() == repos.RecordNotFound {
			h.sendErrorResponse(c, http.StatusNotFound, ErrorEmpNotFound)
		} else {
			h.sendErrorResponse(c, http.StatusInternalServerError)
		}
		return
	}

	var emp Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		l.Error().Err(err).Msg("failed to bind json")
		h.sendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	if err := h.jsonValidator.Struct(emp); err != nil {
		l.Error().Err(err).Msg("failed to validate json")

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			h.handleValidationErrors(c, validationErrors)
		} else {
			h.sendErrorResponse(c, http.StatusInternalServerError)

		}
		return
	}

	existingEmpRec.FirstName = emp.FirstName
	existingEmpRec.LastName = emp.LastName
	existingEmpRec.DateOfBirth = emp.DateOfBirth
	existingEmpRec.Email = emp.Email
	existingEmpRec.IsActive = emp.IsActive

	// We don't want to wipe these fields if they are not passed in the request body
	if emp.Department != "" {
		existingEmpRec.Department = emp.Department
	}

	if emp.Role != "" {
		existingEmpRec.Role = emp.Role
	}

	empRec, err := h.svc.EmpRepo.UpdateEmployee(existingEmpRec)
	if err != nil {
		l.Error().Err(err).Str("id", id).Msg("failed to update employee")
		h.sendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	// Discard what was passed in the request body
	emp.ID = empRec.ID
	emp.Role = empRec.Role
	emp.Department = empRec.Department

	c.JSON(http.StatusOK, emp)
}

// DeleteEmployee deletes an employee
func (h Handler) DeleteEmployee(c *gin.Context) {
	l := h.logger.With().Str("package", packageName).Str("func", "DeleteEmployee").Logger()

	// Check if id is provided
	id := c.Param("id")
	if id == "" {
		l.Error().Msg("id is empty")
		h.sendErrorResponse(c, http.StatusBadRequest, ErrorIDRequired)
		return
	}

	// Check if employee exists
	_, err := h.svc.EmpRepo.GetEmployee(id)
	if err != nil {
		l.Error().Err(err).Str("id", id).Msg("failed to get employee")

		if err.Error() == repos.RecordNotFound {
			h.sendErrorResponse(c, http.StatusNotFound, ErrorEmpNotFound)
		} else {
			h.sendErrorResponse(c, http.StatusInternalServerError)
		}
		return
	}

	// Delete employee
	if err := h.svc.EmpRepo.DeleteEmployee(id); err != nil {
		l.Error().Err(err).Msg("failed to delete employee")
		h.sendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h Handler) handleValidationErrors(c *gin.Context, err validator.ValidationErrors) {
	errDtls := make([]ValidationErrorDtl, 0)
	for _, e := range err {
		errDtls = append(errDtls, ValidationErrorDtl{
			Field:   e.Field(),
			Message: e.Error(),
		})
	}

	errResponse := &ValidationError{
		Code:   http.StatusBadRequest,
		Title:  http.StatusText(http.StatusBadRequest),
		Errors: errDtls,
	}

	c.JSON(http.StatusBadRequest, errResponse)
}

// Send error response based on http status code and title if any
func (h Handler) sendErrorResponse(c *gin.Context, code int, title ...string) {

	errorTitle := http.StatusText(code)
	if len(title) > 0 {
		errorTitle = title[0]
	}
	c.JSON(code, &ValidationError{
		Code:  code,
		Title: errorTitle,
	})
}

// Custom validation for date of birth
func dateOfBirthFormat(fl validator.FieldLevel) bool {
	dateStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return false
	}
	str := date.Format(layout)
	parsed, _ := time.Parse(layout, str)
	return parsed == date
}
