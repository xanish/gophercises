package choose_your_adventure

import (
	"encoding/json"
	"io"
)

type Story map[string]Arc

type Arc struct {
	Title       string   `json:"title"`
	Description []string `json:"story"`
	Options     []Choice `json:"options"`
}

type Choice struct {
	Description string `json:"text"`
	Arc         string `json:"arc"`
}

func ParseJSON(r io.Reader) (Story, error) {
	story := Story{}
	d := json.NewDecoder(r)
	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}
