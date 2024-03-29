package main

import (
	jar "21d/cookiejar"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

const (
	privateAPIURL = "http://www.deezer.com/ajax/gw-light.php"
	publicAPIURL  = "https://api.deezer.com"
)

type headers map[string][]string

// PrivateClient is a client that connects to the private API. The private API is mainly used to handle logging in and
// retrieving hidden track information such as `MD5_ORIGIN` to generate download links.
type PrivateClient struct {
	// basePath is the API host, this gets prepended to every request.
	basePath string
	// path is the url path.
	path string
	// client provides access to the original http.client functions.
	client *http.Client
	// headers are custom set headers.
	headers map[string]string
	// jar is the cookie jar
	jar *jar.Jar
}

//
// Private client
//

// NewPrivateClient returns a new client which talks to the private Deezer API. The private client sends default headers
// with each request and uses a local cookie jar for authentication.
func NewPrivateClient() (*PrivateClient, error) {
	// create private client
	c := &PrivateClient{}

	// create new cookie jar with persistent storage
	// loads cookies if found
	jarOpts := jar.Options{Filename: ".cookies"}
	var err error
	c.jar, err = jar.New(&jarOpts)
	// return error if cookies are not loaded properly
	if err != nil {
		return nil, err
	}

	// set url related info
	u, err := url.Parse(privateAPIURL)
	if err != nil {
		return nil, err
	}
	c.basePath = u.Host
	c.path = u.Path

	// set default headers
	c.headers = map[string]string{
		"User-Agent": "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Content-Language": "en-US",
		"Cache-Control": "max-age=0",
		"Accept": "*/*",
		"Accept-Charset": "utf-8,ISO-8859-1;q=0.7,*;q=0.3",
		"Accept-Language": "en-US,en;q=0.9,en-US;q=0.8,en;q=0.7",
	}

	c.client = &http.Client{Jar: c.jar}
	return c, nil
}

// getAPIPath constructs the URL to query. Each request requires the `api_version`, `input`, and `api_token` fields and
// are built into each URL.
func (c *PrivateClient) getAPIPath(query url.Values) string {
	// set default query values if not specified
	if query.Get("api_version") == "" {
		query.Set("api_version", "1.0")
	}
	if query.Get("input") == "" {
		query.Set("input", "3")
	}
	if query.Get("api_token") == "" {
		query.Set("api_token", "null")
	}
	return (&url.URL{Path: c.path, RawQuery: query.Encode()}).String()
}

// addHeaders is called when building a request to add headers
func (c *PrivateClient) addHeaders(req *http.Request, h headers) *http.Request {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	if h != nil {
		for k, v := range h {
			req.Header[k] = v
		}
	}
	return req
}

func (c *PrivateClient) buildRequest(method, path string, body io.Reader, headers headers) (*http.Request, error) {
	expectedPayload := (method == "POST" || method == "PUT")
	if expectedPayload && body == nil {
		body = bytes.NewReader([]byte{})
	}
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req = c.addHeaders(req, headers)
	req.URL.Host, req.URL.Scheme = c.basePath, "http"
	return req, nil
}

func (c *PrivateClient) sendRequest(ctx context.Context, method string, query url.Values, body io.Reader, headers headers) (*http.Response, error) {
	req, err := c.buildRequest(method, c.getAPIPath(query), body, headers)
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *PrivateClient) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *PrivateClient) get(ctx context.Context, query url.Values, headers map[string][]string) (*http.Response, error) {
	return c.sendRequest(ctx, "GET", query, nil, headers)
}

func (c *PrivateClient) post(ctx context.Context, query url.Values, body io.Reader, headers map[string][]string) (*http.Response, error) {
	return c.sendRequest(ctx, "POST", query, body, headers)
}

//
// Public client
//

// PublicClient is a client that connects to the public API. The public API is mainly used to retrieve search results.
type PublicClient struct {
	// basePath is the API host, this gets prepended to every request.
	basePath string
	// client provides access to the original http.client functions.
	client *http.Client
	// headers are custom set headers.
	headers map[string]string
}

// NewPublicClient returns a new client which talks to the public Deezer API
func NewPublicClient() (*PublicClient, error) {
	c := &PublicClient{}

	u, err := url.Parse(publicAPIURL)
	if err != nil {
		return nil, err
	}
	c.basePath = u.Host
	// set default headers
	c.headers = map[string]string{
		"User-Agent": "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Content-Language": "en-US",
		"Cache-Control": "max-age=0",
		"Accept": "*/*",
		"Accept-Charset": "utf-8,ISO-8859-1;q=0.7,*;q=0.3",
		"Accept-Language": "en-US,en;q=0.9,en-US;q=0.8,en;q=0.7",
	}

	c.client = &http.Client{}
	return c, nil
}

func (c *PublicClient) getAPIPath(path string, query url.Values) string {
	return (&url.URL{Path: path, RawQuery: query.Encode()}).String()
}

// addHeaders is called when building a request to add headers
func (c *PublicClient) addHeaders(req *http.Request, h headers) *http.Request {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	if h != nil {
		for k, v := range h {
			req.Header[k] = v
		}
	}
	return req
}

func (c *PublicClient) buildRequest(method, path string, body io.Reader, headers headers) (*http.Request, error) {
	expectedPayload := (method == "POST" || method == "PUT")
	if expectedPayload && body == nil {
		body = bytes.NewReader([]byte{})
	}
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req = c.addHeaders(req, headers)
	req.URL.Host, req.URL.Scheme = c.basePath, "https"
	return req, nil
}

func (c *PublicClient) sendRequest(ctx context.Context, method string, path string, query url.Values, body io.Reader, headers headers) (*http.Response, error) {
	req, err := c.buildRequest(method, c.getAPIPath(path, query), body, headers)
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *PublicClient) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *PublicClient) get(ctx context.Context, path string, query url.Values, headers map[string][]string) (*http.Response, error) {
	return c.sendRequest(ctx, "GET", path, query, nil, headers)
}

func (c *PublicClient) post(ctx context.Context, path string, query url.Values, body io.Reader, headers map[string][]string) (*http.Response, error) {
	return c.sendRequest(ctx, "POST", path, query, body, headers)
}