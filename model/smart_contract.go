package model

import "gorm.io/datatypes"

type SmartContract struct {
	ID                  int            `json:"id"`
	MarketplaceId       int            `json:"marketplaceId"`
	Marketplace         Marketplace    `json:"marketplace"`
	Address             string         `json:"address"`
	NodeUrl             string         `json:"nodeUrl"`
	Event               datatypes.JSON `json:"event"`
	ResourceAddress     string         `json:"resourceAddress"`
	ResourceNodeUrl     string         `json:"resourceNodeUrl"`
	ResourceOffsetStack int            `json:"resourceOffsetStack"`
}

type tblSmartContractColumns struct {
	ID                  string
	MarketplaceId       string
	Address             string
	NodeUrl             string
	Event               string
	ResourceAddress     string
	ResourceNodeUrl     string
	ResourceOffsetStack string
}

type tblSmartContract struct {
	Name    string
	Columns tblSmartContractColumns
}

var tblSmartContractDefine = tblSmartContract{
	Name: "smart_contracts",
	Columns: tblSmartContractColumns{
		ID:                  "id",
		MarketplaceId:       "marketplace_id",
		Address:             "address",
		NodeUrl:             "node_url",
		Event:               "event",
		ResourceAddress:     "resource_address",
		ResourceNodeUrl:     "resource_node_url",
		ResourceOffsetStack: "resource_offset_stack",
	},
}

func (t tblSmartContract) GetColumns() []string {
	return []string{
		"id",
		"marketplace_id",
		"address",
		"node_url",
		"event",
		"resource_address",
		"resource_node_url",
		"resource_offset_stack",
	}
}

func (r SmartContract) T() tblSmartContract {
	return tblSmartContractDefine
}
