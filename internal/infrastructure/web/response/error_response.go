package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewErrorResponse(c *gin.Context, code int, message string, err string) {
	response := ErrorResponse{
		Code:    code,
		Message: message,
		Error:   err,
	}
	c.JSON(code, response)
}

func NewBadRequestResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusBadRequest, message, err)
}

func NewUnauthorizedResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusUnauthorized, message, err)
}

func NewForbiddenResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusForbidden, message, err)
}

func NewNotFoundResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusNotFound, message, err)
}

func NewInternalServerErrorResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusInternalServerError, message, err)
}

func NewConflictResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusConflict, message, err)
}

func NewUnprocessableEntityResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusUnprocessableEntity, message, err)
}

func NewTooManyRequestsResponse(c *gin.Context, message string, err string) {
	NewErrorResponse(c, http.StatusTooManyRequests, message, err)
}
