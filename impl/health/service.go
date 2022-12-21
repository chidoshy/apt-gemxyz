package health

import (
	"context"
	"git.xantus.network/apt-gemxyz/pkg/db"
)

type ServiceImpl struct {
	container db.Container
}

func NewServiceImpl(container db.Container) *ServiceImpl {
	return &ServiceImpl{
		container: container,
	}
}

func (s *ServiceImpl) isLive(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *ServiceImpl) isReady(ctx context.Context) (bool, error) {
	dbConnection, err := s.container.Conn.DB()
	if err != nil {
		return false, err
	}
	err = dbConnection.PingContext(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}
