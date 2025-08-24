package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"telegram-informer/cmd/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		fmt.Println("🤖 Bot is running. Press CTRL+C to stop.")

		err := bot.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	fmt.Println("🛑 Shutting down gracefully...")
}
