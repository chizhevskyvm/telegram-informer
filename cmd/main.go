package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"log"
	"os"
	"os/signal"
	lCache "telegram-informer/internal/cache/redis"
	appconfig "telegram-informer/internal/config"
	"telegram-informer/internal/db/postgres"
	"telegram-informer/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Load config
	config, err := appconfig.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Init Redis
	redisCache, err := lCache.New(ctx, config.Redis.Address, config.Redis.Password, config.Redis.Db)
	if err != nil {
		log.Fatalf("Ошибка инициализации кеша: %v", err)
	}

	// Init Postgres
	db, err := postgres.New(config.Database.Connection)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	// Init Telegram Bot
	b, err := bot.New(config.Telegram.Token)
	if err != nil {
		log.Fatalf("Ошибка создания Telegram-бота: %v", err)
	}

	// Register handlers
	server.RegisterHandlers(b, db, redisCache)

	// Run bot
	go func() {
		fmt.Println("🤖 Bot is running. Press CTRL+C to stop.")
		b.Start(ctx)
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("🛑 Shutting down gracefully...")
}
