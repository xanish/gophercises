package url_shortener

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"net/http"
)

type redirectionRecord struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		redirectTo, ok := pathsToUrls[path]
		if ok {
			http.Redirect(writer, request, redirectTo, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirectionRecords []redirectionRecord
	err := yaml.Unmarshal(yml, &redirectionRecords)
	if err != nil {
		return nil, err
	}

	pathToUrls := redirectionRecordsToMap(redirectionRecords)

	return MapHandler(pathToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//	{
//	   "path": "/urlshort",
//	   "url": "https://github.com/gophercises/urlshort"
//	 }
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(bytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirectionRecords []redirectionRecord
	err := json.Unmarshal(bytes, &redirectionRecords)
	if err != nil {
		return nil, err
	}

	pathToUrls := redirectionRecordsToMap(redirectionRecords)

	return MapHandler(pathToUrls, fallback), nil
}

// DatabaseHandler will read the url from provided Database and
// then return an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the Database, then the
// fallback http.Handler will be called instead.
func DatabaseHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Print(string(request.URL.Path))
		var url string
		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("redirectionUrls"))
			value := bucket.Get([]byte(request.URL.Path))
			if value != nil {
				url = string(value)
			}
			return nil
		})

		fmt.Println(err, url)

		if err == nil && url != "" {
			http.Redirect(writer, request, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

func redirectionRecordsToMap(redirectionRecords []redirectionRecord) map[string]string {
	pathToUrls := map[string]string{}
	for _, record := range redirectionRecords {
		pathToUrls[record.Path] = record.URL
	}

	return pathToUrls
}
