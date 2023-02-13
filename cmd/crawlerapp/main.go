package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beccafields/crawler/crawler"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Wrong number of arguments found - see usage below\n Provide only one URL to crawl from - e.g. ./main example.com")
	}
	url := os.Args[1]

	log.Println("Crawling page:", url)
	client := crawler.MakeDefaultClient()
	links, err := crawler.CrawlURL(url, client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Here are the links for", url)
	prettyPrint(links)
}

// prettyPrint is a helper func to print out all values of a slice on new lines
func prettyPrint(slice []string) {
	for _, s := range slice {
		fmt.Println(s)
	}
}
