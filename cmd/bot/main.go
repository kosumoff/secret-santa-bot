package main

import (
	"context"
	"log"

	"github.com/kosumoff/secret-santa-bot/config"
	"github.com/kosumoff/secret-santa-bot/internal/adapter/sqlite"
	"github.com/kosumoff/secret-santa-bot/internal/adapter/telegram"

	"github.com/go-telegram/bot"
)

func main() {
	cfg := config.InitConfig()

	db, err := sqlite.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	pr := sqlite.NewParticipantRepo(db)
	ar := sqlite.NewAssignmentRepo(db)

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	handler := telegram.NewHandler(pr, ar)
	handler.Register(b)

	log.Println("bot started...")
	b.Start(context.Background())
}
