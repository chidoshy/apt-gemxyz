package model

import "time"

type CollectionStat struct {
	ID           int        `json:"id"`
	CollectionId int        `json:"collectionId"`
	TotalVolume  float64    `json:"totalVolume"`
	TotalBuy     int        `json:"totalBuy"`
	TotalListing int        `json:"totalListing"`
	FloorPrice   float64    `json:"floorPrice"`
	CreatedAt    *time.Time `json:"createdAt" gorm:"autoUpdateTime:false"`
	UpdatedAt    *time.Time `json:"updatedAt" gorm:"autoUpdateTime:false"`
	DeletedAt    *time.Time `json:"deletedAt"`
}

type tblCollectionStatColumns struct {
	ID           string
	CollectionId string
	TotalVolume  string
	TotalBuy     string
	TotalListing string
	FloorPrice   string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}

type tblCollectionStat struct {
	Name    string
	Columns tblCollectionStatColumns
}

var tblCollectionStatDefine = tblCollectionStat{
	Name: "collection_stats",
	Columns: tblCollectionStatColumns{
		ID:           "id",
		CollectionId: "collection_id",
		TotalVolume:  "total_volume",
		TotalBuy:     "total_buy",
		TotalListing: "total_listing",
		FloorPrice:   "floor_price",
		CreatedAt:    "created_at",
		UpdatedAt:    "updated_at",
		DeletedAt:    "deleted_at",
	},
}

func (t tblCollectionStat) GetColumns() []string {
	return []string{
		"id",
		"collection_id",
		"total_volume",
		"total_buy",
		"total_listing",
		"floor_price",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (r CollectionStat) T() tblCollectionStat {
	return tblCollectionStatDefine
}
