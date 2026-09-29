package v1

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"hexagonal-go-backend/internal/modules/courses/domain"
	"hexagonal-go-backend/internal/platform/http/response"

	"github.com/gin-gonic/gin"
)

// Service is satisfied by the courses application service.
type Service interface {
	CreateCourse(context.Context, domain.Course) (domain.Course, error)
	GetCourse(context.Context, string) (domain.Course, error)
	CreateCategory(context.Context, domain.Category) (domain.Category, error)
	ListCategories(context.Context, string) ([]domain.Category, error)
	CreateTemplate(context.Context, domain.Template) (domain.Template, error)
	ListTemplates(context.Context) ([]domain.Template, error)
	GetTemplate(context.Context, string) (domain.Template, error)
	CreateVersion(context.Context, domain.Version) (domain.Version, error)
	ListVersions(context.Context, string) ([]domain.Version, error)
	SetSyllabus(context.Context, domain.Syllabus) (domain.Syllabus, error)
	GetSyllabus(context.Context, string) (domain.Syllabus, error)
}

type Controller struct{ service Service }

func NewController(service Service) *Controller { return &Controller{service: service} }

func writeError(c *gin.Context, err error) {
	status, code, message := http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
	switch {
	case errors.Is(err, domain.ErrInvalid):
		status, code, message = http.StatusUnprocessableEntity, "INVALID_COURSE_DATA", err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, code, message = http.StatusNotFound, "COURSE_ITEM_NOT_FOUND", err.Error()
	case errors.Is(err, domain.ErrConflict):
		status, code, message = http.StatusConflict, "COURSE_ITEM_CONFLICT", err.Error()
	}
	c.JSON(status, response.ErrorResponse{Error: response.ErrorBody{Code: code, Message: message}, RequestID: c.GetString("request_id")})
}

func requiredParam(c *gin.Context, name string) string {
	value := strings.TrimSpace(c.Param(name))
	if value == "" {
		response.ValidationError(c, errors.New("missing "+name))
	}
	return value
}

// CreateCourse godoc
// @Summary Create course with its full content tree
// @Tags courses
// @Accept json
// @Produce json
// @Param body body CreateCourseRequest true "Course"
// @Success 201 {object} response.Envelope
// @Failure 400,409,422,500 {object} response.ErrorResponse
// @Router /api/v1/courses [post]
func (h *Controller) CreateCourse(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	result, err := h.service.CreateCourse(c.Request.Context(), req.toDomain())
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Envelope{Data: courseResponse(result)})
}

// GetCourse godoc
// @Summary Get course by version ID with its full content tree
// @Tags courses
// @Produce json
// @Param id path string true "Version ID"
// @Success 200 {object} response.Envelope
// @Failure 400,404,500 {object} response.ErrorResponse
// @Router /api/v1/courses/{id} [get]
func (h *Controller) GetCourse(c *gin.Context) {
	id := requiredParam(c, "id")
	if id == "" {
		return
	}
	result, err := h.service.GetCourse(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Data: courseResponse(result)})
}

// CreateCategory godoc
// @Summary Create course category
// @Tags courses
// @Accept json
// @Produce json
// @Param body body CreateCategoryRequest true "Category"
// @Success 201 {object} response.Envelope
// @Failure 400,409,422,500 {object} response.ErrorResponse
// @Router /api/v1/courses/categories [post]
func (h *Controller) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	result, err := h.service.CreateCategory(c.Request.Context(), domain.Category{TenantID: req.TenantID, Code: req.Code, Name: req.Name, Description: req.Description, ParentID: req.ParentID})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Envelope{Data: categoryResponse(result)})
}

// ListCategories godoc
// @Summary List course categories
// @Tags courses
// @Produce json
// @Param tenant_id query string true "Tenant ID"
// @Success 200 {object} response.Envelope
// @Failure 400,500 {object} response.ErrorResponse
// @Router /api/v1/courses/categories [get]
func (h *Controller) ListCategories(c *gin.Context) {
	tenantID := strings.TrimSpace(c.Query("tenant_id"))
	if tenantID == "" {
		response.ValidationError(c, errors.New("missing tenant_id"))
		return
	}
	items, err := h.service.ListCategories(c.Request.Context(), tenantID)
	if err != nil {
		writeError(c, err)
		return
	}
	output := make([]CategoryResponse, len(items))
	for i, item := range items {
		output[i] = categoryResponse(item)
	}
	c.JSON(http.StatusOK, response.Envelope{Data: output})
}

// CreateTemplate godoc
// @Summary Create course template
// @Tags courses
// @Accept json
// @Produce json
// @Param body body CreateTemplateRequest true "Template"
// @Success 201 {object} response.Envelope
// @Failure 400,409,422,500 {object} response.ErrorResponse
// @Router /api/v1/courses/templates [post]
func (h *Controller) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	result, err := h.service.CreateTemplate(c.Request.Context(), domain.Template{Code: req.Code, Title: req.Title, Description: req.Description})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Envelope{Data: templateResponse(result)})
}

// ListTemplates godoc
// @Summary List course templates
// @Tags courses
// @Produce json
// @Success 200 {object} response.Envelope
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/courses/templates [get]
func (h *Controller) ListTemplates(c *gin.Context) {
	items, err := h.service.ListTemplates(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	output := make([]TemplateResponse, len(items))
	for i, item := range items {
		output[i] = templateResponse(item)
	}
	c.JSON(http.StatusOK, response.Envelope{Data: output})
}

// GetTemplate godoc
// @Summary Get course template
// @Tags courses
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} response.Envelope
// @Failure 400,404,500 {object} response.ErrorResponse
// @Router /api/v1/courses/templates/{id} [get]
func (h *Controller) GetTemplate(c *gin.Context) {
	id := requiredParam(c, "id")
	if id == "" {
		return
	}
	result, err := h.service.GetTemplate(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Data: templateResponse(result)})
}

// CreateVersion godoc
// @Summary Create course version
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param body body CreateVersionRequest true "Version"
// @Success 201 {object} response.Envelope
// @Failure 400,404,409,422,500 {object} response.ErrorResponse
// @Router /api/v1/courses/templates/{id}/versions [post]
func (h *Controller) CreateVersion(c *gin.Context) {
	id := requiredParam(c, "id")
	if id == "" {
		return
	}
	var req CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	result, err := h.service.CreateVersion(c.Request.Context(), domain.Version{TemplateID: id, Tag: req.Tag, Status: req.Status, EstimatedHours: req.EstimatedHours})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Envelope{Data: versionResponse(result)})
}

// ListVersions godoc
// @Summary List course versions
// @Tags courses
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} response.Envelope
// @Failure 400,404,500 {object} response.ErrorResponse
// @Router /api/v1/courses/templates/{id}/versions [get]
func (h *Controller) ListVersions(c *gin.Context) {
	id := requiredParam(c, "id")
	if id == "" {
		return
	}
	items, err := h.service.ListVersions(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	output := make([]VersionResponse, len(items))
	for i, item := range items {
		output[i] = versionResponse(item)
	}
	c.JSON(http.StatusOK, response.Envelope{Data: output})
}

// SetSyllabus godoc
// @Summary Set course syllabus
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Version ID"
// @Param body body SetSyllabusRequest true "Syllabus"
// @Success 200 {object} response.Envelope
// @Failure 400,404,422,500 {object} response.ErrorResponse
// @Router /api/v1/courses/versions/{id}/syllabus [put]
func (h *Controller) SetSyllabus(c *gin.Context) {
	id := requiredParam(c, "id")
	if id == "" {
		return
	}
	var req SetSyllabusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	result, err := h.service.SetSyllabus(c.Request.Context(), domain.Syllabus{VersionID: id, Objectives: req.Objectives, EntryProfile: req.EntryProfile, ExitProfile: req.ExitProfile, Methodology: req.Methodology, DurationHours: req.DurationHours})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Data: syllabusResponse(result)})
}

// GetSyllabus godoc
// @Summary Get course syllabus
// @Tags courses
// @Produce json
// @Param id path string true "Version ID"
// @Success 200 {object} response.Envelope
// @Failure 400,404,500 {object} response.ErrorResponse
// @Router /api/v1/courses/versions/{id}/syllabus [get]
func (h *Controller) GetSyllabus(c *gin.Context) {
	id := requiredParam(c, "id")
	if id == "" {
		return
	}
	result, err := h.service.GetSyllabus(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Data: syllabusResponse(result)})
}
