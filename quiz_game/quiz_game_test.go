package quiz_game

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"
)

const baseTimeDuration = time.Duration(1) * time.Second

func TestReadQuestions(t *testing.T) {
	source, want := dummyQuestions()
	got, err := readQuestions(&source)

	if err != nil {
		t.Error("expected err to be nil, got " + err.Error())
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestNewQuiz(t *testing.T) {
	source, questions := dummyQuestions()
	want := Quiz{questions: questions, timeLimit: baseTimeDuration, score: 0}
	got, err := NewQuiz(&source, time.Duration(2)*time.Second, false)

	if err != nil {
		t.Error("expected err to be nil, got " + err.Error())
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestQuiz_Start(t *testing.T) {
	t.Run("successful completion", func(t *testing.T) {
		_, questions := dummyQuestions()
		quiz := Quiz{questions: questions, timeLimit: baseTimeDuration, score: 0}

		dummyInput := "4\n2\n20"
		funcDefer, err := mockStdin(t, dummyInput)
		if err != nil {
			t.Fatal(err)
		}
		defer funcDefer()

		want := Quiz{questions: questions, timeLimit: baseTimeDuration, score: 2}
		got, err := quiz.Start()

		if err != nil {
			t.Error("expected err to be nil, got " + err.Error())
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("times up", func(t *testing.T) {
		_, questions := dummyQuestions()
		quiz := Quiz{questions: questions, timeLimit: baseTimeDuration, score: 0}

		dummyInput := "4\n"
		funcDefer, err := mockStdin(t, dummyInput)
		if err != nil {
			t.Fatal(err)
		}
		defer funcDefer()

		_, err = quiz.Start()

		// TODO: dig more into this
		// not sure why this does not work, apparently Scanln returns EOF as an error
		// don't know why for now
		if err == nil {
			t.Error("expected err to be not nil")
		} else if err.Error() != "Times up! You scored 1 out of 3" {
			t.Errorf("want \"Times up! You scored 1 out of 3\", got %v", err)
		}
	})
}

func dummyQuestions() (bytes.Buffer, []Question) {
	var buffer bytes.Buffer
	buffer.WriteString("2+2,4\n5-3,2\n8*5,40")

	questions := []Question{
		{
			description: "2+2",
			answer:      "4",
		},
		{
			description: "5-3",
			answer:      "2",
		},
		{
			description: "8*5",
			answer:      "40",
		},
	}

	return buffer, questions
}

func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin

	tmpFile, err := os.CreateTemp(t.TempDir(), "mockTestInput.*")
	if err != nil {
		return nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpFile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpFile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}, nil
}
