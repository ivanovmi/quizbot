package main

import (
	// Import develop version, because stable version with poll feature not fucking released yet
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ivanovmi/quizbot/pkg/quiz"
	"github.com/jasonlvhit/gocron"
)

// TOKEN is telegram bot token
const TOKEN = os.Getenv("BOT_TOKEN")

func main() {
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	bot.Debug = true
	/*  u := tgbotapi.NewUpdate(0) */
	// u.Timeout = 60
	// updates := bot.GetUpdatesChan(u)
	// for update := range updates {
	// if update.Message == nil {
	// continue
	// }
	// fmt.Println(update.Message.Chat.ID)
	/* } */
	gocron.Every(1).Day().At("11:00").Do(SendRuMsg, bot)
	gocron.Every(1).Day().At("18:00").Do(SendEnMsg, bot)
	<-gocron.Start()
}
