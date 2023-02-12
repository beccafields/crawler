package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/beccafields/crawler/crawler"
)

func init() {
	crawler.Client = &http.Client{}
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Wrong number of arguments found - see usage below\n Provide only one URL to crawl from - e.g. ./main example.com")
	}
	url := os.Args[1]

	log.Println("Crawling page:", url)
	links, err := crawler.CrawlURL(url, crawler.Client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Here are the links for", url)
	prettyPrint(links)
}

// prettyPrint is a helper func to print out all keys of a given map on new lines
func prettyPrint(m map[string]bool) {
	for k := range m {
		fmt.Println(k)
	}
}
