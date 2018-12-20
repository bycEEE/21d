package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Login handles logging in.
func (c *PrivateClient) Login(username, password, checkFormLogin string) (*http.Response, error) {
	// set parameters
	data := url.Values{}
	data.Add("type", "login")
	data.Add("mail", username)
	data.Add("password", password)
	data.Add("checkFormLogin", checkFormLogin)
	// create new request
	req, err := http.NewRequest("POST", "http://www.deezer.com/ajax/action.php", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	// add form header
	req = c.addHeaders(req, headers{"Content-Type": {"application/x-www-form-urlencoded"}})
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("request failed with error code %d, %s", resp.StatusCode, string(body))
	}
	// check if login succeeded
	if string(body) != "success" {
		return resp, fmt.Errorf( "invalid username/password combination")
	}

	// retrieve cookies from login and save if found
	if len(resp.Cookies()) < 1 {
		return resp, fmt.Errorf("no cookies found in login response")
	} else {
		err = c.jar.Save()
		if err != nil {
			return nil, err
		}
		fmt.Println("Cookies saved")
	}

	// get token
	v := url.Values{}
	v.Set("method", "deezer.getUserData")
	privateResp, err := c.GetPrivateResponse(v)
	if err != nil {
		return resp, err
	}
	if privateResp.Results.CheckForm == "" {
		return resp, fmt.Errorf("token value is empty")
	}
	f, err := os.OpenFile(".token", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return resp, err
	}
	defer f.Close()
	f.WriteString(privateResp.Results.CheckForm)
	fmt.Println("Token saved")

	return resp, nil
}

// GetPrivateResponse parses an http response from a GET request to the private API.
func (c *PrivateClient) GetPrivateResponse(v url.Values) (*PrivateResponse, error) {
	ctx := context.Background()
	resp, err := c.get(ctx, v, nil)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with error code %d, %s", resp.StatusCode, string(body))
	}

	//
	// implement error handling with PrivateError
	//

	// unmarshal
	var pr PrivateResponse
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, err
}

// PostPrivateResponse parses an http response from a POST request to the private API.
func (c *PrivateClient) PostPrivateResponse(v url.Values, body io.Reader) (*PrivateResponse, error) {
	ctx := context.Background()
	// read token if not set
	if v.Get("api_token") == "" {
		b, err := ioutil.ReadFile(".token")
		if err != nil {
			return nil, fmt.Errorf("error reading token, try logging in again: %+v", err)
		}
		v.Set("api_token", string(b))
	}
	resp, err := c.post(ctx, v, body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	// do error handling here when token expires
	// {"error":{"VALID_TOKEN_REQUIRED":"Invalid CSRF token"},"results":{},"payload":null}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with error code %d, %s", resp.StatusCode, string(respBody))
	}

	//
	// implement error handling with PrivateError
	//

	// unmarshal
	var pr PrivateResponse
	err = json.Unmarshal(respBody, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, err
}
