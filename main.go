package main

import (
	"fmt"
)

type job struct {
	url   string
	depth int
}

const maxDepth = 1
const maxNumWorkers = 3

func main() {
	url := "http://bbc.com/"

	total := Crawl(url, maxDepth, maxNumWorkers)

	fmt.Printf("Total urls %d\n", total)
	fmt.Println("All jobs processed successfully!")
}
