package quiz

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"strings"
	//	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"github.com/ivanovmi/telegram-bot-api"
)

func getDoc(URL string) (*goquery.Document, error) {
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, nil
}

func parseRating(doc *goquery.Document) (*Rating, error) {
	var r Rating
	doc.Find(".rating-table-row").Each(func(i int, s *goquery.Selection) {
		s.Find(".hidden-xs").Each(func(i int, s *goquery.Selection) {
			r.Place = strings.TrimSpace(strings.ReplaceAll(s.Text(), ".", ""))
		})
		s.Find(".rating-table-kol-game").Each(func(i int, s *goquery.Selection) {
			r.Games = strings.TrimSpace(strings.ReplaceAll(s.Text(), "Игры", ""))
		})
		s.Find(".rating-table-points").Each(func(i int, s *goquery.Selection) {
			r.Points = strings.TrimSpace(strings.ReplaceAll(s.Text(), "Баллы", ""))
		})
		s.Find(".rating-table-rang-block").Each(func(i int, s *goquery.Selection) {
			rang, _ := s.Find("img").Attr("alt")
			r.Rang = ratingsMap[rang]
		})
	})
	return &r, nil
}

func getGlobalRating() (*Rating, error) {
	d, _ := getDoc(GlobalRatingURL)
	r, _ := parseRating(d)
	return r, nil
}

func getSeasonRating() (*Rating, error) {
	d, _ := getDoc(SeasonRatingURL)
	r, _ := parseRating(d)
	return r, nil
}

func getRating() (map[string]*Rating, error) {
	gRating, _ := getGlobalRating()
	sRating, _ := getSeasonRating()
	return map[string]*Rating{"Global": gRating, "Season": sRating}, nil
}

// SendRatingMsg for sending weekly team rating
func SendRatingMsg(bot *tgbotapi.BotAPI) {
	var ratings []string
	r, _ := getRating()
	spew.Dump(r)
	for t, rat := range r {
		ratings = append(ratings, fmt.Sprintf("%s rating:\nPlace: %s, after %s games, with %s points (Rang: %s)", t, rat.Place, rat.Games, rat.Points, rat.Rang))
	}
	msg := tgbotapi.NewMessage(CHATID, strings.Join(ratings, "\n"))
	bot.Send(msg)
}
