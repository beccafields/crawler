package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/beccafields/crawler/crawler"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatal("Wrong number of arguments found - see usage below\n Provide one URL to crawl from and a limit of URLs to be returned - e.g. ./main example.com 50")
	}
	url := os.Args[1]
	limit, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("The given limit is not a valid int - see usage below\n Provide one URL to crawl from and a limit of URLs to be returned - e.g. ./main example.com 50")
	}

	log.Println("Crawling from page:", url)

	ch := make(chan string)
	totalVisited := []string{}

	client := crawler.MakeDefaultClient()
	initialSet, err := crawler.CrawlURL(url, client)
	if err != nil {
		log.Fatal(err)
	}

	// Split up the cralwing into one goroutine per URL returned from the starting page
	for _, url := range initialSet {
		go crawler.CrawlWeb(url, limit, ch)
	}

	// Retreive the successfully visited URLs from each of the goroutines
	for i := 0; i <= limit; i++ {
		visited := <-ch
		totalVisited = append(totalVisited, visited)
	}

	prettyPrint(totalVisited)
}

// prettyPrint is a helper func to print out all values of a slice on new lines
func prettyPrint(slice []string) {
	for _, s := range slice {
		fmt.Println(s)
	}
}
