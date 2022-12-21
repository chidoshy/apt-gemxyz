package model

import (
	"gorm.io/datatypes"
	"time"
)

type Transaction struct {
	ID               int            `json:"id"`
	TxHash           string         `json:"txHash"`
	TxCreationNumber string         `json:"txCreationNumber"`
	TxSequenceNumber string         `json:"txSequenceNumber"`
	MarketplaceId    int            `json:"marketplaceId"`
	CollectionId     int            `json:"collectionId"`
	WalletAddress    string         `json:"walletAddress"`
	Event            string         `json:"event"`
	NftId            int            `json:"nftId"`
	FromAddress      string         `json:"fromAddress"`
	ToAddress        string         `json:"toAddress"`
	Price            float64        `json:"price"`
	Data             datatypes.JSON `json:"data"`
	CreatedAt        *time.Time     `json:"createdAt" gorm:"autoUpdateTime:false"`
	UpdatedAt        *time.Time     `json:"updatedAt" gorm:"autoUpdateTime:false"`
	DeletedAt        *time.Time     `json:"deletedAt"`
}

type tblTransactionColumns struct {
	ID               string
	TxHash           string
	TxCreationNumber string
	TxSequenceNumber string
	MarketplaceId    string
	CollectionId     string
	WalletAddress    string
	Event            string
	NftId            string
	FromAddress      string
	ToAddress        string
	Price            string
	Data             string
	CreatedAt        string
	UpdatedAt        string
	DeletedAt        string
}

type tblTransaction struct {
	Name    string
	Columns tblTransactionColumns
}

var tblTransactionDefine = tblTransaction{
	Name: "transactions",
	Columns: tblTransactionColumns{
		ID:               "id",
		TxHash:           "tx_hash",
		TxCreationNumber: "tx_creation_number",
		TxSequenceNumber: "tx_sequence_number",
		MarketplaceId:    "marketplace_id",
		CollectionId:     "collection_id",
		WalletAddress:    "wallet_address",
		Event:            "event",
		NftId:            "nft_id",
		FromAddress:      "from_address",
		ToAddress:        "to_address",
		Price:            "price",
		Data:             "data",
		CreatedAt:        "created_at",
		UpdatedAt:        "updated_at",
		DeletedAt:        "deleted_at",
	},
}

func (t tblTransaction) GetColumns() []string {
	return []string{
		"id",
		"tx_hash",
		"tx_creation_number",
		"tx_sequence_number",
		"marketplace_id",
		"collection_id",
		"wallet_address",
		"event",
		"nft_id",
		"from_address",
		"to_address",
		"price",
		"data",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (r Transaction) T() tblTransaction {
	return tblTransactionDefine
}
