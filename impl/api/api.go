package api

import (
	"git.xantus.network/apt-gemxyz/config"
	"git.xantus.network/apt-gemxyz/impl/health"
	"git.xantus.network/apt-gemxyz/pkg/db"
	"github.com/gin-gonic/gin"
)

func NewAPI(r *gin.Engine, db db.Container, conf *config.Config) {
	//Health
	healthHandler := health.Initial(db)
	healthHandler.Router(r)
}
