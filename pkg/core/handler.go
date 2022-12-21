package core

import (
	"git.xantus.network/apt-gemxyz/lerror"
	"git.xantus.network/apt-gemxyz/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *Handler) Error(c *gin.Context, err error, statusCode ...int) {
	if !lerror.IsLError(err) {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ler := lerror.Unwrap(err)
	res := Response{
		Code:    ler.Code,
		Message: ler.Message,
	}
	c.JSON(ler.Status, res)
}

func (h *Handler) Response(c *gin.Context, data interface{}, status ...int) {
	statusCode := http.StatusOK
	if len(status) > 0 {
		statusCode = status[0]
	}
	c.JSON(statusCode, Response{
		Code:    2000,
		Message: "Success",
		Data:    data,
	})
}
