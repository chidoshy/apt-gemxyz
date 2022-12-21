package model

type Collection struct {
	ID            int    `json:"id"`
	MarketplaceId int    `json:"marketplaceId"`
	Name          string `json:"name"`
	Creator       string `json:"creator"`
}

type tblCollectionColumns struct {
	ID            string
	MarketplaceId string
	Name          string
	Creator       string
}

type tblCollection struct {
	Name    string
	Columns tblCollectionColumns
}

var tblCollectionDefine = tblCollection{
	Name: "collections",
	Columns: tblCollectionColumns{
		ID:            "id",
		MarketplaceId: "marketplace_id",
		Name:          "name",
		Creator:       "creator",
	},
}

func (t tblCollection) GetColumns() []string {
	return []string{
		"id",
		"marketplace_id",
		"name",
		"creator",
	}
}

func (r Collection) T() tblCollection {
	return tblCollectionDefine
}
