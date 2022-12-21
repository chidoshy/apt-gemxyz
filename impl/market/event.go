package market

import (
	"context"
	"encoding/json"
	"fmt"
	"git.xantus.network/apt-gemxyz/inf/adapter"
	"git.xantus.network/apt-gemxyz/model"
	"git.xantus.network/apt-gemxyz/pkg/aptconvert"
	"git.xantus.network/apt-gemxyz/pkg/const_data"
	"git.xantus.network/apt-gemxyz/pkg/log"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

type (
	MarketplaceEventImpl struct {
		smartContractRepo    SmartContractRepo
		collectionRepo       CollectionRepo
		collectionStatRepo   CollectionStatRepo
		nftRepo              NftRepo
		transactionRepo      TransactionRepo
		walletRepo           WalletRepo
		transformDataAdapter map[string]TransformDataAdapter
	}
	SmartContractRepo interface {
		GetAll(
			ctx context.Context) ([]model.SmartContract, error)
		Get(
			ctx context.Context,
			id int) (*model.SmartContract, error)
		Save(
			ctx context.Context,
			sm *model.SmartContract) error
	}
	CollectionRepo interface {
		GetOrCreate(
			ctx context.Context,
			marketplaceId int,
			name string,
			creator string) (*model.Collection, error)
	}
	CollectionStatRepo interface {
		GetOrCreateStatByDate(
			ctx context.Context,
			collectionId int,
			txDate *time.Time) (*model.CollectionStat, error)
		Save(
			ctx context.Context,
			collectionStat *model.CollectionStat) error
	}
	NftRepo interface {
		GetByCollectionName(
			ctx context.Context,
			marketplaceId int,
			collectionId int,
			name string) (*model.Nft, error)
		Create(
			ctx context.Context,
			nft *model.Nft) error
		Save(
			ctx context.Context,
			nft *model.Nft) error
	}
	TransactionRepo interface {
		GetByTxHash(
			ctx context.Context,
			txHash string,
			creationNumber string,
			sequenceNumber string) (*model.Transaction, error)
		Create(
			ctx context.Context,
			tx *model.Transaction) (*model.Transaction, error)
		Save(
			ctx context.Context,
			tx *model.Transaction) error
	}
	WalletRepo interface {
		GetOrCreate(
			ctx context.Context,
			address string) (*model.Wallet, error)
	}
	TransformDataAdapter interface {
		TransformMarketplaceEvent(
			sm *model.SmartContract,
			transaction adapter.APTTransaction,
			eventName string,
			event adapter.Event) adapter.MarketplaceEvent
	}
)

func NewMarketplaceEventImpl(
	smartContractRepo SmartContractRepo,
	collectionRepo CollectionRepo,
	collectionStatRepo CollectionStatRepo,
	nftRepo NftRepo,
	transactionRepo TransactionRepo,
	walletRepo WalletRepo) *MarketplaceEventImpl {
	return &MarketplaceEventImpl{
		smartContractRepo:    smartContractRepo,
		collectionRepo:       collectionRepo,
		collectionStatRepo:   collectionStatRepo,
		nftRepo:              nftRepo,
		transactionRepo:      transactionRepo,
		walletRepo:           walletRepo,
		transformDataAdapter: make(map[string]TransformDataAdapter),
	}
}

func (s MarketplaceEventImpl) AddTransformDataAdapter(marketName string, adapter TransformDataAdapter) {
	s.transformDataAdapter[marketName] = adapter
}

func (s *MarketplaceEventImpl) Init() {
	// init context
	ctx := context.Background()
	// get all smart contract
	smartContracts, err := s.smartContractRepo.GetAll(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("Get all smart contract error - %v", err))
	}
	// init go routine to sync all contract
	wg := sync.WaitGroup{}
	for _, sm := range smartContracts {
		go s.Run(ctx, sm.ID)
		wg.Add(1)
	}
	wg.Wait()
}

func (s *MarketplaceEventImpl) Run(
	ctx context.Context,
	smartContractId int) {
	for {
		// get smart contract
		smartContract, err := s.smartContractRepo.Get(ctx, smartContractId)
		if err != nil {
			log.Fatal(fmt.Sprintf("Get smart contract error - %v", err))
		}
		// get transaction
		transactions, offsetStack, err := adapter.GetTransactionFilter(
			ctx,
			smartContract.NodeUrl,
			smartContract.ResourceAddress,
			smartContract.ResourceNodeUrl,
			smartContract.ResourceOffsetStack)
		if err != nil {
			log.Debug(fmt.Sprintf("Get transaction error - %v", err))
			log.Debug("Something wrong - sleep 30s")
			time.Sleep(30)
			continue
		}
		// create map event
		mapEventHandle := make(map[string]string)
		err = json.Unmarshal(smartContract.Event, &mapEventHandle)
		if err != nil {
			log.Fatal(fmt.Sprintf("Create map event error - %v", err))
		}
		// handle transaction
		for _, tx := range transactions {
			// loop event
			for _, event := range tx.Events {
				if mapEventHandle[event.Type] != "" {
					// transform event
					marketplaceEventData := s.transformDataAdapter[smartContract.Marketplace.Name].
						TransformMarketplaceEvent(smartContract, tx, mapEventHandle[event.Type], event)
					// check duplicate event
					txHistory, err := s.transactionRepo.GetByTxHash(
						ctx, marketplaceEventData.TxHash, marketplaceEventData.TxCreationNumber, marketplaceEventData.TxSequenceNumber)
					if err != nil {
						log.Fatal(fmt.Sprintf("Get tx history error - %v", err))
					}
					if txHistory != nil {
						log.Debug("Duplicate event")
						continue
					}
					// get collection
					collection, err := s.collectionRepo.GetOrCreate(
						ctx,
						marketplaceEventData.MarketplaceId,
						marketplaceEventData.CollectionName,
						marketplaceEventData.Creator)
					if err != nil {
						log.Fatal(fmt.Sprintf("Get collection error - %v", err))
					}
					// get transaction date
					txDate, err := s.TransactionDate(marketplaceEventData.Timestamp)
					if err != nil {
						log.Fatal(fmt.Sprintf("Get transaction date error - %v", err))
					}
					// get collection stat
					collectionStat, err := s.collectionStatRepo.GetOrCreateStatByDate(ctx, collection.ID, txDate)
					if err != nil {
						log.Fatal(fmt.Sprintf("Get collection stat error - %v", err))
					}
					// create transaction history
					txHistory, err = s.transactionRepo.Create(ctx, &model.Transaction{
						TxHash:           marketplaceEventData.TxHash,
						TxCreationNumber: marketplaceEventData.TxCreationNumber,
						TxSequenceNumber: marketplaceEventData.TxSequenceNumber,
						MarketplaceId:    marketplaceEventData.MarketplaceId,
						CollectionId:     collection.ID,
						Event:            marketplaceEventData.EventName,
						CreatedAt:        txDate,
						UpdatedAt:        txDate,
					})
					if err != nil {
						log.Fatal(fmt.Sprintf("Create transaction history error - %v", err))
					}
					switch mapEventHandle[event.Type] {
					case const_data.LISTING_EVENT:
						err = s.HandleListingEvent(ctx, marketplaceEventData, collection, collectionStat, txHistory, txDate)
						if err != nil {
							log.Fatal(fmt.Sprintf("HandleListingEvent error - %v", err))
						}
					case const_data.BUY_EVENT:
						err = s.HandleBuyEvent(ctx, marketplaceEventData, collection, collectionStat, txHistory, txDate)
						if err != nil {
							log.Fatal(fmt.Sprintf("HandleListingEvent error - %v", err))
						}
					case const_data.DELIST_EVENT:
						err = s.HandleDeListEvent(ctx, marketplaceEventData, collection, collectionStat, txHistory, txDate)
						if err != nil {
							log.Fatal(fmt.Sprintf("HandleListingEvent error - %v", err))
						}
					case const_data.CHANGE_PRICE_EVENT:
						err = s.HandleChangePriceEvent(ctx, marketplaceEventData, collection, collectionStat, txHistory, txDate)
						if err != nil {
							log.Fatal(fmt.Sprintf("HandleListingEvent error - %v", err))
						}
					case const_data.ACCEPT_OFFER_EVENT:
						err = s.HandleAcceptOfferEvent(ctx, marketplaceEventData, collection, collectionStat, txHistory, txDate)
						if err != nil {
							log.Fatal(fmt.Sprintf("HandleListingEvent error - %v", err))
						}
					default:
						continue
					}
				}
			}
		}
		// save offset stack
		smartContract.ResourceOffsetStack = offsetStack
		err = s.smartContractRepo.Save(ctx, smartContract)
		if err != nil {
			log.Fatal(fmt.Sprintf("Save smart contract error - %v", err))
		}
	}
}

// func support

func (s *MarketplaceEventImpl) HandleListingEvent(
	ctx context.Context,
	marketplaceEventData adapter.MarketplaceEvent,
	collection *model.Collection,
	collectionStat *model.CollectionStat,
	transactionHistory *model.Transaction,
	txDate *time.Time) error {
	// convert price
	price, err := aptconvert.FromWeiStrToAptFloat(marketplaceEventData.Price)
	// get or create nft
	nft, err := s.nftRepo.GetByCollectionName(ctx,
		marketplaceEventData.MarketplaceId,
		collection.ID,
		marketplaceEventData.NftName)
	if err != nil {
		return err
	}
	wallet, err := s.walletRepo.GetOrCreate(ctx, marketplaceEventData.SellerAddress)
	if err != nil {
		return err
	}
	if nft == nil {
		// create nft
		err = s.nftRepo.Create(ctx, &model.Nft{
			MarketplaceId:      marketplaceEventData.MarketplaceId,
			CollectionId:       collection.ID,
			OwnerWalletAddress: wallet.Address,
			Name:               marketplaceEventData.NftName,
			ListingPrice:       price,
			Status:             const_data.LISTING_STATUS,
		})
		if err != nil {
			return err
		}
	} else {
		// update owner, price
		nft.OwnerWalletAddress = marketplaceEventData.SellerAddress
		nft.ListingPrice = price
		nft.Status = const_data.LISTING_STATUS
		err = s.nftRepo.Save(ctx, nft)
		if err != nil {
			return err
		}
	}
	// check floor price
	if collectionStat.FloorPrice == 0 || collectionStat.FloorPrice > price {
		collectionStat.FloorPrice = price
	}
	// increase total buy event
	collectionStat.TotalListing = collectionStat.TotalListing + 1
	err = s.collectionStatRepo.Save(ctx, collectionStat)
	if err != nil {
		return err
	}

	// save user address
	transactionHistory.CollectionId = collection.ID
	transactionHistory.WalletAddress = marketplaceEventData.SellerAddress
	transactionHistory.FromAddress = marketplaceEventData.SellerAddress
	transactionHistory.NftId = nft.ID
	transactionHistory.Price = price
	err = s.transactionRepo.Save(ctx, transactionHistory)
	if err != nil {
		return err
	}
	return nil
}

func (s *MarketplaceEventImpl) HandleDeListEvent(
	ctx context.Context,
	marketplaceEventData adapter.MarketplaceEvent,
	collection *model.Collection,
	collectionStat *model.CollectionStat,
	transactionHistory *model.Transaction,
	txDate *time.Time) error {
	// get or create nft
	nft, err := s.nftRepo.GetByCollectionName(ctx,
		marketplaceEventData.MarketplaceId,
		collection.ID,
		marketplaceEventData.NftName)
	if err != nil {
		return err
	}
	if nft != nil {
		// update status nft
		nft.Status = const_data.AVAILABLE_STATUS
		err = s.nftRepo.Save(ctx, nft)
		if err != nil {
			return err
		}
		// save user address
		transactionHistory.CollectionId = collection.ID
		transactionHistory.WalletAddress = marketplaceEventData.SellerAddress
		transactionHistory.FromAddress = marketplaceEventData.SellerAddress
		transactionHistory.NftId = nft.ID
		err = s.transactionRepo.Save(ctx, transactionHistory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MarketplaceEventImpl) HandleBuyEvent(
	ctx context.Context,
	marketplaceEventData adapter.MarketplaceEvent,
	collection *model.Collection,
	collectionStat *model.CollectionStat,
	transactionHistory *model.Transaction,
	txDate *time.Time) error {
	// convert price
	price, err := aptconvert.FromWeiStrToAptFloat(marketplaceEventData.Price)
	// get or create nft
	nft, err := s.nftRepo.GetByCollectionName(ctx,
		marketplaceEventData.MarketplaceId,
		collection.ID,
		marketplaceEventData.NftName)
	if err != nil {
		return err
	}
	wallet, err := s.walletRepo.GetOrCreate(ctx, marketplaceEventData.BuyerAddress)
	if err != nil {
		return err
	}

	if nft != nil {
		sellerAddress := nft.OwnerWalletAddress
		// update owner, price, status
		nft.OwnerWalletAddress = wallet.Address
		nft.ListingPrice = price
		nft.Status = const_data.AVAILABLE_STATUS
		err = s.nftRepo.Save(ctx, nft)
		if err != nil {
			return err
		}
		// increase total buy event
		collectionStat.TotalVolume = collectionStat.TotalVolume + price
		collectionStat.TotalBuy = collectionStat.TotalBuy + 1
		err = s.collectionStatRepo.Save(ctx, collectionStat)
		if err != nil {
			return err
		}

		// save user address
		transactionHistory.CollectionId = collection.ID
		transactionHistory.WalletAddress = marketplaceEventData.BuyerAddress
		transactionHistory.FromAddress = sellerAddress
		transactionHistory.ToAddress = marketplaceEventData.BuyerAddress
		transactionHistory.NftId = nft.ID
		transactionHistory.Price = price
		err = s.transactionRepo.Save(ctx, transactionHistory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MarketplaceEventImpl) HandleChangePriceEvent(
	ctx context.Context,
	marketplaceEventData adapter.MarketplaceEvent,
	collection *model.Collection,
	collectionStat *model.CollectionStat,
	transactionHistory *model.Transaction,
	txDate *time.Time) error {
	// convert price
	price, err := aptconvert.FromWeiStrToAptFloat(marketplaceEventData.Price)
	// get or create nft
	nft, err := s.nftRepo.GetByCollectionName(ctx,
		marketplaceEventData.MarketplaceId,
		collection.ID,
		marketplaceEventData.NftName)
	if err != nil {
		return err
	}
	if nft != nil {
		// update owner, price, status
		nft.ListingPrice = price
		err = s.nftRepo.Save(ctx, nft)
		if err != nil {
			return err
		}

		// save user address
		transactionHistory.CollectionId = collection.ID
		transactionHistory.WalletAddress = marketplaceEventData.SellerAddress
		transactionHistory.FromAddress = marketplaceEventData.SellerAddress
		transactionHistory.NftId = nft.ID
		transactionHistory.Price = price
		err = s.transactionRepo.Save(ctx, transactionHistory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MarketplaceEventImpl) HandleAcceptOfferEvent(
	ctx context.Context,
	marketplaceEventData adapter.MarketplaceEvent,
	collection *model.Collection,
	collectionStat *model.CollectionStat,
	transactionHistory *model.Transaction,
	txDate *time.Time) error {
	// convert price
	price, err := aptconvert.FromWeiStrToAptFloat(marketplaceEventData.OffererAmountPerItem)
	// get or create nft
	nft, err := s.nftRepo.GetByCollectionName(ctx,
		marketplaceEventData.MarketplaceId,
		collection.ID,
		marketplaceEventData.NftName)
	if err != nil {
		return err
	}
	// get wallet
	wallet, err := s.walletRepo.GetOrCreate(ctx, marketplaceEventData.OffererAddress)
	if err != nil {
		return err
	}

	if nft != nil {
		sellerAddress := nft.OwnerWalletAddress
		// update owner
		nft.OwnerWalletAddress = wallet.Address
		nft.Status = const_data.AVAILABLE_STATUS
		err = s.nftRepo.Save(ctx, nft)
		if err != nil {
			return err
		}
		// increase total buy event
		collectionStat.TotalVolume = collectionStat.TotalVolume + price
		collectionStat.TotalBuy = collectionStat.TotalBuy + 1
		err = s.collectionStatRepo.Save(ctx, collectionStat)
		if err != nil {
			return err
		}

		// save user address
		transactionHistory.CollectionId = collection.ID
		transactionHistory.WalletAddress = marketplaceEventData.OffererAddress
		transactionHistory.FromAddress = sellerAddress
		transactionHistory.ToAddress = marketplaceEventData.OffererAddress
		transactionHistory.NftId = nft.ID
		transactionHistory.Price = price
		err = s.transactionRepo.Save(ctx, transactionHistory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MarketplaceEventImpl) TransactionDate(timestamp string) (*time.Time, error) {
	value, err := decimal.NewFromString(timestamp)
	if err != nil {
		return nil, err
	}
	tm := time.UnixMicro(value.IntPart())
	return &tm, nil
}
