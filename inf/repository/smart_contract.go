package repository

import (
	"context"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/db"
)

type (
	SmartContractRepositoryImpl struct {
		conn db.Container
	}
)

func NewSmartContractRepositoryImpl(conn db.Container) *SmartContractRepositoryImpl {
	return &SmartContractRepositoryImpl{conn: conn}
}

var SmartContractT = (&model.SmartContract{}).T()
var SmartContractC = SmartContractT.Columns

func (r *SmartContractRepositoryImpl) GetAll(
	ctx context.Context) ([]model.SmartContract, error) {
	smContracts := []model.SmartContract{}
	err := r.conn.Conn.Table(SmartContractT.Name).Find(&smContracts).Error
	if err != nil {
		return nil, err
	}
	return smContracts, err
}

func (r *SmartContractRepositoryImpl) Get(
	ctx context.Context,
	id int) (*model.SmartContract, error) {
	sm := &model.SmartContract{}
	err := r.conn.Conn.
		Preload("Marketplace").
		Where(SmartContractC.ID+" = ?", id).Take(&sm).Error
	if err != nil {
		return nil, err
	}
	return sm, err
}

func (r *SmartContractRepositoryImpl) Save(
	ctx context.Context,
	sm *model.SmartContract) error {
	err := r.conn.Conn.Save(&sm).Error
	if err != nil {
		return err
	}
	return nil
}
