package crawl

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type job struct {
	url   string
	depth int
}

type fetchResult struct {
	j        job
	finished bool
}

func Worker(id int, jobs <-chan job, results chan<- fetchResult, wg *sync.WaitGroup, ctx context.Context, client *http.Client) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Work stopped. Reason: %v\n", ctx.Err())
			return
		case j, ok := <-jobs:
			if !ok {
				return
			}
			Fetcher(client, ctx, j, results)
			results <- fetchResult{finished: true}
		}
	}
}
