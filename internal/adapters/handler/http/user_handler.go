package http

import (
	"net/http"
	"strconv"

	"hexagonal-go-backend/internal/adapters/handler/http/dto"
	"hexagonal-go-backend/internal/core/domain"
	"hexagonal-go-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{ service ports.UserService }

func NewUserHandler(s ports.UserService) *UserHandler { return &UserHandler{s} }
func (h *UserHandler) Create(c *gin.Context) {
	var r dto.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		validationError(c, err)
		return
	}
	u, err := h.service.Create(c.Request.Context(), ports.CreateUserInput{Name: r.Name, Email: r.Email, Password: r.Password, Role: r.Role})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.Envelope{Data: dto.User(u)})
}
func (h *UserHandler) Get(c *gin.Context) {
	u, err := h.service.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Envelope{Data: dto.User(u)})
}
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var active *bool
	if v := c.Query("active"); v != "" {
		b, e := strconv.ParseBool(v)
		if e != nil {
			validationError(c, e)
			return
		}
		active = &b
	}
	f := ports.UserFilter{Page: page, Limit: limit, Search: c.Query("search"), Role: domain.Role(c.Query("role")), Active: active, Sort: c.DefaultQuery("sort", "created_at"), Order: c.DefaultQuery("order", "asc")}
	users, total, err := h.service.List(c.Request.Context(), f)
	if err != nil {
		writeError(c, err)
		return
	}
	out := make([]dto.UserResponse, len(users))
	for i := range users {
		out[i] = dto.User(&users[i])
	}
	c.JSON(http.StatusOK, dto.Envelope{Data: out, Meta: gin.H{"page": page, "limit": limit, "total": total}})
}
func (h *UserHandler) Update(c *gin.Context) {
	var r dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		validationError(c, err)
		return
	}
	u, err := h.service.Update(c.Request.Context(), c.Param("id"), ports.UpdateUserInput{Name: r.Name, Email: r.Email, Active: r.Active, Role: r.Role})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Envelope{Data: dto.User(u)})
}
func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
