package health

import "github.com/gin-gonic/gin"

func (h *Handler) Router(r *gin.Engine) {
	r.GET("/health/liveness", h.Liveness)
	r.GET("/health/readiness", h.Readiness)
}
