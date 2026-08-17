package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
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

	visited := map[string]bool{seed: true}

	jobs <- job{url: seed, depth: 0}

	pending := 1
	total := 0
	pendingQueue := []job{}
	jobsClosed := false

loop:
	for {
		if len(pendingQueue) == 0 {
			select {
			case r, ok := <-results:
				if !ok {
					break loop
				}
				handleResult(r, &pendingQueue, &pending, visited, &total)
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
				handleResult(r, &pendingQueue, &pending, visited, &total)
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

func handleResult(r job, pendingQueue *[]job, pending *int, visited map[string]bool, total *int) {
	fmt.Println(r.url)
	*total++

	if r.depth < maxDepth {
		if !visited[r.url] {
			visited[r.url] = true
			*pending++
			*pendingQueue = append(*pendingQueue, job{r.url, r.depth + 1})
		}
	}
}

func Fetcher(j job, results chan<- job) {
	res, err := http.Get(j.url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		href, exists := s.Attr("href")
		if exists {
			if startsWithHTTPS(href) {
				results <- job{url: href, depth: j.depth + 1}
			}
		}
	})
}

func startsWithHTTPS(href string) bool {
	return strings.HasPrefix(href, "https://")
}
