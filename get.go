package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/url"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "You must specify the type of resource to get",
}

var getTrackCmd = &cobra.Command{
	Use:   "track",
	Short: "Get track info",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, id := range args {
			t, err := GetTrack(id)
			if err != nil {
				log.Fatalf("getting track %s failed: %+v", id, err)
			} else {
				fmt.Println(t)
			}
		}
	},
}

// GetTrack retrieves a track and its relevant info.
func GetTrack(songID string) (*PrivateTrack, error) {
	// create private client
	privateClient, err := NewPrivateClient()
	if err != nil {
		log.Fatalf("Error establishing connection to the private Deezer API: %+v", err)
	}
	// private API expects this in the body
	bodyVal := map[string]string{
		"sng_id": songID,
	}
	resp, err := privateClient.GetResource("deezer.pageTrack", songID, bodyVal)
	if err != nil {
		return nil, err
	}
	pt := PrivateTrack{Data: resp.Results.Data, Lyrics: resp.Results.Lyrics}
	if pt.Data.MD5Origin == "" {
		return nil, fmt.Errorf("failed to retrieve track MD5_ORIGIN, value is empty")
	}
	return &pt, nil
}

// GetResource is a wrapper around PostPrivateResponse that sends a post request to retrieve a resource.
func (c *PrivateClient) GetResource(method string, id string, bodyVal map[string]string) (*PrivateResponse, error) {
	v := url.Values{}
	v.Set("method", method)
	jsonVal, err := json.Marshal(bodyVal)
	if err != nil {
		return nil, err
	}
	resp, err := c.PostPrivateResponse(v, bytes.NewBuffer(jsonVal))
	if err != nil {
		return nil, fmt.Errorf("error retrieving resource, %s: %s, %+v", method, id, err)
	}
	return resp, nil
}
