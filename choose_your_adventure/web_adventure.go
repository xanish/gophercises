package choose_your_adventure

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type WebHandler struct {
	story    Story
	template *template.Template
}

func (h WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path string
	if r.URL.Path == "/" {
		path = "intro"
	} else {
		path = strings.TrimLeft(r.URL.Path, "/")
	}

	log.Printf("%s %s", r.Method, path)

	arc := h.story[path]
	_ = h.template.Execute(w, arc)
}

func Web(storyReader io.Reader, templatePath string) (WebHandler, error) {
	story, err := ParseJSON(storyReader)
	if err != nil {
		return WebHandler{}, fmt.Errorf("could not load story: %w", err)
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return WebHandler{}, fmt.Errorf("could not read story template: %w", err)
	}

	return WebHandler{story: story, template: t}, nil
}
