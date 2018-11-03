package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const privateKey = "jo6aey6haid2Teih"
const secret = "g4el58wc0zvf9na1"
const chunkSize = 2048

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
				log.Fatalf("downloading track failed: %+v", err)
			}
			// check to see if file exists
			// remember MP3_256 is an int for some reason and not a string
			if (format == "3" && track.Data.FileSizeMP3320 == "0") || (format == "5" && track.Data.FileSizeMP3256 == 0) || (format == "1" && track.Data.FileSizeMP3128 == "0" || (format == "9" && track.Data.FileSizeFLAC == "0")) {
				log.Fatalf("track not found in the desired bitrate") // add song title and bitrate here
			}
			DownloadTrack("test.mp3", track, format)
		}
	},
}

//func getPath()

// DownloadTrack will download the track
func DownloadTrack(path string, track*PrivateTrack, format string) error {
	// get download and file size
	u, err := GetDownloadURL(track.Data.SongID, track.Data.MD5Origin, format, track.Data.MediaVersion)
	if err != nil {
		return err
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed with error code %d", resp.StatusCode)
	}
	if len(resp.Header["Content-Length"]) < 1 {
		return fmt.Errorf("no content length provided by header")
	}

	// get decryption key
	key := GetBlowfishKey(track.Data.SongID)

	// create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// create buffer and keep track of how many iterations we've completed, every 3rd 2048 byte block is encrypted
	size, _ := strconv.Atoi(resp.Header["Content-Length"][0])
	count := 0
	cur := 0
	filebuf := make([]byte, size) // file buffer
	partbuf := make([]byte, chunkSize) // part buffer
	fmt.Printf("Downloading %s, size: %d bytes\n", path, size)
	for {
		bytelen, err := io.ReadFull(resp.Body, partbuf)
		if err != nil { // err when fewer bytes than 2048 are read
			lastbuf := make([]byte, size-cur)
			copy(filebuf[cur:size], lastbuf[:])
			cur = size
			break
		}
		// decrypt every third 2048 byte block
		if (count % 3 != 0) || bytelen < chunkSize {
			copy(filebuf[cur:cur+bytelen], partbuf[:bytelen])
		} else {
			copy(filebuf[cur:], decryptChunk(partbuf, key))
		}
		cur += bytelen
		count++
	}
	f.Write(filebuf) // write to file
	return nil
}
//
//// DownloadTrack will download the track.
//func DownloadTrack(path string, track *PrivateTrack, format string) error {
//	chunkSize := 2048
//	size := 0
//	start := 0
//	end := 0
//
//	// get file size
//	u, err := GetDownloadURL(track.Data.SongID, track.Data.MD5Origin, format, track.Data.MediaVersion)
//	if err != nil {
//		return err
//	}
//	resp, err := http.Get(u.String())
//	if err != nil {
//		return err
//	}
//	if resp.StatusCode != 200 {
//		return fmt.Errorf("request failed with error code %d", resp.StatusCode)
//	}
//	if len(resp.Header["Content-Length"]) < 1 {
//		return fmt.Errorf("no content length provided by header")
//	}
//	size, _ = strconv.Atoi(resp.Header["Content-Length"][0])
//
//	// get decryption key
//	key := GetBlowfishKey(track.Data.SongID)
//
//	// create file and lock when writing to make it safe for concurrency
//	//sync.Mutex{}.Lock()
//	//defer sync.Mutex{}.Unlock()
//	f, err := os.Create(path)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	// download and send to file, every 3rd 2048 byte block
//	fmt.Printf("Downloading %s, size: %d\n", path, size)
//	for i := 0; i < int(size/chunkSize); i++ { // round down
//		end = start + chunkSize
//		getChunk(i, int64(start), int64(end), u, key, f) // offset index for i % 3 to return the correct remainder
//		start = end
//	}
//	// if the requested chunk is smaller than the chunk size, get remaining bytes
//	// last bytes are not encrypted
//	if end < size {
//		start = end
//		end = size
//		getChunk(1, int64(start), int64(end), u, key, f) // iter = 1 is just a placeholder since this won't be encrypted
//	}
//	return nil
//}
//
//// getChunk gets the requested chunk size. It will decrypt every 3rd 2048 byte block.
//func getChunk(iter int, start, end int64, u *url.URL, key []byte, f *os.File) error {
//	var client http.Client
//	req, err := http.NewRequest("GET", u.String(), nil)
//	if err != nil {
//		return err
//	}
//	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end-1)) // offset by 1 for index
//	resp, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	// check if server accepts partial content requests
//	if resp.StatusCode != 206 {
//		return fmt.Errorf("server does not accept partial content requests")
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//
//	// decrypt on third block
//	if iter % 3 == 0 {
//		f.WriteAt(decryptChunk(body, key), start)
//	} else {
//		f.WriteAt(body, start)
//	}
//	return nil
//}

// decryptChunk is used to decrypt tracks. Every 3rd 2048 byte block is encrypted. The key is retrieved from
// GetBlowfishKey.
func decryptChunk(et, key []byte) []byte {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	//make initialisation vector to be the first 8 bytes of ciphertext
	div := et[:blowfish.BlockSize]
	// check last slice of encrypted text, if it's not a modulus of cipher block size, we're in trouble
	decrypted := et[blowfish.BlockSize:]
	if len(decrypted)%blowfish.BlockSize != 0 {
		panic("decrypted chunk is not a multiple of blowfish.BlockSize")
	}
	// ok, we're good... create the decrypter
	dcbc := cipher.NewCBCDecrypter(block, div)
	// decrypt!
	dcbc.CryptBlocks(decrypted, decrypted)
	return decrypted
}

// GetBlowfishKey is a magic function that returns a blowfish key from the song ID.
func GetBlowfishKey(songID string) []byte {
	hash := md5.Sum([]byte(songID))
	checksum := hex.EncodeToString(hash[:])
	bfkey := make([]byte, 16)
	for i := 0; i < 16; i++ {
		bfkey[i] += checksum[i] ^ checksum[i+16] ^ secret[i]
	}
	return bfkey
}

// ECBEncrypt is used by GetDownloadURL to retrieve the download URL for a track using AES encryption
// with ECB and no padding. Decrypting is done using CBC.
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
