package repository

import (
	"context"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/db"
)

type (
	CollectionRepositoryImpl struct {
		conn db.Container
	}
)

func NewCollectionRepositoryImpl(conn db.Container) *CollectionRepositoryImpl {
	return &CollectionRepositoryImpl{conn: conn}
}

var CollectionT = (&model.Collection{}).T()
var CollectionC = CollectionT.Columns

func (r *CollectionRepositoryImpl) GetOrCreate(
	ctx context.Context,
	marketplaceId int,
	name string,
	creator string) (*model.Collection, error) {
	collection := &model.Collection{}
	err := r.conn.Conn.
		Where(CollectionC.MarketplaceId+" = ?", marketplaceId).
		Where(CollectionC.Name+" = ?", name).
		Take(&collection).Error
	if err != nil && err != db.ErrNotFound {
		return nil, err
	}
	if err == db.ErrNotFound {
		collection = &model.Collection{
			MarketplaceId: marketplaceId,
			Name:          name,
			Creator:       creator,
		}
		err = r.conn.Conn.Create(&collection).Error
		if err != nil {
			return nil, err
		}
	}
	return collection, nil
}
