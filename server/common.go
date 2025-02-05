package server

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
