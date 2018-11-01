package main

import (
	jar "21d/cookiejar"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PrivateClient is a client that connects to the private API.
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

type headers map[string][]string

// NewPrivateClient returns a new client which talks to the private Deezer API.
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

// getAPIPath creates the URL to query.
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

func main() {
	//// create config file if it does not exist
	//os.OpenFile("config.yaml", os.O_RDONLY|os.O_CREATE, 0666)
	//
	//// read config for login credentials
	//viper.SetConfigFile("config.yaml")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")
	//viper.WatchConfig()
	//err := viper.ReadInConfig()

	// prompt user to log in if credentials are not found or relog if credentials expired
	//if err != nil {
	//	log.Fatalf("Config error, type `21d login` to continue: %s\n", err)
	//} else {
	//	// prompt for credentials and save if non-existent
	//	if !viper.IsSet("deezer.username") || !viper.IsSet("deezer.password") {
	//		log.Fatalf("Deezer login credentials were not found, type `21d` login to continue")
	//	}
	//}

	//// retrieve credentials
	//username := viper.GetString("deezer.username")
	//// on top of aes-256 encryption, the password was base64 encoded
	//encodedPassword, err := base64.StdEncoding.DecodeString(viper.GetString("deezer.password"))
	//if err != nil {
	//	log.Fatalf("Error decoding password: %+v", err)
	//}
	//password, err := decryptCredentials(encodedPassword, localKey)
	//
	//// create private client
	//privateClient, err := NewPrivateClient()
	//if err != nil {
	//	log.Fatalf("Error establishing connection to the private Deezer API: %+v", err)
	//}

	// login
	//resp, err := privateClient.Login(username, string(password))
	//fmt.Println(resp)

	// test getting userdata
	c, _ := NewPrivateClient()
	v := url.Values{}
	v.Set("method", "deezer.getUserData")
	resp, _ := c.GetPrivateResponse(v)
	fmt.Println(resp)
	// initiate cli commands
	Execute()

	//checkFormLogin, err := privateClient.GetCheckFormLogin()
	//if err != nil {
	//	log.Fatalf("Error getting checkFormLogin value: %+v\n", err)
	//}
}