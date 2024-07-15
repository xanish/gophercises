package choose_your_adventure

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ErrAdventureEnded = errors.New("adventure ended")

type CliHandler struct {
	story  Story
	reader *bufio.Reader
}

func (h CliHandler) ServeCLI(path string, out io.Writer) error {
	if path == "" {
		path = "intro"
	}
	arc := h.story[path]

	_, _ = fmt.Fprintf(out, "\n")
	_, _ = fmt.Fprintf(out, "-----------------------------------------------------------\n")
	_, _ = fmt.Fprintf(out, "%s\n", arc.Title)
	_, _ = fmt.Fprintf(out, "-----------------------------------------------------------\n")
	_, _ = fmt.Fprintf(out, "\n")

	for _, desc := range arc.Description {
		_, _ = fmt.Fprintf(out, "%s\n", desc)
	}

	if len(arc.Options) == 0 {
		_, _ = fmt.Fprintf(out, "\n--------------------------The End--------------------------\n")
		return ErrAdventureEnded
	}

	_, _ = fmt.Fprintf(out, "\nYour choices are:\n")
	for i, option := range arc.Options {
		_, _ = fmt.Fprintf(out, "%d) %s\n", i+1, option.Description)
	}
	_, _ = fmt.Fprintf(out, "")

	var choice int
	for {
		_, _ = fmt.Fprintf(out, "What will you do?\n")

		input, err := h.reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintf(out, "Could not read your choice, please retry...\n")
			continue
		}

		choice, err = strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			_, _ = fmt.Fprintf(out, "Invalid choice %d, choice should be between %d and %d\n", choice, 1, len(arc.Options))
			continue
		}

		if choice <= 0 || choice > len(arc.Options) {
			_, _ = fmt.Fprintf(out, "Invalid choice %d, choice should be between %d and %d\n", choice, 1, len(arc.Options))
			continue
		}

		break
	}

	err := h.ServeCLI(arc.Options[choice-1].Arc, out)
	if err != nil {
		return err
	}

	return nil
}

func CLI(storyReader io.Reader, choiceReader io.Reader) (CliHandler, error) {
	story, err := ParseJSON(storyReader)
	if err != nil {
		return CliHandler{}, fmt.Errorf("could not load story: %w", err)
	}

	return CliHandler{story, bufio.NewReader(choiceReader)}, nil
}
