package repository

import (
	"context"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/db"
)

type (
	MarketplaceRepositoryImpl struct {
		conn db.Container
	}
)

func NewMarketplaceRepositoryImpl(conn db.Container) *MarketplaceRepositoryImpl {
	return &MarketplaceRepositoryImpl{conn: conn}
}

var MarketplaceT = (&model.Marketplace{}).T()
var MarketplaceC = MarketplaceT.Columns

func (r *MarketplaceRepositoryImpl) GetByMarketCollection(
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
	return nft, nil
}
