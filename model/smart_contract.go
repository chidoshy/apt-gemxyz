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
	Available           bool           `json:"available"`
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
	Available           string
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
		Available:           "available",
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
		"available",
	}
}

func (r SmartContract) T() tblSmartContract {
	return tblSmartContractDefine
}
