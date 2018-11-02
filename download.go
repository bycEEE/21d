package main

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/charmap"
	"log"
	"net/url"
)

const privateKey = "jo6aey6haid2Teih"

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "You must specify the type of resource to download",
}

var downloadTrackCmd = &cobra.Command{
	Use:   "track",
	Short: "Download track",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// format is an arbitrary number defined by Deezer to correspond to streaming quality
		var format string

		switch downloadQuality {
		case "MP3_320":
			format = "3"
		case "MP3_256":
			format = "5"
		case "MP3_128":
			format = "1"
		case "FLAC":
			format = "9"
		default:
			log.Fatalf("invalid download quality defined")
		}
		for _, id := range args {
			track, err := GetTrack(id)
			if err != nil {
				log.Fatalf("getting track failed: %+v", err)
			}
			downloadURL, err := GetDownloadURL(id, track.Data.MD5Origin, format, track.Data.MediaVersion)
			if err != nil {
				log.Fatalf("getting download url for track %s failed, %+v", id, err)
			}
			fmt.Println(downloadURL.String())
		}
	},
}

// ECBEncrypt is used by GetDownloadURL to retrieve the download URL for a track using AES encryption
// with ECB and no padding.
func ECBEncrypt(pt, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

//func DownloadTrack() {
//
//}

// GetDownloadURL is a magic function that constructs the download URL from SNG_ID,
// MD5_ORIGIN, format (streaming quality), and MEDIA_VERSION.
func GetDownloadURL(songID, md5Origin, format, mediaVersion string) (*url.URL, error) {
	data := md5Origin + "¤" + format + "¤" + songID + "¤" + mediaVersion
	// go encodes using utf-8 by default
	// we need to encode using latin1 iso-8859-1
	dataenc, _ := charmap.ISO8859_1.NewEncoder().Bytes([]byte(data))
	hash := md5.Sum([]byte(dataenc))
	checksum := hex.EncodeToString(hash[:])
	newdata := checksum + "¤" + data + "¤"
	newdataenc, _ := charmap.ISO8859_1.NewEncoder().Bytes([]byte(newdata))
	// no funny characters in our private key
	encrypted := ECBEncrypt([]byte(newdataenc), []byte(privateKey))

	u := "http://e-cdn-proxy-" + md5Origin[0:1] + ".deezer.com/mobile/1/" + hex.EncodeToString(encrypted)
	pu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return pu, nil
}
