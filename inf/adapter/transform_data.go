package adapter

import (
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/const_data"
)

type (
	MarketplaceEvent struct {
		MarketplaceId        int    `json:"marketplaceId"`
		MarketplaceName      string `json:"marketplaceName"`
		Version              string `json:"version"`
		GasUsed              string `json:"gas_used"`
		Sender               string `json:"sender"`
		Timestamp            string `json:"timestamp"`
		EventName            string `json:"eventName"`
		TxHash               string `json:"txHash"`
		TxCreationNumber     string `json:"txCreationNumber"`
		TxSequenceNumber     string `json:"txSequenceNumber"`
		Price                string `json:"price"`
		SellerAddress        string `json:"sellerAddress"`
		BuyerAddress         string `json:"buyerAddress"`
		OffererAddress       string `json:"offererAddress"`
		OffererAmountPerItem string `json:"offererAmountPerItem"`
		CollectionName       string `json:"collectionName"`
		Creator              string `json:"creator"`
		NftName              string `json:"nftName"`
	}

	BlueMoveEventTransformService struct {
	}

	TopazEventTransformService struct {
	}
)

// Blue Move

func NewBlueMoveEventTransformService() *BlueMoveEventTransformService {
	return &BlueMoveEventTransformService{}
}

func (t *BlueMoveEventTransformService) TransformMarketplaceEvent(
	sm *model.SmartContract,
	transaction APTTransaction,
	eventName string,
	event Event) MarketplaceEvent {
	if eventName == const_data.ACCEPT_OFFER_EVENT {
		event.Data.Id.TokenDataId.Collection = event.Data.TokenID.TokenDataId.Collection
		event.Data.Id.TokenDataId.Creator = event.Data.TokenID.TokenDataId.Creator
		event.Data.Id.TokenDataId.Name = event.Data.TokenID.TokenDataId.Name
	}
	return MarketplaceEvent{
		MarketplaceId:        sm.Marketplace.ID,
		MarketplaceName:      sm.Marketplace.Name,
		Version:              transaction.Version,
		GasUsed:              transaction.GasUsed,
		Sender:               transaction.Sender,
		Timestamp:            transaction.Timestamp,
		EventName:            eventName,
		TxHash:               transaction.Hash,
		TxCreationNumber:     event.GUid.CreationNumber,
		TxSequenceNumber:     event.SequenceNumber,
		Price:                event.Data.Amount,
		SellerAddress:        event.Data.SellerAddress,
		BuyerAddress:         event.Data.BuyerAddress,
		OffererAddress:       event.Data.OfferCollectionItem.Offerer,
		OffererAmountPerItem: event.Data.OfferCollectionItem.AmountPerItem,
		CollectionName:       event.Data.Id.TokenDataId.Collection,
		Creator:              event.Data.Id.TokenDataId.Creator,
		NftName:              event.Data.Id.TokenDataId.Name,
	}
}

// Topaz

func NewTopazEventTransformService() *TopazEventTransformService {
	return &TopazEventTransformService{}
}

func (t *TopazEventTransformService) TransformMarketplaceEvent(sm *model.SmartContract, transaction APTTransaction) {

}
