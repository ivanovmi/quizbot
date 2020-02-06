package quiz

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// EnQuizDBURL is URL for get eng questions
const EnQuizDBURL = "https://opentdb.com/api.php?amount=1"

// RuQuizDBURL is URL for russian questions
const RuQuizDBURL = "https://engine.lifeis.porn/api/millionaire.php?q="

// CHATID is id of chat for direct conversation
var CHATID, _ = strconv.Atoi(os.Getenv("CHAT_ID"))

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

// EnQuestion is a structure containing info about questions - answers, difficulty, etc
type EnQuestion struct {
	Category        string   `json:"category"`
	Type            string   `json:"type"`
	Difficulty      string   `json:"difficulty"`
	Question        string   `json:"question"`
	CorrectAnswer   string   `json:"correct_answer"`
	IncorrectAnswer []string `json:"incorrect_answers"`
}

// QuestionsList is a general list of questions
type QuestionsList struct {
	ResponseCode int          `json:"response_code"`
	Data         []EnQuestion `json:"results"`
}

// QuestionData - question itself
type QuestionData struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

// RuQuestion from DB
type RuQuestion struct {
	ID              string       `json:"id"`
	Ok              bool         `json:"ok"`
	Data            QuestionData `json:"data"`
	Category        string
	Difficulty      string
	Question        string
	CorrectAnswer   string
	IncorrectAnswer []string
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

// GetQuestions from QuizDBURL and parse it into QuestionList struct
func GetQuestions() (*QuestionsList, error) {
	var q QuestionsList
	resp, err := http.Get(EnQuizDBURL)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	err = json.Unmarshal(b, &q)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	if q.ResponseCode != 0 {
		return nil, fmt.Errorf("server is unavailable")
	}
	return &q, nil
}

// GetEnQuestion is single question mode
func GetEnQuestion() (*EnQuestion, error) {
	ql, err := GetQuestions()
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	return &ql.Data[0], nil
}

// GetRuQuestion is for ru-lang questions
func GetRuQuestion() (*RuQuestion, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	lvl := RuLevels[r.Intn(len(RuLevels))]
	resp, err := http.Get(fmt.Sprint(RuQuizDBURL, lvl))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	defer resp.Body.Close()
	var q RuQuestion
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	err = json.Unmarshal(b, &q)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	q.Question = q.Data.Question
	q.CorrectAnswer = q.Data.Answers[0]
	q.IncorrectAnswer = q.Data.Answers[1:len(q.Data.Answers)]
	q.Difficulty = RuLevelsMap[lvl]
	return &q, nil
}

func getIndex(sl []string, el string) int {
	for i := 0; i < len(sl); i++ {
		if sl[i] == el {
			return i
		}
	}
	return -1
}

// SendEnMsg is for sending msg with eng question
func SendEnMsg(bot *tgbotapi.BotAPI) {
	q, err := GetEnQuestion()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
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

// SendRuMsg is for sending msg with ru question
func SendRuMsg(bot *tgbotapi.BotAPI) {
	q, err := GetRuQuestion()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
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
