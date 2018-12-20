package main

import (
	"context"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
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
	Short: "Search for a track",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, query := range args {
			results, err := SearchTrack(query)
			if err != nil {
				log.Fatalf("searching track %s failed: %+v", query, err)
			} else {
				data := make([][]string, len(results.Data))
				for i, track := range results.Data {
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

var searchAlbumCmd = &cobra.Command{
	Use:   "album",
	Short: "Search for an album",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, query := range args {
			results, err := SearchAlbum(query)
			if err != nil {
				log.Fatalf("searching album %s failed: %+v", query, err)
			} else {
				data := make([][]string, len(results.Data))
				for i, album := range results.Data {
					data[i] = []string{strconv.Itoa(album.ID), album.Title, strconv.Itoa(album.Artist.ID), album.Artist.Name}
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Album ID", "Album Title", "Artist ID", "Artist Name"})
				table.SetRowLine(true)
				table.AppendBulk(data)
				table.Render()
			}
		}
	},
}

func SearchResponse(objtype, query string) (*http.Response, error) {
	c, err := NewPublicClient()
	if err != nil {
		return nil, err
	}
	path := "/search/" + objtype
	v := url.Values{}
	v.Set("q", query)
	v.Set("strict", "on")
	ctx := context.Background()
	resp, err := c.get(ctx, path, v, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SearchTrack(query string) (*PublicTrackListResults, error) {
	resp, err := SearchResponse("track", query)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var results PublicTrackListResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func SearchAlbum(query string) (*PublicAlbumListResults, error) {
	resp, err := SearchResponse("album", query)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var results PublicAlbumListResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

//func SearchArtist(query string) ([]PublicArtist, error) {
//	artists, err := Search("artist", query)
//	if err != nil {
//		return nil, err
//	}
//	return artists.Data, nil
//}
