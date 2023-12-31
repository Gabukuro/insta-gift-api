package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Gabukuro/insta-gift-api/internal/domain/prediction"
	"github.com/Gabukuro/insta-gift-api/internal/domain/product"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/config"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/database"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/log"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/router"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func main() {
	ctx := context.Background()

	time.Local = time.UTC
	logger := log.New(zerolog.InfoLevel)
	config := config.New(logger)

	snsClient := config.AwsProvider.SNS()

	routerInstance := router.NewRouter(&router.Options{
		AppName: "insta-gift-api",
		Logger:  logger,
	})

	databaseURL := config.GetDataBaseSecret("insta-gift-api")
	databaseInstance := database.New(databaseURL, 1, logger).Connect()
	databaseBun := bun.NewDB(databaseInstance.DB, pgdialect.New())

	productRepo := product.NewRepository(databaseBun, logger)
	productService := product.NewService(product.ServiceParams{
		Ctx:        ctx,
		Repository: productRepo,
		Logger:     logger,
	})

	predictionRepo := prediction.NewRepository(databaseBun, logger)
	predictionService := prediction.NewService(prediction.ServiceParams{
		Ctx:                     ctx,
		Repository:              predictionRepo,
		Logger:                  logger,
		SNSClient:               snsClient,
		PredictionEventTopicArn: config.GetSecretString("aws.sns.insta-gift-api.analysisProfileEvents"),
	})
	prediction.NewHTTPHandler(routerInstance.App(), predictionService, productService, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		routerInstance.Stop()
	}()

	logger.Info().Msg("Starting server")
	routerInstance.Start(":8000")
}
