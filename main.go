package main

import (
	"fmt"
)

type job struct {
	url   string
	depth int
}

const maxDepth = 4
const maxNumWorkers = 16

func main() {
	url := "http://bbc.com/"
	//url := "http://cnn.com"

	total := Crawl(url, maxDepth, maxNumWorkers)

	fmt.Printf("Total urls %d\n", total)
	fmt.Println("All jobs processed successfully!")
}
