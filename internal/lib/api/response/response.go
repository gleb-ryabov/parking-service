package response

import (
	"github.com/gin-gonic/gin"
)

// Response consists of required fields for all responses.
type Response struct {
	Status string `json:"status"`          // Response status
	Error  string `json:"error,omitempty"` // Error if response is error
	Data   any    `json:"data"`            // Data response
}

const (
	// StatusOK - Status without error.
	StatusOK = "OK"
	// StatusError - Status with error.
	StatusError = "Error"
)

// OK returns a response with OK status.
func OK(data any) Response {
	return Response{
		Status: StatusOK,
		Data:   data,
	}
}

// Error returns a response with the error status and error message.
func Error(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// Respond writes HTTP response in JSON format with status code.
func Respond(c *gin.Context, code int, resp Response) {
	c.JSON(code, resp)
}
