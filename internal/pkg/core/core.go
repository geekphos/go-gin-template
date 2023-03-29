package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"phos.cc/yoo/internal/pkg/errno"
)

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		hcode, code, message := errno.Decode(err)
		c.JSON(hcode, ErrResponse{Code: code, Message: message})
		return
	}

	c.JSON(http.StatusOK, data)
}
