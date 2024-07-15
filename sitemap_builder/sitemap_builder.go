package sitemap_builder

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type urlset struct {
	Xmlns string `xml:"xmlns,attr"`
	Urls  []loc  `xml:"url"`
}

type loc struct {
	Value string `xml:"loc"`
}

func BuildSitemap(urls []string, out io.Writer) {
	siteMap := urlset{Xmlns: xmlns}
	for _, url := range urls {
		siteMap.Urls = append(siteMap.Urls, loc{Value: url})
	}

	_, err := fmt.Fprint(out, xml.Header)
	if err != nil {
		log.Fatalf("could not write to output: %v", err)
	}

	enc := xml.NewEncoder(out)
	enc.Indent("", "  ")
	if err = enc.Encode(siteMap); err != nil {
		log.Fatalf("could not write to output: %v", err)
	}

	_, err = fmt.Fprint(out, "\n")
	if err != nil {
		log.Fatalf("could not write to output: %v", err)
	}
}
