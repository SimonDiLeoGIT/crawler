package main

import (
	"fmt"
	"sync"
	"time"
)

func Worker(id int, jobs <-chan job, results chan<- job, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("Worker %d processing job: %s\n", id, j.url)
		Fetcher(j, results)
		done <- struct{}{} // signal: this job is fully done, all its links reported
		time.Sleep(200 * time.Millisecond)
	}
}
