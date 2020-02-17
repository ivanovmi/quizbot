package quiz

import (
	"os"
	"strconv"
)

const (
	// EnQuizURL is URL for questions in english
	EnQuizURL = "https://opentdb.com/api.php?amount=1"
	// RuQuizURL is URL for questions in russian
	RuQuizURL = "https://engine.lifeis.porn/api/millionaire.php?q="
	// RuTrueFalseURL is URL for True-False questions in russian
	RuTrueFalseURL = "https://engine.lifeis.porn/api/true_or_false.php"
	// URL is a quizplease site URL
	URL    = "https://saratov.quizplease.ru"
	house  = "\xF0\x9F\x8F\xA0"
	link   = "\xF0\x9F\x94\x97"
	clock  = "\xF0\x9F\x95\x94"
	finger = "\xF0\x9F\x91\x89"
	cal    = "\xF0\x9F\x93\x85"
)

var (
	// CHATID is id of chat for direct conversation
	CHATID, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	// RuLevels is levels list
	RuLevels = [...]int{
		1, // Easy
		2, // Medium
		3, // Hard
		4, // Child
	}

	// RuLevelsMap is map of difficulty levels
	RuLevelsMap = map[int]string{
		1: "easy",
		2: "medium",
		3: "hard",
		4: "child",
	}
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
	Date  string
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
