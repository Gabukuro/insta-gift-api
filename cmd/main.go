package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/config"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/log"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/router"
	"github.com/rs/zerolog"
)

func main() {
	// ctx := context.Background()

	time.Local = time.UTC
	logger := log.New(zerolog.InfoLevel)
	config := config.New(logger)

	routerInstance := router.NewRouter(&router.Options{
		AppName: "insta-gift-api",
		Logger:  logger,
	})

	databaseURL := config.GetSecretString("aws.database.insta-gift-api")
	logger.Info().Msgf("databaseURL: %s", databaseURL)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		routerInstance.Stop()
	}()

	logger.Info().Msg("Starting server")
	routerInstance.Start(":3000")
}
