package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty" `
}

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    0,
		Message: message,
		Data:    data,
	})
}
func Error(c *gin.Context, httpStatus int, errorCode int, message string) {
	c.JSON(httpStatus, APIResponse{
		Success: false,
		Code:    errorCode,
		Message: message,
	})
}
