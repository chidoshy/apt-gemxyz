package repository

import (
	"context"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/db"
)

type (
	WalletRepositoryImpl struct {
		conn db.Container
	}
)

func NewWalletRepositoryImpl(conn db.Container) *WalletRepositoryImpl {
	return &WalletRepositoryImpl{conn: conn}
}

var WalletT = (&model.Wallet{}).T()
var WalletC = WalletT.Columns

func (r *WalletRepositoryImpl) GetOrCreate(
	ctx context.Context,
	address string) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	err := r.conn.Conn.
		Where(WalletC.Address+" = ?", address).
		Take(&wallet).Error
	if err != nil && err != db.ErrNotFound {
		return nil, err
	}
	if err == db.ErrNotFound {
		wallet = &model.Wallet{
			Address: address,
		}
		err = r.conn.Conn.Create(&wallet).Error
		if err != nil {
			return nil, err
		}
	}
	return wallet, nil
}
