package main

import (
	// Import develop version, because stable version with poll feature not fucking released yet
	//	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	quiz "github.com/ivanovmi/quizbot/pkg/quiz"
	"github.com/ivanovmi/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"os"
)

// TOKEN is telegram bot token
var TOKEN = os.Getenv("BOT_TOKEN")

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
	gocron.Every(1).Day().At("10:00").Do(quiz.SendMsg, bot, "ru")
	gocron.Every(1).Day().At("14:00").Do(quiz.SendMsg, bot, "ru-tf")
	gocron.Every(1).Day().At("18:00").Do(quiz.SendMsg, bot, "en")
	<-gocron.Start()
}
