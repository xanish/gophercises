package main

import (
	"flag"
	"fmt"
	quizgame "github.com/xanish/gophercises/quiz_game"
	"os"
	"time"
)

func main() {
	runQuizGame()
}

func runQuizGame() {
	csvFileNamePtr := flag.String("file", "problems.csv", "CSV file containing questions for quiz")
	timeLimitPtr := flag.Int("time", 30, "Time duration of the quiz")
	shouldShufflePtr := flag.Bool("shuffle", false, "Shuffle the questions in quiz")

	f, err := os.Open(*csvFileNamePtr)
	if err != nil {
		panic(err)
	}

	// since we are reading the file I think it is fine to
	// ignore the error from f.Close(), however, it would be
	// good to handle it correctly in case of writing
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	quiz, err := quizgame.NewQuiz(f, time.Duration(*timeLimitPtr)*time.Second, *shouldShufflePtr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Press ENTER to begin quiz")
	_, err = fmt.Scanln()
	if err != nil {
		panic(err)
	}

	quiz, err = quiz.Start()
	if err != nil {
		panic(err)
	}

	fmt.Println(quiz)
}
