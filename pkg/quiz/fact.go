package quiz

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	//	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ivanovmi/telegram-bot-api"
	"log"
	"net/http"
)

func getFact() (string, error) {
	res, err := http.Get(FactURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fact := ""
	doc.Find("#fact > table > tbody > tr > td").Each(func(i int, s *goquery.Selection) {
		fact = s.Text()
	})
	return fmt.Sprintf("%s %s", bulb, fact), nil
}

func SendFactMsg(bot *tgbotapi.BotAPI) {
	fact, _ := getFact()
	spew.Dump(fact)
	msg := tgbotapi.NewMessage(CHATID, fact)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
