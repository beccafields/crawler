package crawler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrNoURL       = errors.New("no URL given to crawl through - please provide a URL")
	ErrNoBodyFound = errors.New("no body returned from page GET request")
)

/*CrawlURL takes a string URL, fetches the data from it and extracts all the valid URLs within that data.

Only URLs of the form "http://example" or "https://example" are accepted - "www.example" is not

It returns a map where the keys are the URLs found and the values are always true bools, and an error.*/
func CrawlURL(argURL string) (map[string]bool, error) {

	_, err := url.ParseRequestURI(argURL)
	if err != nil {
		return nil, err
	}

	pageData, err := getPageData(argURL)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(pageData))
	if err != nil {
		return nil, err
	}

	links := make(map[string]bool)
	extractLinks(doc, &links)

	return links, nil
}

// PrettyPrint is a helper func to print out all keys of a given map on new lines
func PrettyPrint[M ~map[K]V, K comparable, V any](m M) {
	for k := range m {
		fmt.Println(k)
	}
}

// getPageData fetches the body data of the given URL and returns it as a string.
// If no data is found, ErrNoBodyFound is returned. 
func getPageData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", ErrNoBodyFound
	}
	return string(body), nil
}

// extractLinks takes a *html.Node and traverses it's children, checking for href tags and saving
// any URLs found to the given map.
func extractLinks(node *html.Node, links *map[string]bool) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			_, err := url.ParseRequestURI(a.Val)
			if a.Key == "href" && err == nil {
				(*links)[a.Val] = true
				break
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, links)
	}
}
