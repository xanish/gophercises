package main

import (
	"flag"
	"github.com/xanish/gophercises/sitemap_builder"
	"log"
	"os"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "the maximum number of links deep to traverse")
	flag.Parse()

	crawled, err := sitemap_builder.Crawl(*urlFlag, *maxDepth, os.Stdout)
	if err != nil {
		log.Fatalf("Error building sitemap %v", err)
	}

	sitemap_builder.BuildSitemap(crawled, os.Stdout)
}
