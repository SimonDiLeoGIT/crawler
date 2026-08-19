package crawl

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var results = make(chan fetchResult, 100)
var jobs = make(chan job, 100)

func InitWorkers(client *http.Client, ctx context.Context, numWorkers int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobs, results, &wg, ctx, client)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return &wg
}

func Crawl(seed string, depth, numWorkers int) {
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &http.Client{Timeout: 10 * time.Second}

	InitWorkers(client, ctx, numWorkers)

	pendingQueue := []job{}
	pendingQueue = append(pendingQueue, job{url: seed, depth: 0})

	visited := map[string]bool{seed: true}

	pending := 1
	discoveredUrls := 0
	crawledUrls := 0
	jobsClosed := false

loop:
	for {
		if len(pendingQueue) == 0 {
			select {
			case r, ok := <-results:
				if !ok {
					break loop
				}
				handleResult(r, seed, depth, &pendingQueue, &pending, visited, &discoveredUrls)
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
				handleResult(r, seed, depth, &pendingQueue, &pending, visited, &discoveredUrls)
				if pending == 0 && !jobsClosed {
					close(jobs)
					jobsClosed = true
				}
			case jobs <- pendingQueue[0]:
				crawledUrls++
				if !jobsClosed {
					pendingQueue = pendingQueue[1:]
				}
			}
		}
	}

	elapsed := time.Since(start)

	fmt.Println("")
	fmt.Println("Crawl completed\n")
	fmt.Printf("URLs discovered %d\n", discoveredUrls)
	fmt.Printf("URLs crawled %d\n", crawledUrls)
	//fmt.Printf("URLs skiped %d\n", discoveredUrls-crawledUrls)
	fmt.Println("")
	fmt.Printf("Depth: %d\n", depth)
	fmt.Printf("Workers: %d\n", numWorkers)
	fmt.Printf("Duration: %v\n", elapsed)
}

func handleResult(r fetchResult, seed string, depth int, pendingQueue *[]job, pending *int, visited map[string]bool, discoveredUrls *int) {
	if r.finished {
		*pending--
		return
	}
	if !visited[r.j.url] && IsInSeedDomain(r.j.url, seed) {
		visited[r.j.url] = true
		*discoveredUrls++
		fmt.Println(r.j.url)
		if r.j.depth < depth {
			*pending++
			*pendingQueue = append(*pendingQueue, r.j)
		}
	}
}
