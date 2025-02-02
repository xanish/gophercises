package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/xanish/gophercises/hacker_news"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./hacker_news/cmd/index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	sc := storyCache{
		numStories: numStories,
		duration:   6 * time.Second,
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   6 * time.Second,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-ticker.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	numStories int
	cache      []item
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if time.Now().Sub(sc.expiration) < 0 {
		return sc.cache, nil
	}

	stories, err := getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}
	sc.expiration = time.Now().Add(sc.duration)
	sc.cache = stories

	return sc.cache, nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hacker_news.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}

	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	jobsCh := make(chan int, len(ids))
	resultCh := make(chan jobResult)

	// start a pool with 6 workers
	for i := 0; i < 6; i++ {
		go spawnWorker(jobsCh, resultCh, i)
	}

	for _, storyId := range ids {
		jobsCh <- storyId
	}

	var results []jobResult
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}

	return stories
}

func spawnWorker(jobs <-chan int, results chan<- jobResult, workerId int) {
	i := 0
	for storyId := range jobs {
		fmt.Printf("worker %d starting story %d\n", workerId, storyId)
		go func(idx, id int) {
			var client hacker_news.Client
			hnItem, err := client.GetItem(id)
			if err != nil {
				results <- jobResult{idx: idx, err: err}
			}
			results <- jobResult{idx: idx, item: parseHNItem(hnItem)}
		}(i, storyId)

		i++
	}
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hacker_news.Item) item {
	ret := item{Item: hnItem}
	itemUrl, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(itemUrl.Hostname(), "www.")
	}

	return ret
}

// item is the same as the hacker_news.Item, but adds the Host field
type item struct {
	hacker_news.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type jobResult struct {
	idx  int
	item item
	err  error
}
