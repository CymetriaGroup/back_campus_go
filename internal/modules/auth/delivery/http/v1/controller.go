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
