package model

type Marketplace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type tblMarketplaceColumns struct {
	ID   string
	Name string
}

type tblMarketplace struct {
	Name    string
	Columns tblMarketplaceColumns
}

var tblMarketplaceDefine = tblMarketplace{
	Name: "marketplaces",
	Columns: tblMarketplaceColumns{
		ID:   "id",
		Name: "name",
	},
}

func (t tblMarketplace) GetColumns() []string {
	return []string{
		"id",
		"name",
	}
}

func (r Marketplace) T() tblMarketplace {
	return tblMarketplaceDefine
}
