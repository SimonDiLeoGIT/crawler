package crawl

import (
	"context"
	"log"
	"net/http"
	urls "net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

func Fetcher(ctx context.Context, j job, results chan<- fetchResult) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, j.url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
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
		return
	}

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		href, exists := s.Attr("href")
		if exists && len(href) > 0 {
			normalizedUrl := NormalizeURL(href, j.url)
			results <- fetchResult{j: job{url: normalizedUrl, depth: j.depth + 1}}
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

func IsInSeedDomain(url, seed string) bool {
	parsedURL, err := urls.Parse(url)
	if err != nil {
		log.Println(err)
		return false
	}
	urlRootDomain, err := publicsuffix.EffectiveTLDPlusOne(parsedURL.Hostname())
	if err != nil {
		log.Println(err)
		return false
	}

	parsedSeed, err := urls.Parse(seed)
	if err != nil {
		log.Println(err)
		return false
	}
	seedRootDomain, err := publicsuffix.EffectiveTLDPlusOne(parsedSeed.Hostname())
	if err != nil {
		log.Println(err)
		return false
	}
	return urlRootDomain == seedRootDomain
}
