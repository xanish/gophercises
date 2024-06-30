package quiz_game

import (
	"encoding/csv"
	"io"
	"math/rand"
	"strings"
	"time"
)

func readQuestions(reader io.Reader) ([]Question, error) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		return []Question{}, err
	}

	questions := make([]Question, len(records))
	for i, record := range records {
		questions[i] = Question{
			description: strings.TrimSpace(record[0]),
			answer:      strings.TrimSpace(record[1]),
		}
	}

	return questions, nil
}

func NewQuiz(reader io.Reader, timeLimit time.Duration, shouldShuffle bool) (Quiz, error) {
	questions, err := readQuestions(reader)
	if err != nil {
		return Quiz{}, err
	}

	if shouldShuffle {
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}

	quiz := Quiz{questions: questions, timeLimit: timeLimit, score: 0}

	return quiz, nil
}
