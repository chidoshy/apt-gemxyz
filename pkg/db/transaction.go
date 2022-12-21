package db

import (
	"context"
	"gorm.io/gorm"
)

type (
	TransactionImpl struct {
		Transaction *gorm.DB
	}
)

func NewTransactionImpl(db Container) *TransactionImpl {
	return &TransactionImpl{
		Transaction: db.Conn.Begin(),
	}
}

func (r *TransactionImpl) TxCommit(ctx context.Context) {
	r.Transaction.Commit()
}

func (r *TransactionImpl) TxRollBack(ctx context.Context) {
	r.Transaction.Rollback()
}
