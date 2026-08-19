package crawl

import (
	"context"
	"fmt"
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

func Worker(id int, jobs <-chan job, results chan<- fetchResult, wg *sync.WaitGroup, ctx context.Context) {
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
			Fetcher(ctx, j, results)
			results <- fetchResult{finished: true}
		}
	}
}
