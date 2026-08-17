package main

import (
	"fmt"
	"log"
	"net/http"
	urls "net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

var results = make(chan job, 100)
var jobs = make(chan job, 100)
var done = make(chan struct{}, 100)

func InitWorkers(numWorkers int) {
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobs, results, done, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
}

func Crawl(seed string, depth, numWorkers int) int {

	InitWorkers(numWorkers)

	pendingQueue := []job{}
	pendingQueue = append(pendingQueue, job{url: seed, depth: 0})

	visited := map[string]bool{seed: true}

	pending := 1
	total := 0
	jobsClosed := false

loop:
	for {
		if len(pendingQueue) == 0 {
			select {
			case r, ok := <-results:
				if !ok {
					break loop
				}
				handleResult(r, seed, &pendingQueue, &pending, visited, &total)
			case <-done:
				pending--
				if pending == 0 && !jobsClosed {
					close(jobs)
					jobsClosed = true
				}
			}
		} else {
			select {
			case r, ok := <-results:
				if !ok {
					break loop
				}
				handleResult(r, seed, &pendingQueue, &pending, visited, &total)
			case <-done:
				pending--
				if pending == 0 && !jobsClosed {
					close(jobs)
					jobsClosed = true
				}
			case jobs <- pendingQueue[0]:
				if !jobsClosed {
					pendingQueue = pendingQueue[1:]
				}
			}
		}
	}
	return total
}

func handleResult(r job, seed string, pendingQueue *[]job, pending *int, visited map[string]bool, total *int) {
	if !visited[r.url] && IsInSeedDomain(r.url, seed) {
		visited[r.url] = true
		*total++
		fmt.Println(r.url)
		if r.depth < maxDepth {
			*pending++
			*pendingQueue = append(*pendingQueue, job{r.url, r.depth})
		}
	}
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
