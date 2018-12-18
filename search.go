package main

import (
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
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
				data := make([][]string, len(results))
				for i, track := range results {
					data[i] = []string{strconv.Itoa(track.ID), track.Title, strconv.Itoa(track.Artist.ID), track.Artist.Name, strconv.Itoa(track.Album.ID), track.Album.Title}
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Track ID", "Track Title", "Artist ID", "Artist Name", "Album ID", "Album Title"})
				table.SetRowLine(true)
				table.AppendBulk(data)
				table.Render()
			}
		}
	},
}

func SearchTrack(query string) (PublicTrackList, error) {
	c, err := NewPublicClient()
	if err != nil {
		return nil, err
	}
	u := &url.URL{Scheme: "https", Host: c.basePath, Path: "search/track"}
	q := u.Query()
	q.Set("q", query)
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
	var results PublicTrackListResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	return results.Data, nil
}

//func SearchArtist(query string) ([]PublicArtist, error) {
//	artists, err := Search("artist", query)
//	if err != nil {
//		return nil, err
//	}
//	return artists.Data, nil
//}
