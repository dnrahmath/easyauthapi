package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//===================================================================

type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(response.StatusCode, response)
}

//===================================================================

type ResponseByte struct {
	StatusCode  int    `json:"-"`
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	ContentType string `json:"contenttype"`
	Data        []byte `json:"data,omitempty"`
}

func (response *ResponseByte) SendByteWithContentType(c *gin.Context) {
	c.Header("Content-Type", response.ContentType)
	c.Data(response.StatusCode, response.ContentType, response.Data)
}

//===================================================================

func SendResponseData(c *gin.Context, data gin.H) {
	response := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
	response.SendResponse(c)
}

//==================================

func SendErrorResponse(c *gin.Context, status int, message string) {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    message,
	}
	response.SendResponse(c)
}

//===================================================================
