package main

import (
	"context"
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

// GetCheckFormLogin retrieves the checkFormLogin parameter that gets sent to
// https://www.deezer.com/ajax/action.php when logging in. This is used for
// enabling auto-logins.
func (c *PrivateClient) GetCheckFormLogin() (string, error) {
	ctx := context.Background()
	v := url.Values{}
	v.Set("method", "deezer.getUserData")
	resp, err := c.get(ctx, v, nil)
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
	// errors come in a map or empty array. Need to address inconsistent types
	//if pb.Error.GatewayError != "" {
	//	return "", fmt.Errorf(pb.Error.GatewayError)
	//}
	if pb.Results.CheckFormLogin == "" {
		return "", fmt.Errorf("checkFormLogin value is empty")
	}
	return pb.Results.CheckFormLogin, nil
}

//func (c *PrivateClient) GetCheckFormLogin() (l string, err error) {
//	// Create URL string
//	u, err := NewDefaultPrivateAPIURL()
//	if err != nil {
//		return "", err
//	}
//	q := u.Query()
//	q.Set("method", "deezer.getUserData")
//	u.RawQuery = q.Encode()
//
//	// Create request and add headers
//	req, err := http.NewRequest("GET", u.String(), nil)
//	if err != nil {
//		return "", err
//	}
//	req = SetDefaultHeaders(req)
//
//	// Send request and read body
//	resp, err := c.httpClient.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//
//	// Unmarshal
//	var pb PrivateBody
//	err = json.Unmarshal(body, &pb)
//	if err != nil {
//		return "", err
//	}
//	if pb.Results.CheckFormLogin == "" {
//		return "", fmt.Errorf("checkFormLogin value is empty")
//	}
//	return pb.Results.CheckFormLogin, nil
//}
