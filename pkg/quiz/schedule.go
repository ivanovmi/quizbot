package quiz

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"strings"
	//	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ivanovmi/telegram-bot-api"
	"html"
)

const (
	// URL is a quizplease site URL
	URL    = "https://saratov.quizplease.ru"
	house  = "\xF0\x9F\x8F\xA0"
	link   = "\xF0\x9F\x94\x97"
	clock  = "\xF0\x9F\x95\x94"
	finger = "\xF0\x9F\x91\x89"
)

// Schedule is game schedule
type Schedule struct {
	Games []Game
}

// Game is just a game :)
type Game struct {
	Title string
	URL   string
	Place string
	Time  string
}

func getSchedule() (*Schedule, error) {
	res, err := http.Get("http://saratov.quizplease.ru/schedule")
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
	var games []Game
	doc.Find(".schedule-column .schedule-block").Each(func(i int, s *goquery.Selection) {
		var g Game
		title := s.Find(".h2-game-card").Text()
		title = strings.ReplaceAll(title, "[", "")
		title = strings.ReplaceAll(title, "]", "")
		g.Title = title
		href := s.Find("a[href].schedule-block-head")
		href.Each(func(i int, s *goquery.Selection) {
			h, _ := s.Attr("href")
			g.URL = URL + h
		})
		var info []string
		s.Find(".schedule-info").Find(".techtext").Each(func(i int, s *goquery.Selection) {
			info = append(info, strings.TrimSpace(s.Text()))
		})
		g.Place = fmt.Sprintf("%s (%s)", strings.Join(strings.Fields(info[0])[0:3], " "), info[1])
		g.Time = info[2]
		games = append(games, g)
	})
	schedule := Schedule{
		Games: games,
	}
	return &schedule, nil
}

// SendScheduleMsg for sending schedule updates
func SendScheduleMsg(bot *tgbotapi.BotAPI) {
	var games []string
	s, _ := getSchedule()
	spew.Dump(s)
	for _, g := range s.Games {
		spew.Dump(g)
		games = append(games, fmt.Sprintf("%s [%s](%s)\n%s %s\n%s %s", finger, html.UnescapeString(g.Title), g.URL, house, g.Place, clock, g.Time))
	}
	spew.Dump(games)
	msg := tgbotapi.NewMessage(CHATID, strings.Join(games, "\n\n"))
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
