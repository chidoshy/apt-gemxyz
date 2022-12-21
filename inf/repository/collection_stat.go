package repository

import (
	"context"
	"fmt"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/db"
	"time"
)

type (
	CollectionStatRepositoryImpl struct {
		conn db.Container
	}
)

func NewCollectionStatRepositoryImpl(conn db.Container) *CollectionStatRepositoryImpl {
	return &CollectionStatRepositoryImpl{conn: conn}
}

var CollectionStatT = (&model.CollectionStat{}).T()
var CollectionStatC = CollectionStatT.Columns

func (r *CollectionStatRepositoryImpl) GetOrCreateStatByDate(
	ctx context.Context,
	collectionId int,
	txDate *time.Time) (*model.CollectionStat, error) {
	collectionStat := &model.CollectionStat{}
	err := r.conn.Conn.
		Where(CollectionStatC.CollectionId+" = ?", collectionId).
		Where(fmt.Sprintf(" DATE(%s) = DATE('%s') ", CollectionStatC.CreatedAt, txDate.Format("2006-01-02"))).
		Take(&collectionStat).
		Error
	if err != nil && err != db.ErrNotFound {
		return nil, err
	}
	if err == db.ErrNotFound {
		collectionStat.CollectionId = collectionId
		collectionStat.CreatedAt = txDate
		collectionStat.UpdatedAt = txDate
		err = r.conn.Conn.Create(&collectionStat).Error
		if err != nil {
			return nil, err
		}
	}
	return collectionStat, nil
}

func (r *CollectionStatRepositoryImpl) Save(
	ctx context.Context,
	collectionStat *model.CollectionStat) error {
	err := r.conn.Conn.Save(&collectionStat).Error
	if err != nil {
		return err
	}
	return nil
}
