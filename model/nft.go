package model

import (
	"gorm.io/datatypes"
	"time"
)

type Nft struct {
	ID                 int            `json:"id"`
	MarketplaceId      int            `json:"marketplaceId"`
	CollectionId       int            `json:"collectionId"`
	OwnerWalletAddress string         `json:"ownerWalletAddress"`
	Name               string         `json:"name"`
	Metadata           datatypes.JSON `json:"metadata"`
	ListingPrice       float64        `json:"listingPrice"`
	ListingStartDate   *time.Time     `json:"listingStartDate"`
	ListingEndDate     *time.Time     `json:"listingEndDate"`
	Status             string         `json:"status"`
}

type tblNftColumns struct {
	ID                 string
	MarketplaceId      string
	CollectionId       string
	OwnerWalletAddress string
	Name               string
	Metadata           string
	ListingPrice       string
	ListingStartDate   string
	ListingEndDate     string
	Status             string
}

type tblNft struct {
	Name    string
	Columns tblNftColumns
}

var tblNftDefine = tblNft{
	Name: "nfts",
	Columns: tblNftColumns{
		ID:                 "id",
		MarketplaceId:      "marketplace_id",
		CollectionId:       "collection_id",
		OwnerWalletAddress: "owner_wallet_address",
		Name:               "name",
		Metadata:           "metadata",
		ListingPrice:       "listing_price",
		ListingStartDate:   "listing_start_date",
		ListingEndDate:     "listing_end_date",
		Status:             "status",
	},
}

func (t tblNft) GetColumns() []string {
	return []string{
		"id",
		"marketplace_id",
		"collection_id",
		"owner_wallet_address",
		"name",
		"metadata",
		"listing_price",
		"listing_start_date",
		"listing_end_date",
		"status",
	}
}

func (r Nft) T() tblNft {
	return tblNftDefine
}
