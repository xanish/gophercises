package quiz_game

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Question struct {
	description string
	answer      string
}

func (q Question) String() string {
	return fmt.Sprintf("Q) %s", q.description)
}

type Quiz struct {
	questions []Question
	timeLimit time.Duration
	score     int
}

func (q Quiz) Start() (Quiz, error) {
	// not sure if either Start or String should use the reference to the quiz
	// maybe it's a good idea to modify the original quiz passed? but then again
	// it might just be better to take a copy of the original quiz and return
	// the updated one for the user taking it

	quizResult := make(chan Quiz)

	// not sure if creating an error channel is the right call here
	// just want to be able to notify the caller that an error happened
	// while we were trying to fetch user input
	quizError := make(chan error)

	// start the quiz flow in a goroutine, this way we can use a channel
	// to notify once the quiz has ended and use it along with a timer and
	// stop the quiz as soon as either of them is completed
	go func() {
		var answer string
		for _, question := range q.questions {
			fmt.Println(question)

			_, err := fmt.Scanln(&answer)
			if err != nil {
				quizError <- err
			}

			if strings.ToLower(strings.TrimSpace(answer)) == strings.ToLower(question.answer) {
				q.score++
			}
		}

		quizResult <- q
	}()

	// stop the quiz on either timeout or error or quiz ended successfully
	select {
	case <-time.After(q.timeLimit):
		return q, errors.New(fmt.Sprintf("Times up! %s", q))
	case err := <-quizError:
		return Quiz{}, err
	case <-quizResult:
		return q, nil
	}
}

func (q Quiz) String() string {
	return fmt.Sprintf("You scored %d out of %d", q.score, len(q.questions))
}
