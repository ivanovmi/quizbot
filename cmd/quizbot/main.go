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

// QuizSchedule is schedule for questions
var QuizSchedule = map[string]string{
	"10:00": "ru",
	"14:00": "ru-tf",
	"18:00": "en",
}

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
	for time, t := range QuizSchedule {
		gocron.Every(1).Day().At(time).Do(quiz.SendMsg, bot, t)
	}
	gocron.Every(1).Sunday().At("15:00").Do(quiz.SendScheduleMsg, bot)
	gocron.Every(1).Thursday().At("15:00").Do(quiz.SendRatingMsg, bot)
	<-gocron.Start()
}
