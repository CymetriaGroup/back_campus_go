package http

import (
	"hexagonal-go-backend/internal/adapters/handler/http/dto"
	"hexagonal-go-backend/internal/core/ports"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ service ports.AuthService }

func NewAuthHandler(s ports.AuthService) *AuthHandler { return &AuthHandler{s} }
func (h *AuthHandler) Login(c *gin.Context) {
	var r dto.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		validationError(c, err)
		return
	}
	p, err := h.service.Login(c.Request.Context(), r.Email, r.Password)
	if err != nil {
		writeError(c, err)
		return
	}
	tokens(c, p)
}
func (h *AuthHandler) Refresh(c *gin.Context) {
	var r dto.TokenRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		validationError(c, err)
		return
	}
	p, err := h.service.Refresh(c.Request.Context(), r.RefreshToken)
	if err != nil {
		writeError(c, err)
		return
	}
	tokens(c, p)
}
func (h *AuthHandler) Logout(c *gin.Context) {
	var r dto.TokenRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		validationError(c, err)
		return
	}
	if err := h.service.Logout(c.Request.Context(), r.RefreshToken); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func tokens(c *gin.Context, p ports.TokenPair) {
	c.JSON(http.StatusOK, dto.Envelope{Data: dto.TokenResponse{AccessToken: p.AccessToken, RefreshToken: p.RefreshToken, TokenType: "Bearer", ExpiresIn: int64(time.Until(p.AccessExpiresAt).Seconds())}})
}
