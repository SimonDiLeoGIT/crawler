package crawl

import (
	"sync"
)

type job struct {
	url   string
	depth int
}

func Worker(id int, jobs <-chan job, results chan<- job, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		Fetcher(j, results)
		done <- struct{}{} // signal: this job is fully done, all its links reported
	}
}
