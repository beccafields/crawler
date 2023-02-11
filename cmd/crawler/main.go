package main

import (
	"log"
	"os"

	"example.com/crawler/crawler"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal(crawler.ErrNoURL)
	}
	url := os.Args[1]

	log.Println("Crawling page:", url)
	links, err := crawler.CrawlURL(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Here are the links for", url)
	crawler.PrettyPrint(links)
}
