package crawler

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MakeDefaultClient provides a default HTTP client
func MakeDefaultClient() HTTPClient {
	return &http.Client{}
}

/*
CrawlURL takes a string URL, fetches the data from it and extracts all the valid URLs within that data.

It requires a http client to be supplied, in order to make a GET request.

It returns a slice of the URLs found and an error.
*/
func CrawlURL(url string, client HTTPClient) ([]string, error) {
	pageData, err := getPageData(url, client)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(pageData))
	if err != nil {
		return nil, err
	}

	links := []string{}
	extractLinks(doc, url, &links)

	return links, nil
}

// getPageData fetches the body data of the given URL and returns it as a string.
func getPageData(url string, client HTTPClient) (string, error) {

	// GET requests don't work with www.example.com URLs
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := (client).Do(request)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("get request didn't return a 200 status - got " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// extractLinks takes a *html.Node and traverses it's children, checking for href tags and storing
// any URLs found in the given slice.
func extractLinks(node *html.Node, argURL string, links *[]string) {
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
				*links = (append(*links, parsedURL.Host+parsedURL.Path+parsedURL.RawQuery))
				break
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, argURL, links)
	}
}
