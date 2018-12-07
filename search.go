package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/url"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search",
}

var searchTrackCmd = &cobra.Command{
	Use:   "track",
	Short: "Search for track",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, query := range args {
			results, err := SearchTrack(query)
			if err != nil {
				log.Fatalf("searching track %s failed: %+v", query, err)
			} else {
				for i, track := range results {
					fmt.Printf("%d | %s (%d) by %s (%d) in %s (%d)\n", i, track.Title, track.ID, track.Artist.Name, track.Artist.ID, track.Album.Title, track.Album.ID)
				}
			}
		}
	},
}

func Search(objtype, query string) (*PublicResults, error) {
	c, err := NewPublicClient()
	if err != nil {
		return nil, err
	}
	u := &url.URL{Scheme: "https", Host: c.basePath, Path: fmt.Sprintf("search/%s", objtype)}
	q := u.Query()
	q.Set("q", fmt.Sprintf("%s:\"%s\"", objtype, query))
	u.RawQuery = q.Encode()

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results PublicResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func SearchTrack(query string) ([]PublicTrack, error) {
	tracks, err := Search("track", query)
	if err != nil {
		return nil, err
	}
	return tracks.Data, nil
}
