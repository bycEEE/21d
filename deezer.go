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

// Login handles logging in.
func (c *PrivateClient) Login(username, password string) (*http.Response, error) {
	req, err := http.NewRequest("POST", privateAPIURL, nil)
	if err != nil {
		return nil, err
	}
	checkFormLogin, err := c.GetCheckFormLogin()
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Add("type", "login")
	form.Add("mail", username)
	form.Add("password", password)
	form.Add("checkFormLogin", checkFormLogin)
	req.PostForm = form
	req = c.addHeaders(req, headers{"Content-Type": {"application/x-www-form-url-encoded"}})
	resp, err := c.client.Do(req)
	return resp, nil
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

