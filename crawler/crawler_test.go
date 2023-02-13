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

	case strings.Contains(req.URL.String(), "not200example"):
		resp.Body = io.NopCloser(bytes.NewReader([]byte("")))
		resp.StatusCode = http.StatusCreated
		resp.Status = "201 Error"
		err = nil
	}

	return resp, err
}

func TestCrawlURL(t *testing.T) {
	cases := []struct {
		name            string
		url             string
		expectedLinks   []string
		expectedErrText string
	}{
		{
			"success on full URL",
			"https://successexample.com",
			[]string{
				"https://example2.com",
				"https://example3.com",
			},
			"",
		},
		{
			"success on empty page",
			"www.emptybody.com",
			[]string{},
			"",
		},
		{
			"success on partial URL",
			"successexample.com",
			[]string{
				"https://example2.com",
				"https://example3.com",
			},
			"",
		},
		{
			"fail on empty URL",
			"",
			[]string(nil),
			"get request didn't return a 200 status - got 0",
		},
		{
			"fail on invalid URL",
			"failexample.com.com",
			[]string(nil),
			"could not make request",
		},
		{
			"fail on get request status",
			"https://not200example.co.uk",
			[]string(nil),
			"get request didn't return a 200 status - got 201",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := mockClient{}
			links, err := CrawlURL(c.url, client)
			if c.expectedErrText == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), c.expectedErrText)
			}
			assert.Equal(t, c.expectedLinks, links)
		})
	}
}

func TestGetPageData(t *testing.T) {
	cases := []struct {
		name            string
		url             string
		expectedData    string
		expectedErrText string
	}{
		{
			"success fetched data",
			"https://successexample.com",
			`<p>Links:</p><ul><li><a href="https://example2.com">Example 2</a><li><a href="https://example3.com">Example 3</a></ul>`,
			"",
		},
		{
			"success empty page data",
			"https://emptybody.com",
			"",
			"",
		},
		{
			"fail invalid URL",
			"failexample.com",
			"",
			"could not make request",
		},
		{
			"fail no error 201 status",
			"https://not200example.com",
			"",
			"get request didn't return a 200 status - got 201",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := mockClient{}
			pageData, err := getPageData(c.url, client)
			if c.expectedErrText == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), c.expectedErrText)
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
		completeLinks []string
	}{
		{
			"success with full links found",
			`<p>Links:</p><ul><li><a href="https://example2.com">Example 2</a><li><a href="https://example3.com">Example 3</a></ul>`,
			"https://example.com",
			[]string{
				"https://example2.com",
				"https://example3.com",
			},
		},
		{
			"success with relative links found",
			`<p>Links:</p><ul><li><a href="/relative1">Example 2</a><li><a href="/relative2">Example 3</a></ul>`,
			"https://example.com",
			[]string{
				"https://example.com/relative1",
				"https://example.com/relative2",
			},
		},
		{
			"success with invalid links ignored",
			`<p>Links:</p><ul><li><a href="#main">Main</a><li><a href="https://example2.com">Example 2</a></ul>`,
			"https://example.com",
			[]string{
				"https://example2.com",
			},
		},
		{
			"success with query strings handled",
			`<p>Links:</p><ul><li><a href="https://example2.com?this=added">Main</a><li></ul>`,
			"https://example.com",
			[]string{
				"https://example2.com?this=added",
			},
		},
		{
			"no URLs found",
			`<p>Links:</p><ul><li><a>Main</a><li></ul>`,
			"https://example.com",
			[]string{},
		},
		{
			"success same link found twice",
			`<p>Links:</p><ul><li><a href="https://example2.com">Example 2</a><li><a href="https://example2.com">Example 2</a></ul>`,
			"https://example.com",
			[]string{
				"https://example2.com",
				"https://example2.com",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			node, err := html.Parse(strings.NewReader(c.pageData))
			require.NoError(t, err)

			links := []string{}
			extractLinks(node, c.argURL, &links)

			assert.Equal(t, c.completeLinks, links)
		})
	}
}
