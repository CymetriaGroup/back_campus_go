package http

import (
	"errors"
	"hexagonal-go-backend/internal/adapters/handler/http/dto"
	"hexagonal-go-backend/internal/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	status, code, message := http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		status, code, message = 404, "USER_NOT_FOUND", err.Error()
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		status, code, message = 409, "EMAIL_ALREADY_EXISTS", err.Error()
	case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrInvalidToken):
		status, code, message = 401, "UNAUTHORIZED", "Invalid credentials or token"
	case errors.Is(err, domain.ErrForbidden):
		status, code, message = 403, "FORBIDDEN", err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		status, code, message = 422, "INVALID_INPUT", err.Error()
	}
	c.JSON(status, dto.ErrorResponse{Error: dto.ErrorBody{Code: code, Message: message}, RequestID: c.GetString("request_id")})
}
func validationError(c *gin.Context, _ error) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: "Invalid request"}, RequestID: c.GetString("request_id")})
}
