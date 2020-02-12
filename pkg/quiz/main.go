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
	"os"
	"strconv"
	"time"
)

const (
	// EnQuizURL is URL for questions in english
	EnQuizURL = "https://opentdb.com/api.php?amount=1"
	// RuQuizURL is URL for questions in russian
	RuQuizURL = "https://engine.lifeis.porn/api/millionaire.php?q="
	// RuTrueFalseURL is URL for True-False questions in russian
	RuTrueFalseURL = "https://engine.lifeis.porn/api/true_or_false.php"
)

// CHATID is id of chat for direct conversation
var CHATID, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)

// RuLevels is levels list
var RuLevels = [...]int{
	1, // Easy
	2, // Medium
	3, // Hard
	4, // Child
}

// RuLevelsMap is map of difficulty levels
var RuLevelsMap = map[int]string{
	1: "easy",
	2: "medium",
	3: "hard",
	4: "child",
}

// Question is a structure containing info about question
type Question struct {
	Category        string
	Difficulty      string
	Question        string
	CorrectAnswer   string
	IncorrectAnswer []string
}

type ruTrueFalseData struct {
	Ok   bool          `json:"ok"`
	Data []ruTrueFalse `json:"data"`
}

type ruTrueFalse struct {
	Fact   string `json:"fact"`
	IsTrue string `json:"is_true"`
}

type ruQuestionData struct {
	Ok   bool       `json:"ok"`
	Data ruQuestion `json:"data"`
}

type ruQuestion struct {
	Question   string   `json:"question"`
	Answers    []string `json:"answers"`
	Difficulty string
}

type enQuestionData struct {
	ResponseCode int          `json:"response_code"`
	Data         []enQuestion `json:"results"`
}

type enQuestion struct {
	Category        string   `json:"category"`
	Type            string   `json:"type"`
	Difficulty      string   `json:"difficulty"`
	Question        string   `json:"question"`
	CorrectAnswer   string   `json:"correct_answer"`
	IncorrectAnswer []string `json:"incorrect_answers"`
}

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

func getRuTrueFalseQuestion() (*ruTrueFalse, error) {
	var q ruTrueFalseData
	jsonBody, err := get(RuTrueFalseURL)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	err = json.Unmarshal(jsonBody, &q)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return &q.Data[0], nil
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
	case "ru-tf":
		var correctAnswer string
		var incorrectAnswer []string
		fmt.Println("ru-tf")
		question, err := getRuTrueFalseQuestion()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		if question.IsTrue == "1" {
			correctAnswer = "Правда"
			incorrectAnswer = []string{"Ложь"}
		} else {
			correctAnswer = "Ложь"
			incorrectAnswer = []string{"Правда"}
		}
		q = Question{
			Question:        question.Fact,
			CorrectAnswer:   correctAnswer,
			IncorrectAnswer: incorrectAnswer,
		}
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
	bot.Send(pollMsg)
}
