package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.xantus.network/apt-gemxyz/lerror"
	"io"
	"net/http"
)

type (
	APTTransaction struct {
		Version   string  `json:"version"`
		Hash      string  `json:"hash"`
		GasUsed   string  `json:"gas_used"`
		Sender    string  `json:"sender"`
		Timestamp string  `json:"timestamp"`
		Events    []Event `json:"events"`
	}
	Event struct {
		GUid           GUid   `json:"guid"`
		SequenceNumber string `json:"sequence_number"`
		Type           string `json:"type"`
		Data           Data   `json:"data"`
	}
	GUid struct {
		CreationNumber string `json:"creation_number"`
		AccountAddress string `json:"account_address"`
	}

	Data struct {
		Amount              string              `json:"amount"`
		Id                  NFTInformation      `json:"id"`
		TokenID             NFTInformation      `json:"token_id"`
		OfferCollectionItem OfferCollectionItem `json:"offer_collection_item"`
		RoyaltyDenominator  string              `json:"royalty_denominator"`
		RoyaltyNumerator    string              `json:"royalty_numerator"`
		RoyaltyPayee        string              `json:"royalty_payee"`
		SellerAddress       string              `json:"seller_address"`
		BuyerAddress        string              `json:"buyer_address"`
		Price               string              `json:"price"`
		Seller              string              `json:"seller"`
	}

	NFTInformation struct {
		PropertyVersion string      `json:"property_version"`
		TokenDataId     TokenDataId `json:"token_data_id"`
	}

	OfferCollectionItem struct {
		AmountPerItem string `json:"amount_per_item"`
		Offerer       string `json:"offerer"`
	}

	TokenDataId struct {
		Collection string `json:"collection"`
		Creator    string `json:"creator"`
		Name       string `json:"name"`
	}

	VersionFilter struct {
		OperationName string    `json:"operationName"`
		Variables     Variables `json:"variables"`
		Query         string    `json:"query"`
	}
	Variables struct {
		Address string `json:"address"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
	}
	Transaction struct {
		TransactionVersion int64 `json:"transaction_version"`
	}
	Resource struct {
		MoveResources []Transaction `json:"move_resources"`
	}
	FilterResponse struct {
		Data Resource `json:"data"`
	}
)

func GetTransactionFilter(
	ctx context.Context,
	nodeUrl string,
	filterAddress string,
	filterNodeUrl string,
	offsetStack int) ([]APTTransaction, int, error) {
	req := VersionFilter{
		OperationName: "AccountTransactionsData",
		Variables: Variables{
			Address: filterAddress,
			Limit:   25,
			Offset:  offsetStack,
		},
		Query: "query AccountTransactionsData($address: String, $limit: Int, $offset: Int) { move_resources(   where: {address: {_eq: $address}}  order_by: {transaction_version: asc}   distinct_on: transaction_version  limit: $limit offset: $offset ) {  transaction_version   __typename  }}",
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, 0, err
	}

	url := filterNodeUrl
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(reqJson),
	)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != 200 {
		return nil, 0, lerror.InternalServer.ToError("Get data error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	respData := FilterResponse{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, 0, err
	}
	var transaction []APTTransaction
	if len(respData.Data.MoveResources) > 0 {
		for _, item := range respData.Data.MoveResources {
			tx, err := GetTransactionByVersion(ctx, nodeUrl, item.TransactionVersion)
			if err != nil {
				return nil, 0, nil
			}
			transaction = append(transaction, *tx)
		}
		return transaction, offsetStack + len(respData.Data.MoveResources), nil
	}
	return nil, offsetStack + len(respData.Data.MoveResources), nil
}

func GetTransactionByVersion(
	ctx context.Context,
	nodeUrl string,
	version int64) (*APTTransaction, error) {
	url := fmt.Sprintf("%s/v1/transactions/by_version/%d", nodeUrl, version)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, lerror.InternalServer.ToError(fmt.Sprintf("Get transaction error - url:%s", url))
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var transactions APTTransaction
	err = json.Unmarshal(responseData, &transactions)
	if err != nil {
		return nil, err
	}
	return &transactions, nil
}
