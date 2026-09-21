package v1

import (
	"net/http"
	"strconv"

	"hexagonal-go-backend/internal/modules/users/application"
	"hexagonal-go-backend/internal/modules/users/domain"
	"hexagonal-go-backend/internal/platform/http/response"

	"github.com/gin-gonic/gin"
)

type Controller struct{ service application.UserService }

func NewController(service application.UserService) *Controller {
	return &Controller{service: service}
}

func (h *Controller) Create(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}
	user, err := h.service.Create(c.Request.Context(), application.CreateUserInput{Name: request.Name, Email: request.Email, Password: request.Password, Role: request.Role})
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Envelope{Data: userResponse(user)})
}

func (h *Controller) Get(c *gin.Context) {
	user, err := h.service.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Data: userResponse(user)})
}

func (h *Controller) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var active *bool
	if value := c.Query("active"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			response.ValidationError(c, err)
			return
		}
		active = &parsed
	}
	filter := application.UserFilter{Page: page, Limit: limit, Search: c.Query("search"), Role: domain.Role(c.Query("role")), Active: active, Sort: c.DefaultQuery("sort", "created_at"), Order: c.DefaultQuery("order", "asc")}
	users, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	output := make([]UserResponse, len(users))
	for i := range users {
		output[i] = userResponse(&users[i])
	}
	c.JSON(http.StatusOK, response.Envelope{Data: output, Meta: gin.H{"page": page, "limit": limit, "total": total}})
}

func (h *Controller) Update(c *gin.Context) {
	var request UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}
	user, err := h.service.Update(c.Request.Context(), c.Param("id"), application.UpdateUserInput{Name: request.Name, Email: request.Email, Active: request.Active, Role: request.Role})
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Data: userResponse(user)})
}

func (h *Controller) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
