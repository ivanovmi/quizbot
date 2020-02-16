all:
	go build -o quizbot cmd/quizbot/main.go

clean:
	rm -rf quizbot
