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

// Create godoc
// @Summary      Crear usuario
// @Description  Registra un nuevo usuario en la plataforma.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateUserRequest true "Payload de creación de usuario"
// @Success      201 {object} response.Envelope{data=UserResponse}
// @Failure      400 {object} response.ErrorResponse "Petición inválida o datos malformados"
// @Failure      401 {object} response.ErrorResponse "No autenticado"
// @Failure      403 {object} response.ErrorResponse "Requiere rol de Administrador"
// @Failure      409 {object} response.ErrorResponse "El email ya está registrado"
// @Router       /users [post]
func (h *Controller) Create(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}
	user, err := h.service.Create(c.Request.Context(), application.CreateUserInput{AuthID: request.AuthID, Name: request.Name, Email: request.Email, Role: request.Role})
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Envelope{Data: userResponse(user)})
}

// Get godoc
// @Summary      Obtener usuario
// @Description  Retorna la información de perfil de un usuario por su ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del usuario"
// @Success      200 {object} response.Envelope{data=UserResponse}
// @Failure      401 {object} response.ErrorResponse "No autenticado"
// @Failure      404 {object} response.ErrorResponse "Usuario no encontrado"
// @Router       /users/{id} [get]
func (h *Controller) Get(c *gin.Context) {
	user, err := h.service.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Data: userResponse(user)})
}

// List godoc
// @Summary      Listar usuarios
// @Description  Retorna una lista paginada de usuarios con filtros opcionales.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query int    false "Número de página" default(1)
// @Param        limit  query int    false "Límite por página" default(20)
// @Param        search query string false "Término de búsqueda"
// @Param        role   query string false "Filtrar por rol" Enums(admin, user)
// @Param        active query bool   false "Filtrar por estado activo"
// @Success      200 {object} response.Envelope{data=[]UserResponse,meta=object}
// @Failure      401 {object} response.ErrorResponse "No autenticado"
// @Failure      403 {object} response.ErrorResponse "Requiere rol de Administrador"
// @Router       /users [get]
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

// Update godoc
// @Summary      Actualizar usuario
// @Description  Actualiza la información de un usuario existente.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del usuario"
// @Param        request body UpdateUserRequest true "Datos a actualizar del usuario"
// @Success      200 {object} response.Envelope{data=UserResponse}
// @Failure      400 {object} response.ErrorResponse "Petición inválida"
// @Failure      401 {object} response.ErrorResponse "No autenticado"
// @Failure      403 {object} response.ErrorResponse "Requiere rol de Administrador"
// @Failure      404 {object} response.ErrorResponse "Usuario no encontrado"
// @Router       /users/{id} [put]
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

// Delete godoc
// @Summary      Eliminar usuario
// @Description  Elimina un usuario por su ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del usuario"
// @Success      204 "Sin contenido"
// @Failure      401 {object} response.ErrorResponse "No autenticado"
// @Failure      403 {object} response.ErrorResponse "Requiere rol de Administrador"
// @Failure      404 {object} response.ErrorResponse "Usuario no encontrado"
// @Router       /users/{id} [delete]
func (h *Controller) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
