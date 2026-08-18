package main

import (
	crawl "crawler/internal"
	"fmt"
)

func main() {
	url := "http://bbc.com/"
	//url := "http://cnn.com"

	total := crawl.Crawl(url)

	fmt.Printf("Total urls %d\n", total)
	fmt.Println("All jobs processed successfully!")
}
