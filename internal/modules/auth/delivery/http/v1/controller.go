package v1

import (
	"net/http"
	"time"

	"hexagonal-go-backend/internal/modules/auth/application"
	"hexagonal-go-backend/internal/platform/http/response"

	"github.com/gin-gonic/gin"
)

type Controller struct{ service application.AuthService }

func NewController(service application.AuthService) *Controller {
	return &Controller{service: service}
}

// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica credenciales de usuario y genera un par de tokens JWT (Access y Refresh).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Credenciales de acceso"
// @Success      200 {object} response.Envelope{data=TokenResponse}
// @Failure      400 {object} response.ErrorResponse "Formato de request inválido"
// @Failure      401 {object} response.ErrorResponse "INVALID_CREDENTIALS: email o contraseña incorrectos"
// @Failure      403 {object} response.ErrorResponse "ACCOUNT_DISABLED: cuenta deshabilitada"
// @Failure      423 {object} response.ErrorResponse "ACCOUNT_LOCKED: cuenta bloqueada temporalmente"
// @Failure      503 {object} response.ErrorResponse "IDENTITY_PROVIDER_UNAVAILABLE: servicio de autenticación no disponible"
// @Router       /auth/login [post]
func (h *Controller) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}
	pair, err := h.service.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	writeTokens(c, pair)
}

// Refresh godoc
// @Summary      Renovar token de acceso
// @Description  Genera un nuevo par de tokens JWT utilizando un refresh token válido.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body TokenRequest true "Refresh token"
// @Success      200 {object} response.Envelope{data=TokenResponse}
// @Failure      400 {object} response.ErrorResponse "Formato de request inválido"
// @Failure      401 {object} response.ErrorResponse "Token inválido o expirado"
// @Router       /auth/refresh [post]
func (h *Controller) Refresh(c *gin.Context) {
	var request TokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}
	pair, err := h.service.Refresh(c.Request.Context(), request.RefreshToken)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	writeTokens(c, pair)
}

// Logout godoc
// @Summary      Cerrar sesión
// @Description  Invalida el refresh token del usuario.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body TokenRequest true "Refresh token a revocar"
// @Success      204 "Sin contenido"
// @Failure      400 {object} response.ErrorResponse "Formato de request inválido"
// @Failure      401 {object} response.ErrorResponse "Token inválido"
// @Router       /auth/logout [post]
func (h *Controller) Logout(c *gin.Context) {
	var request TokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}
	if err := h.service.Logout(c.Request.Context(), request.RefreshToken); err != nil {
		response.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeTokens(c *gin.Context, pair application.TokenPair) {
	c.JSON(http.StatusOK, response.Envelope{Data: TokenResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken, TokenType: "Bearer", ExpiresIn: int64(time.Until(pair.AccessExpiresAt).Seconds())}})
}
