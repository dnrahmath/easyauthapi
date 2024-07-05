package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code,omitempty"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type UtilsResponse struct{}

func (u *UtilsResponse) Success(message string, data interface{}) Response {
	return Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func (u *UtilsResponse) Error(code int, message string) Response {
	return Response{
		Code:    code,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

func (u *UtilsResponse) JsonResponse(ctx *gin.Context, response Response) *gin.Context {
	ctx.JSON(response.Code, response)
	return ctx
}
