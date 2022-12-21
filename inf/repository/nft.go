package repository

import (
	"context"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/db"
)

type (
	NftRepositoryImpl struct {
		conn db.Container
	}
)

func NewNftRepositoryImpl(conn db.Container) *NftRepositoryImpl {
	return &NftRepositoryImpl{conn: conn}
}

var NftT = (&model.Nft{}).T()
var NftC = NftT.Columns

func (r *NftRepositoryImpl) GetByCollectionName(
	ctx context.Context,
	marketplaceId int,
	collectionId int,
	name string) (*model.Nft, error) {
	nft := &model.Nft{}
	err := r.conn.Conn.
		Where(NftC.MarketplaceId+" = ?", marketplaceId).
		Where(NftC.CollectionId+" = ?", collectionId).
		Where(NftC.Name+" = ?", name).
		Take(&nft).
		Error
	if err != nil && err != db.ErrNotFound {
		return nil, err
	}
	if err == db.ErrNotFound {
		return nil, nil
	}
	return nft, err
}

func (r *NftRepositoryImpl) Create(
	ctx context.Context,
	nft *model.Nft) error {
	err := r.conn.Conn.Create(&nft).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *NftRepositoryImpl) Save(
	ctx context.Context,
	nft *model.Nft) error {
	err := r.conn.Conn.Save(&nft).Error
	if err != nil {
		return err
	}
	return nil
}
