package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	privateAPIURL = "http://www.deezer.com/ajax/gw-light.php"
	publicAPIURL  = "https://api.deezer.com"
)

// SetDefaultHeaders sets the default headers on a request to the private API.
// This should be set on all requests or reworked into the client.
func SetDefaultHeaders(r *http.Request) *http.Request {
	h := r.Header
	h.Set("User-Agent", "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0")
	h.Set("Content-Language", "en-US")
	h.Set("Cache-Control", "max-age=0")
	h.Set("Accept", "*/*")
	h.Set("Accept-Charset", "utf-8,ISO-8859-1;q=0.7,*;q=0.3")
	h.Set("Accept-Language", "en-US,en;q=0.9,en-US;q=0.8,en;q=0.7")
	return r
}

// NewDefaultPrivateAPIURL returns a URL type containing default parameters the
// Deezer Private API accepts. This is used to quickly set parameters
// such as `api_token` and `method`.
func NewDefaultPrivateAPIURL() (u *url.URL, err error) {
	u, err = url.Parse(privateAPIURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("api_version", "1.0")
	q.Set("api_token", "null")
	q.Set("input", "3")
	// q.Set("method", "null")
	u.RawQuery = q.Encode()
	return u, nil
}

// GetCheckFormLogin retrieves the checkFormLogin parameter that gets sent to
// https://www.deezer.com/ajax/action.php when logging in. This is used for
// enabling auto-logins.
func (c *PrivateClient) GetCheckFormLogin() (l string, err error) {
	// Create URL string
	u, err := NewDefaultPrivateAPIURL()
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("method", "deezer.getUserData")
	u.RawQuery = q.Encode()

	// Create request and add headers
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req = SetDefaultHeaders(req)

	// Send request and read body
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal
	var pb PrivateBody
	err = json.Unmarshal(body, &pb)
	if err != nil {
		return "", err
	}
	if pb.Results.CheckFormLogin == "" {
		return "", fmt.Errorf("checkFormLogin value is empty")
	}
	return pb.Results.CheckFormLogin, nil
}

// Deezer.prototype.getToken = async function(){
// 	const res = await request.get({
// 		url: this.apiUrl,
// 		headers: this.httpHeaders,
// 		strictSSL: false,
// 		qs: {
// 			api_version: "1.0",
// 			api_token: "null",
// 			input: "3",
// 			method: 'deezer.getUserData'
// 		},
// 		json: true,
// 		jar: true,
// 	})
// 	return res.body.results.checkForm;
