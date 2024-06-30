package url_shortener

import (
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	fallbackPath          = "/test"
	fallbackNotFoundPath  = "/test-not-found"
	fallbackUrl           = "https://google.com"
	fallbackBody          = "Hello, world!"
	defaultDatabaseBucket = "redirectionUrls"
)

func TestMapHandler(t *testing.T) {
	pathsToUrls := map[string]string{
		fallbackPath: fallbackUrl,
	}
	mapHandler := MapHandler(pathsToUrls, http.HandlerFunc(fallbackHandler))

	t.Run("should return fallback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackNotFoundPath, nil)
		response := httptest.NewRecorder()

		mapHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusOK)
		assertBody(t, result, fallbackBody)
	})

	t.Run("should return redirected route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackPath, nil)
		response := httptest.NewRecorder()

		mapHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusPermanentRedirect)
		assertUrl(t, result, pathsToUrls[fallbackPath])
	})
}

func TestYAMLHandler(t *testing.T) {
	paths := `- path: /test
  url: http://google.com`
	mapHandler, err := YAMLHandler([]byte(paths), http.HandlerFunc(fallbackHandler))
	if err != nil {
		t.Fatalf("could not create YAMLHandler %v", err)
	}

	t.Run("should return fallback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackNotFoundPath, nil)
		response := httptest.NewRecorder()

		mapHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusOK)
		assertBody(t, result, fallbackBody)
	})

	t.Run("should return redirected route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackPath, nil)
		response := httptest.NewRecorder()

		mapHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusPermanentRedirect)
		assertUrl(t, result, fallbackUrl)
	})
}

func TestJSONHandler(t *testing.T) {
	paths := `[
  {
    "path": "/test",
    "url": "http://google.com"
  }
]`
	mapHandler, err := JSONHandler([]byte(paths), http.HandlerFunc(fallbackHandler))
	if err != nil {
		t.Fatalf("could not create JSONHandler %v", err)
	}

	t.Run("should return fallback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackNotFoundPath, nil)
		response := httptest.NewRecorder()

		mapHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusOK)
		assertBody(t, result, fallbackBody)
	})

	t.Run("should return redirected route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackPath, nil)
		response := httptest.NewRecorder()

		mapHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusPermanentRedirect)
		assertUrl(t, result, fallbackUrl)
	})
}

func TestDatabaseHandler(t *testing.T) {
	db, dbCloser := setupDatabase(t)
	databaseHandler := DatabaseHandler(db, http.HandlerFunc(fallbackHandler))
	defer dbCloser(db)

	t.Run("should return fallback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackNotFoundPath, nil)
		response := httptest.NewRecorder()

		databaseHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusOK)
		assertBody(t, result, fallbackBody)
	})

	t.Run("should return redirected route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fallbackPath, nil)
		response := httptest.NewRecorder()

		databaseHandler(response, request)

		result := response.Result()

		assertStatusCode(t, result, http.StatusPermanentRedirect)
		assertUrl(t, result, fallbackUrl)
	})
}

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, fallbackBody)
}

func assertStatusCode(t *testing.T, response *http.Response, expected int) {
	t.Helper()

	if response.StatusCode != expected {
		t.Errorf("expected status to be %d, got %d", expected, response.StatusCode)
	}
}

func assertUrl(t *testing.T, response *http.Response, expected string) {
	t.Helper()

	url, err := response.Location()
	if err != nil {
		t.Fatal("could not read location from response", err)
	}

	if url.String() != expected {
		t.Errorf("expected url to be %s, got %s", url, expected)
	}
}

func assertBody(t *testing.T, response *http.Response, expected string) {
	t.Helper()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		t.Fatal("could not read body from response", err)
	}

	got := string(body)
	if expected != got {
		t.Errorf("expected response body to be %s, got %s", expected, got)
	}
}

func setupDatabase(t *testing.T) (*bolt.DB, func(db *bolt.DB)) {
	db, err := bolt.Open(tempFile(), 0666, nil)
	if err != nil {
		t.Fatal(err)
	} else if db == nil {
		t.Fatal("expected db")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(defaultDatabaseBucket))
		if err != nil {
			t.Fatalf("create bucket failed: %s", err)
		}

		err = b.Put([]byte(fallbackPath), []byte(fallbackUrl))
		if err != nil {
			t.Fatalf("put key in bucket failed: %s", err)
		}

		return nil
	})

	if err != nil {
		t.Fatal("could not create bucket")
	}

	return db, func(db *bolt.DB) {
		_ = db.Close()
	}
}

func tempFile() string {
	f, err := os.CreateTemp("", "bolt-*")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
