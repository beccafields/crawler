package crawler

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

/*
CrawlURL takes a string URL, fetches the data from it and extracts all the valid URLs within that data.

It requires a http client to be supplied, in order to make a GET request.

It returns a map where the keys are the URLs found and the values are always true bools, and an error.
*/
func CrawlURL(url string, client HTTPClient) (map[string]bool, error) {
	pageData, err := getPageData(url, client)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(pageData))
	if err != nil {
		return nil, err
	}

	links := make(map[string]bool)
	extractLinks(doc, url, &links)

	return links, nil
}

// getPageData fetches the body data of the given URL and returns it as a string.
func getPageData(url string, client HTTPClient) (string, error) {

	// GET requests don't work for www.example.com URLs
	if !strings.Contains(url, "https") {
		url = "https://" + url
	}

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := (client).Do(request)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("get request didn't return a 200 status - " + resp.Status)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", errors.New("no body returned from page GET request")
	}
	return string(body), nil
}

// extractLinks takes a *html.Node and traverses it's children, checking for href tags and saving
// any URLs found to the given map.
func extractLinks(node *html.Node, argURL string, links *map[string]bool) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			parsedURL, err := url.ParseRequestURI(a.Val)
			if a.Key == "href" && err == nil {
				// handle any relative links found
				if parsedURL.Host == "" {
					parsedURL.Host = argURL
				}
				if !strings.Contains(parsedURL.Host, "http") {
					parsedURL.Host = "https://" + parsedURL.Host
				}
				if parsedURL.RawQuery != "" {
					parsedURL.RawQuery = "?" + parsedURL.RawQuery
				}
				(*links)[parsedURL.Host+parsedURL.Path+parsedURL.RawQuery] = true
				break
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, argURL, links)
	}
}
