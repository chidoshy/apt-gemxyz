package main

import (
	"git.xantus.network/apt-gemxyz/config"
	"git.xantus.network/apt-gemxyz/impl/api"
	"git.xantus.network/apt-gemxyz/impl/market"
	"git.xantus.network/apt-gemxyz/inf/adapter"
	"git.xantus.network/apt-gemxyz/inf/repository"
	"git.xantus.network/apt-gemxyz/pkg/const_data"
	"git.xantus.network/apt-gemxyz/pkg/core/middleware"
	"git.xantus.network/apt-gemxyz/pkg/database/migrate"
	"git.xantus.network/apt-gemxyz/pkg/db"
	"git.xantus.network/apt-gemxyz/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func run(args []string) error {
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	logger := cfg.Log.Build()

	conn, err := gorm.Open(mysql.Open(cfg.MySQL.DSN()), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Fatal("Error connect to database")
	}
	sqlDB, err := conn.DB()

	if err != nil {
		logger.WithError(err).Fatal("Error get db")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if cfg.MySQL.Debug {
		conn = conn.Debug()
	}

	cmd := &cobra.Command{
		Use: "server",
	}

	container := db.Container{Conn: conn}

	// init repository
	smartContractRepo := repository.NewSmartContractRepositoryImpl(container)
	collectionRepo := repository.NewCollectionRepositoryImpl(container)
	collectionStatRepo := repository.NewCollectionStatRepositoryImpl(container)
	nftRepo := repository.NewNftRepositoryImpl(container)
	transactionRepo := repository.NewTransactionRepositoryImpl(container)
	walletRepo := repository.NewWalletRepositoryImpl(container)
	// transform data adapter
	blueMoveTransformDataAdapter := adapter.NewBlueMoveEventTransformService()
	marketplaceService := market.NewMarketplaceEventImpl(
		smartContractRepo,
		collectionRepo,
		collectionStatRepo,
		nftRepo,
		transactionRepo,
		walletRepo)
	marketplaceService.AddTransformDataAdapter(const_data.BLUEMOVE_MARKETPLACE, blueMoveTransformDataAdapter)

	cmd.AddCommand(&cobra.Command{
		Use: "marketplace-sync",
		Run: func(cmd *cobra.Command, args []string) {
			marketplaceService.Init()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			r := gin.Default()
			r.Use(middleware.CORSMiddleware())
			api.NewAPI(r, container, cfg)
			if err := r.Run("0.0.0.0:8888"); err != nil {
				log.Fatal(err)
			}
		},
	}, migrate.Command(cfg.MigrationsFolder, cfg.MySQL.String()))

	return cmd.Execute()

}
