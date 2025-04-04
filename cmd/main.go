package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"os"
	"os/signal"
	appconfig "telegram-informer/internal/config"
	"telegram-informer/internal/db/sqllite"
	"telegram-informer/internal/server"
)

func main() {
	config, err := appconfig.LoadConfig("./configs/config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := sqllite.New(config.DatabaseConfig.Connection)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(config.Telegram.Token)
	if err != nil {
		panic(err)
	}
	server.RegisterHandlers(b, db)

	go func() {
		fmt.Println("Bot is now running. Press CTRL-C to exit.")
		b.Start(ctx)
	}()

	<-ctx.Done()

	fmt.Println("Shutting down gracefully...")
}
