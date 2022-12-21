package model

type Wallet struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

type tblWalletColumns struct {
	ID      string
	Address string
}

type tblWallet struct {
	Name    string
	Columns tblWalletColumns
}

var tblWalletDefine = tblWallet{
	Name: "wallets",
	Columns: tblWalletColumns{
		ID:      "id",
		Address: "address",
	},
}

func (t tblWallet) GetColumns() []string {
	return []string{
		"id",
		"address",
	}
}

func (r Wallet) T() tblWallet {
	return tblWalletDefine
}
