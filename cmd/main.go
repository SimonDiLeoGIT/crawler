package main

import (
	crawl "crawler/internal"
	"flag"
)

func main() {

	url := flag.String("url", "http://bbc.com/", "Crawler seed")
	depth := flag.Int("depth", 3, "Hops limit")
	numWorkers := flag.Int("workers", 8, "Num of workers")
	//url := "http://bbc.com/"
	//url := "http://cnn.com"

	// Parse parameters from the command line
	flag.Parse()

	crawl.Crawl(*url, *depth, *numWorkers)
}
