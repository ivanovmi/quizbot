package quiz

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	//	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ivanovmi/telegram-bot-api"
	"log"
	"net/http"
	"strings"
)

func getFavArticle() (*FavArticle, error) {
	res, err := http.Get(WikiURL)
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
	art := FavArticle{}
	doc.Find("#main-tfa").Each(func(i int, s *goquery.Selection) {
		href := s.Find(".mw-headline").Find("a")
		h, _ := href.Attr("href")
		art.Title = href.Text()
		art.URL = WikiURL + h
	})
	return &art, nil
}

func getDykArticle() (*[]string, error) {
	res, err := http.Get(WikiURL)
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
	artList := []string{}
	doc.Find("#main-dyk>ul").Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, s *goquery.Selection) {
			htmlText, _ := s.Html()
			htmlText = strings.ReplaceAll(htmlText, "<span class=\"nowrap\"><i>(на илл.)</i></span>", "")
			artList = append(artList, fmt.Sprintf("%s %s", bulb, htmlText))
		})
	})
	return &artList, nil
}

func sendFavArtMsg(bot *tgbotapi.BotAPI) {
	a, _ := getFavArticle()
	spew.Dump(a)
	text := fmt.Sprintf("%s Статья дня: %s\n%s %s", page, a.Title, link, a.URL)
	msg := tgbotapi.NewMessage(CHATID, text)
	bot.Send(msg)
}

func sendDykMsg(bot *tgbotapi.BotAPI) {
	dyk, _ := getDykArticle()
	spew.Dump(dyk)
	msg := tgbotapi.NewMessage(CHATID, strings.Join(*dyk, "\n"))
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func SendWikiMsg(bot *tgbotapi.BotAPI, t string) {
	switch t {
	case "art":
		sendFavArtMsg(bot)
	case "dyk":
		sendDykMsg(bot)
	}
}
