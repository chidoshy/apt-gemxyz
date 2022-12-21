package aptconvert

import (
	"git.xantus.network/apt-gemxyz/pkg/const_data"
	"github.com/shopspring/decimal"
)

func FromWeiStrToAptFloat(amount string) (float64, error) {
	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return 0, err
	}

	value, err := decimal.NewFromString(const_data.APT_BASE)
	if err != nil {
		return 0, err
	}

	aptFloat, _ := amountDec.Div(value).BigFloat().Float64()
	return aptFloat, nil
}

func FromAptFloatToWei(amount float64) (string, error) {
	valueBase, err := decimal.NewFromString(const_data.APT_BASE)
	if err != nil {
		return "", err
	}
	valueFloat := decimal.NewFromFloat(amount)
	return valueFloat.Mul(valueBase).String(), nil
}
