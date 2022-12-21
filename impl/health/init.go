package health

import (
	"git.xantus.network/apt-gemxyz/pkg/db"
)

func Initial(db db.Container) *Handler {
	srv := NewServiceImpl(db)
	return NewHandler(srv)
}
