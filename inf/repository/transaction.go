package repository

import (
	"context"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/db"
)

type (
	TransactionRepositoryImpl struct {
		conn db.Container
	}
)

func NewTransactionRepositoryImpl(conn db.Container) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{conn: conn}
}

var TransactionT = (&model.Transaction{}).T()
var TransactionC = TransactionT.Columns

func (r *TransactionRepositoryImpl) GetByTxHash(
	ctx context.Context,
	txHash string,
	creationNumber string,
	sequenceNumber string) (*model.Transaction, error) {
	fh := model.Transaction{}
	err := r.conn.Conn.
		Where(TransactionC.TxHash+" = ?", txHash).
		Where(TransactionC.TxCreationNumber+" = ?", creationNumber).
		Where(TransactionC.TxSequenceNumber+" = ?", sequenceNumber).
		Take(&fh).Error
	if err != nil && err != db.ErrNotFound {
		return nil, err
	}
	if err == db.ErrNotFound {
		return nil, nil
	}
	return &fh, nil
}

func (r *TransactionRepositoryImpl) Create(
	ctx context.Context,
	tx *model.Transaction) (*model.Transaction, error) {
	err := r.conn.Conn.Create(&tx).Error
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *TransactionRepositoryImpl) Save(
	ctx context.Context,
	tx *model.Transaction) error {
	err := r.conn.Conn.Save(&tx).Error
	if err != nil {
		return err
	}
	return nil
}
