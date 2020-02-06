all:
	go build -o quizbot *.go

clean:
	rm -rf quizbot
