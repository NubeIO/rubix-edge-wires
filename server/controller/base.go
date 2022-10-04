package controller

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"github.com/NubeIO/rubix-edge-wires/flow"
	"github.com/gin-gonic/gin"
	"net/http"
)

var fileUtils = fileutils.New()

type Controller struct {
	Flow *flow.Flow
}

type Response struct {
	StatusCode   int         `json:"status_code"`
	ErrorMessage string      `json:"error_message"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

func reposeHandler(body interface{}, err error, c *gin.Context, statusCode ...int) {
	var code int
	if err != nil {
		if len(statusCode) > 0 {
			code = statusCode[0]
		} else {
			code = http.StatusNotFound
		}
		msg := Message{
			Message: fmt.Sprintf("flow: %s", err.Error()),
		}
		c.JSON(code, msg)
	} else {
		if len(statusCode) > 0 {
			code = statusCode[0]
		} else {
			code = http.StatusOK
		}
		c.JSON(code, body)
	}
}

type Message struct {
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
