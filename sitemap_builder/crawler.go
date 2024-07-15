package sitemap_builder

import (
	"fmt"
	"github.com/xanish/gophercises/html_link_parser"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Crawl(seedUrl string, maxDepth int, out io.Writer) ([]string, error) {
	visited := make(map[string]bool)
	pending := make([]html_link_parser.Link, 1)
	pending[0] = html_link_parser.Link{Href: seedUrl}

	// extract the base url from the seedUrl to filter out
	// targets that should not be crawled
	parsedUrl, err := url.Parse(seedUrl)
	fmt.Print(parsedUrl.String())
	if err != nil {
		return nil, err
	}
	baseUrl := fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)

	for i := 0; i < maxDepth; i++ {
		// make a list of all urls to crawl at current depth and
		// copy all pending ones to it
		processing := make([]html_link_parser.Link, len(pending))

		// stop if we can't process anymore
		if len(processing) == 0 {
			break
		}
		copy(processing, pending)

		// reset the pending crawls
		pending = make([]html_link_parser.Link, 0)

		log.Printf("at depth %d, pending %d and processing %d\n", i, len(pending), len(processing))

		// keep working until we have any url left to crawl
		for len(processing) > 0 {
			// get rid of any trailing "/" characters so avoid
			// crawling /abc and /abc/ multiple times
			// leading "/" is fine since we use it to figure out
			// if the url is from the current site
			currentUrl := strings.TrimRight(processing[0].Href, "/")
			if currentUrl == "" {
				currentUrl = baseUrl
			} else if strings.HasPrefix(currentUrl, "/") {
				currentUrl = baseUrl + currentUrl
			}

			// pop out the crawled entry from processing queue
			processing = processing[1:]

			// if the url is already crawled then pop it out of
			// the processing queue and move to next
			if _, ok := visited[currentUrl]; ok {
				continue
			}

			log.Printf("starting crawl for %s\n", currentUrl)

			// extract all links available on the page to crawl
			links, err := extract(currentUrl)
			if err != nil {
				log.Printf("encountered %v\n", err)
				continue
			}

			// filter out the urls that were discovered now and should
			// be crawled at next depth
			pending = append(pending, filter(links, func(s string) bool {
				return strings.HasPrefix(s, baseUrl) || strings.HasPrefix(s, "/")
			})...)

			// mark the url as visited to avoid re-crawling
			visited[currentUrl] = true
		}
	}

	crawledUrls := make([]string, 0, len(visited))
	for k, _ := range visited {
		crawledUrls = append(crawledUrls, k)
	}

	return crawledUrls, nil
}

func extract(url string) ([]html_link_parser.Link, error) {
	// ignore http redirects
	// if CheckRedirect returns ErrUseLastResponse, then the most recent
	// response is returned with its body unclosed, along with a nil error
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	links, err := html_link_parser.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func filter(links []html_link_parser.Link, filterFn func(string) bool) []html_link_parser.Link {
	var filtered []html_link_parser.Link
	for _, link := range links {
		if filterFn(link.Href) {
			filtered = append(filtered, link)
		}
	}

	return filtered
}
