package quiz

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	//	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ivanovmi/telegram-bot-api"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Shuffle an array
func Shuffle(vals []string) []string {
	rand.Seed(time.Now().UnixNano())
	for i := len(vals) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		vals[i], vals[j] = vals[j], vals[i]
	}
	return vals
}

func getIndex(sl []string, el string) int {
	for i := 0; i < len(sl); i++ {
		if sl[i] == el {
			return i
		}
	}
	return -1
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return b, nil
}

func getEnQuestion() (*enQuestion, error) {
	var q enQuestionData
	jsonBody, err := get(EnQuizURL)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	err = json.Unmarshal(jsonBody, &q)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return &q.Data[0], nil
}

func getRuQuestion() (*ruQuestion, error) {
	var q ruQuestionData
	r := rand.New(rand.NewSource(time.Now().Unix()))
	lvl := RuLevels[r.Intn(len(RuLevels))]
	jsonBody, err := get(fmt.Sprint(RuQuizURL, lvl))
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	err = json.Unmarshal(jsonBody, &q)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	q.Data.Difficulty = RuLevelsMap[lvl]
	return &q.Data, nil
}

func beautifyString(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.Split(s, "=")[1], "\"", ""), ";", "")
}

func getPddQuestion() *Question {
	var q Question
	q.Category = "ПДД"
	rand.Seed(time.Now().UnixNano())
	d, _ := getDoc(fmt.Sprintf("%s%s", PDDURL, "/?rand"))
	sel := d.Find("script")
	script := sel.Eq(1).Text()
	rIndex := rand.Intn(9)
	qPattern := fmt.Sprintf("(qL\\[%d\\]=\".*)", rIndex)
	aPattern := fmt.Sprintf("(aL\\[%d\\]\\[\\d\\]=\".*)", rIndex)
	caPattern := fmt.Sprintf("(tp\\[%d\\]\\[\\d\\]=\\d)", rIndex)
	imgPattern := fmt.Sprintf("(imgq\\[%d\\].*;n)", rIndex)
	qr, _ := regexp.Compile(qPattern)
	ar, _ := regexp.Compile(aPattern)
	car, _ := regexp.Compile(caPattern)
	imgr, _ := regexp.Compile(imgPattern)
	q.Question = beautifyString(qr.FindString(script))
	a := ar.FindAllString(script, -1)
	ca := car.FindAllString(script, -1)
	for i := range a {
		if strings.HasSuffix(ca[i], "1") {
			q.CorrectAnswer = beautifyString(a[i])
		} else {
			q.IncorrectAnswer = append(q.IncorrectAnswer, beautifyString(a[i]))
		}
	}
	img := imgr.FindString(script)
	if !strings.Contains(img, "<br>") {
		q.ImgURL = fmt.Sprintf("%s%s", PDDURL, strings.Split(strings.Fields(img)[3], "'")[1])
	}
	return &q
}

func getQuestion(t string) (*Question, error) {
	var q Question
	switch t {
	case "en":
		fmt.Println("en")
		question, err := getEnQuestion()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		q = Question{
			Category:        question.Category,
			Difficulty:      question.Difficulty,
			Question:        question.Question,
			CorrectAnswer:   question.CorrectAnswer,
			IncorrectAnswer: question.IncorrectAnswer,
		}
	case "ru":
		fmt.Println("ru")
		question, err := getRuQuestion()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		q = Question{
			Question:        question.Question,
			CorrectAnswer:   question.Answers[0],
			IncorrectAnswer: question.Answers[1:len(question.Answers)],
			Difficulty:      question.Difficulty,
		}
	case "pdd":
		fmt.Println("pdd")
		q = *getPddQuestion()
	}
	return &q, nil
}

// SendMsg is bot send message command
func SendMsg(bot *tgbotapi.BotAPI, t string) {
	q, _ := getQuestion(t)
	spew.Dump(q)
	var answers []string
	answers = append(answers, html.UnescapeString(q.CorrectAnswer))
	for answer := range q.IncorrectAnswer {
		answers = append(answers, html.UnescapeString(q.IncorrectAnswer[answer]))
	}
	Shuffle(answers)
	correctIndex := getIndex(answers, q.CorrectAnswer)
	pollMsg := tgbotapi.SendPollConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: CHATID,
		},
		Question:        fmt.Sprintf("%s (category: '%s', difficulty: %s)", html.UnescapeString(q.Question), q.Category, q.Difficulty),
		Options:         answers,
		IsAnonymous:     false,
		Type:            "quiz",
		CorrectOptionID: int64(correctIndex),
	}
	if q.ImgURL != "" {
		imgMsg := tgbotapi.NewMessage(CHATID, q.ImgURL)
		bot.Send(imgMsg)
	}
	bot.Send(pollMsg)
}
