package crawl

import (
	"fmt"
	"log"
	urls "net/url"
	"sync"

	"golang.org/x/net/publicsuffix"
)

const maxDepth = 2
const maxNumWorkers = 16

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

func Crawl(seed string) int {

	InitWorkers(maxNumWorkers)

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
