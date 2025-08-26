package bot

import (
	"context"
	"log"
	"telegram-informer/config"
	"telegram-informer/internal/bot"
	"telegram-informer/internal/infra/cache/redis"
	"telegram-informer/internal/infra/db/postgres"
	"telegram-informer/internal/repo"

	tgBot "github.com/go-telegram/bot"
)

func Run(ctx context.Context) error {
	loadConfig, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	redisCache, err := redis.New(ctx, loadConfig.Redis.Address, loadConfig.Redis.Password, loadConfig.Redis.Db)
	if err != nil {
		log.Fatalf("Ошибка инициализации кеша: %v", err)
	}

	db, err := postgres.New(loadConfig.Database.Connection)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := postgres.RunMigrations(db); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	eventRepository := repo.NewEventRepository(db)

	b, err := tgBot.New(loadConfig.Telegram.Token)
	if err != nil {
		log.Fatalf("Ошибка создания Telegram-бота: %v", err)
	}

	bot.RegisterHandlers(b, eventRepository, redisCache)

	b.Start(ctx)

	return nil
}
