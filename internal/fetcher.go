package crawl

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Fetcher(j job, results chan<- job) {
	res, err := http.Get(j.url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Could not load html from %s\n", j.url)
		log.Printf("Error: %s\n", err)
	}

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		href, exists := s.Attr("href")
		if exists && len(href) > 0 {
			normalizedUrl := NormalizeURL(href, j.url)
			results <- job{url: normalizedUrl, depth: j.depth + 1}
		}
	})
}

func NormalizeURL(href string, host string) string {
	if !strings.HasPrefix(href, "https://") && !strings.HasPrefix(href, "http://") {
		if !strings.HasPrefix(href, "/") {
			href = host
		}
		parts := strings.Split(href, "/")
		if len(parts) > 1 {
			if strings.Contains(host, parts[1]) && href != "/" {
				href = host[:strings.Index(host, parts[1])] + "/" + href
			}
		}
		href = host
	}
	parts := strings.Split(href, "?")
	response := ""
	if len(parts) > 1 {
		response = parts[0]
	}
	if response != "" {
		parts = strings.Split(response, "#")
	} else {
		parts = strings.Split(href, "#")
	}
	if len(parts) > 1 {
		response = parts[0]
	}
	if response != "" {
		return response
	}
	return href
}
