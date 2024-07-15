package choose_your_adventure

import (
	"bytes"
	"errors"
	approvals "github.com/approvals/go-approval-tests"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseJSON(t *testing.T) {
	story := dummyStory()
	got, err := ParseJSON(&story)
	if err != nil {
		t.Fatal(err)
	}

	want := Story{}
	want["intro"] = Arc{
		Title:       "intro",
		Description: []string{"intro-story"},
		Options: []Choice{
			{
				Description: "go-to-arc-1",
				Arc:         "arc1",
			},
			{
				Description: "go-to-arc-2",
				Arc:         "arc2",
			},
		},
	}
	want["arc1"] = Arc{
		Title:       "arc1",
		Description: []string{"arc1-story"},
		Options:     []Choice{},
	}
	want["arc2"] = Arc{
		Title:       "arc2",
		Description: []string{"arc2-story"},
		Options: []Choice{
			{
				Description: "go-to-intro",
				Arc:         "intro",
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestWeb(t *testing.T) {
	readBytes := dummyStory()
	filePath, _ := filepath.Abs("./templates/main.gohtml")

	handler, err := Web(&readBytes, filePath)
	if err != nil {
		t.Fatalf("failed to create web handler: %v", err)
	}

	t.Run("should return intro", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		got := response.Body.String()

		approvals.VerifyString(t, got)
	})

	t.Run("should return arc2", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/arc2", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		got := response.Body.String()

		approvals.VerifyString(t, got)
	})
}

func TestCLI(t *testing.T) {
	t.Run("should traverse intro to arc2 to intro to arc1 to end", func(t *testing.T) {
		buf := bytes.Buffer{}
		readBytes := dummyStory()

		dummyInput := "2\n1\n1\n"
		inputReader, funcDefer, err := mockStdin(t, dummyInput)
		if err != nil {
			t.Fatal(err)
		}
		defer funcDefer()

		handler, err := CLI(&readBytes, inputReader)
		if err != nil {
			t.Fatalf("failed to create cli handler: %v", err)
		}

		err = handler.ServeCLI("", &buf)
		if !errors.Is(err, ErrAdventureEnded) {
			t.Fatalf("failed to execute cli: %v", err)
		}
		approvals.VerifyString(t, buf.String())
	})
}

func dummyStory() bytes.Buffer {
	var buffer bytes.Buffer

	buffer.WriteString(`{
	"intro": {
		"title": "intro", 
		"story": ["intro-story"], 
		"options": [
			{"text": "go-to-arc-1", "arc": "arc1"},
			{"text": "go-to-arc-2", "arc": "arc2"}
		]
	},
	"arc1": {
		"title": "arc1", 
		"story": ["arc1-story"], 
		"options": []
	},
	"arc2": {
		"title": "arc2", 
		"story": ["arc2-story"], 
		"options": [
			{"text": "go-to-intro", "arc": "intro"}
		]
	}
}`)

	return buffer
}

func mockStdin(t *testing.T, dummyInput string) (tmpFile *os.File, funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin

	tmpFile, err = os.CreateTemp(t.TempDir(), "mockTestInput.*")
	if err != nil {
		return nil, nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpFile.Write(content); err != nil {
		return nil, nil, err
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpFile

	return tmpFile, func() {
		// clean up
		os.Stdin = oldOsStdin
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}, nil
}
