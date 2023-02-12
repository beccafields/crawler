package crawler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

// mockClient is used to mock the http.Do() request and return a constructed response for testing.
// Used in TestCrawlURL() and TestGetPageData()
type mockClient struct{}

func (mockClient) Do(req *http.Request) (*http.Response, error) {

	resp := &http.Response{}
	var err error

	switch {
	case strings.Contains(req.URL.String(), "successexample"):
		resp.Body = io.NopCloser(bytes.NewReader([]byte(`<p>Links:</p><ul><li><a href="https://example2.com">Example 2</a><li><a href="https://example3.com">Example 3</a></ul>`)))
		resp.StatusCode = http.StatusOK
		resp.Status = "200 OK"
		err = nil

	case strings.Contains(req.URL.String(), "emptybody"):
		resp.Body = io.NopCloser(bytes.NewReader([]byte("")))
		resp.StatusCode = http.StatusOK
		resp.Status = "200 OK"
		err = nil

	case strings.Contains(req.URL.String(), "failexample"):
		resp.Body = io.NopCloser(bytes.NewReader([]byte("")))
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = "500 Error"
		err = errors.New("could not make request")
	}

	return resp, err
}

// Full successful test of the crawl package - more detailed cases handled under the tests for helper funcs
func TestCrawlURL(t *testing.T) {
	cases := []struct {
		name          string
		url           string
		expectedLinks map[string]bool
		expectedErr   bool
	}{
		{
			"success full URL",
			"https://successexample.com",
			map[string]bool{
				"https://example2.com": true,
				"https://example3.com": true,
			},
			false,
		},
		{
			"success partial URL",
			"successexample.com",
			map[string]bool{
				"https://example2.com": true,
				"https://example3.com": true,
			},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := mockClient{}
			links, err := CrawlURL(c.url, client)
			if !c.expectedErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, c.expectedLinks, links)
		})
	}
}

func TestGetPageData(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		expectedData string
		expectedErr  bool
	}{
		{
			"success fetched data",
			"https://successexample.com",
			`<p>Links:</p><ul><li><a href="https://example2.com">Example 2</a><li><a href="https://example3.com">Example 3</a></ul>`,
			false,
		},
		{
			"fail empty page data",
			"https://emptybody.com",
			"",
			true,
		},
		{
			"fail invalid URL",
			"failexample.com",
			"",
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := mockClient{}
			pageData, err := getPageData(c.url, client)
			if !c.expectedErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, c.expectedData, pageData)
		})
	}
}

func TestExtractLinks(t *testing.T) {
	cases := []struct {
		name          string
		pageData      string
		argURL        string
		completeLinks map[string]bool
	}{
		{
			"success with full links found",
			`<p>Links:</p><ul><li><a href="https://example2.com">Example 2</a><li><a href="https://example3.com">Example 3</a></ul>`,
			"https://example.com",
			map[string]bool{
				"https://example2.com": true,
				"https://example3.com": true,
			},
		},
		{
			"success with relative links found",
			`<p>Links:</p><ul><li><a href="/relative1">Example 2</a><li><a href="/relative2">Example 3</a></ul>`,
			"https://example.com",
			map[string]bool{
				"https://example.com/relative1": true,
				"https://example.com/relative2": true,
			},
		},
		{
			"success with invalid links ignored",
			`<p>Links:</p><ul><li><a href="#main">Main</a><li><a href="https://example2.com">Example 2</a></ul>`,
			"https://example.com",
			map[string]bool{
				"https://example2.com": true,
			},
		},
		{
			"success with query strings handled",
			`<p>Links:</p><ul><li><a href="https://example2.com?this=added">Main</a><li></ul>`,
			"https://example.com",
			map[string]bool{
				"https://example2.com?this=added": true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			node, err := html.Parse(strings.NewReader(c.pageData))
			require.NoError(t, err)

			links := make(map[string]bool)
			extractLinks(node, c.argURL, &links)

			assert.Equal(t, c.completeLinks, links)
		})
	}
}
