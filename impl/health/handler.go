package health

import (
	"context"
	"git.xantus.network/apt-gemxyz/pkg/core"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Handler struct {
	srv Service
	core.Handler
}

func NewHandler(srv Service) *Handler {
	return &Handler{srv: srv}
}

type Service interface {
	isLive(ctx context.Context) (bool, error)
	isReady(ctx context.Context) (bool, error)
}

func (h *Handler) Liveness(c *gin.Context) {
	isLive, err := h.srv.isLive(c)
	if err != nil {
		h.Error(c, err)
		return
	}
	if !isLive {
		h.Error(c, errors.New("Service is unavailable!"))
		return
	}
	h.Response(c, err)
}
func (h *Handler) Readiness(c *gin.Context) {
	isReady, err := h.srv.isReady(c)
	if err != nil {
		h.Error(c, err)
		return
	}
	if !isReady {
		h.Error(c, errors.New("Service is not ready!"))
		return
	}
	h.Response(c, err)
}
