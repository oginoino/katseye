package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

func NewCreatedResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusCreated, response)
}

func NewNoContentResponse(c *gin.Context, message string) {
	response := SuccessResponse{
		Code:    http.StatusNoContent,
		Message: message,
	}
	c.JSON(http.StatusNoContent, response)
}

func NewAcceptedResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusAccepted,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusAccepted, response)
}

func NewOKResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

func NewPartialContentResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusPartialContent,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusPartialContent, response)
}

func NewResetContentResponse(c *gin.Context, message string) {
	response := SuccessResponse{
		Code:    http.StatusResetContent,
		Message: message,
	}
	c.JSON(http.StatusResetContent, response)
}

func NewMultiStatusResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusMultiStatus,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusMultiStatus, response)
}

func NewAlreadyReportedResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusAlreadyReported,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusAlreadyReported, response)
}

func NewIMUsedResponse(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Code:    http.StatusIMUsed,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusIMUsed, response)
}

func NewTeapotResponse(c *gin.Context, message string) {
	response := SuccessResponse{
		Code:    http.StatusTeapot,
		Message: message,
	}
	c.JSON(http.StatusTeapot, response)
}

func NewDeleteSuccessResponse(c *gin.Context, resourceType string, id string) {
	response := SuccessResponse{
		Code:    http.StatusOK,
		Message: resourceType + " deleted successfully",
		Data: map[string]string{
			"id":      id,
			"message": resourceType + " with ID " + id + " has been deleted",
		},
	}
	c.JSON(http.StatusOK, response)
}

func NewNoResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
