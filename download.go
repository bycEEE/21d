package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/bogem/id3v2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

const privateKey = "jo6aey6haid2Teih"
const secret = "g4el58wc0zvf9na1"

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
			path := filepath.Join(".", "downloads", track.Data.Artists[0].Name, track.Data.AlbumTitle, fmt.Sprintf("%s - %s.mp3", track.Data.TrackNumber, track.Data.SongTitle))
			DownloadTrack(path, track, format)
		}
	},
}

//func getPath()

// DownloadTrack will download the track
func DownloadTrack(path string, track *PrivateTrack, format string) error {
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

	// create directory and file
	os.MkdirAll(filepath.Dir(path), 0777)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// create buffer and keep track of how many iterations we've completed, every 3rd 2048 byte block is encrypted
	size, _ := strconv.Atoi(resp.Header["Content-Length"][0])
	count := 0
	cur := 0
	filebuf := make([]byte, size) // file buffer
	fmt.Printf("Downloading %s, size: %d bytes\n", path, size)

	for cur < size {
		var partsize int
		if size - cur >= 2048 {
			partsize = 2048
		} else {
			partsize = size - cur
		}
		partbuf := make([]byte, partsize)
		_, err := io.ReadFull(resp.Body, partbuf)
		if err != nil {
			return err
		}
		if (count % 3 > 0) || partsize < 2048 {
			copy(filebuf[cur:cur+partsize], partbuf[:])
		} else {
			decrypted := decryptChunk(partbuf[:], key)
			copy(filebuf[cur:cur+partsize], decrypted)
		}
		cur += partsize
		count++
	}
	f.Write(filebuf) // write to file
	WriteTag(path, track)
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
func decryptChunk(ct, key []byte) []byte {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// make initialisation vector 8 bytes zero to seven instead of first 8 bytes of ciphertext
	// this is a static iv
	iv := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	// check last slice of encrypted text, if it's not a modulus of cipher block size, we're in trouble
	if len(ct[blowfish.BlockSize:])%blowfish.BlockSize != 0 {
		panic("decrypted chunk is not a multiple of blowfish.BlockSize")
	}
	decrypted := ct[:]
	// ok, we're good... create the decrypteaws_cloudwatch_event_target"r
	dcbc := cipher.NewCBCDecrypter(block, iv)
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

// WriteTag will write the id3v2 metadata tags
func WriteTag(file string, track *PrivateTrack) error {
	tag, err := id3v2.Open(file, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("error opening mp3 file", err)
	}
	tag.AddTextFrame(tag.CommonID("Title"), tag.DefaultEncoding(), track.Data.SongTitle)
	tag.AddTextFrame(tag.CommonID("Lead artist/Lead performer/Soloist/Performing group"), tag.DefaultEncoding(), track.Data.Artists[0].Name)
	tag.AddTextFrame(tag.CommonID("Album/Movie/Show title"), tag.DefaultEncoding(), track.Data.AlbumTitle)
	//tag.AddTextFrame(tag.CommonID("Band/Orchestra/Accompaniment"), tag.DefaultEncoding(), track.Data.)
	tag.AddTextFrame(tag.CommonID("Track number/Position in set"), tag.DefaultEncoding(), track.Data.TrackNumber)
	tag.AddTextFrame(tag.CommonID("Part of a set"), tag.DefaultEncoding(), track.Data.DiskNumber)
	tag.AddTextFrame(tag.CommonID("ISRC"), tag.DefaultEncoding(), track.Data.ISRC)
	//tag.AddTextFrame(tag.CommonID("Length"), tag.DefaultEncoding(), track.Data.)
	//tag.AddTextFrame(tag.CommonID("Attached picture"), tag.DefaultEncoding(), track.Data.)
	//tag.AddTextFrame(tag.CommonID("Unsynchronised lyrics/text transcription"), tag.DefaultEncoding(), track.Lyrics.LyricsText)
	//tag.AddTextFrame(tag.CommonID("Publisher"), tag.DefaultEncoding(), track.Data.)
	//tag.AddTextFrame(tag.CommonID("Genre"), tag.DefaultEncoding(), track.Data.)
	tag.AddTextFrame(tag.CommonID("Copyright message"), tag.DefaultEncoding(), track.Data.Copywrite)
	//tag.AddTextFrame(tag.CommonID("Date"), tag.DefaultEncoding(), track.Data.PhysicalReleaseDate)
	//tag.AddTextFrame(tag.CommonID("Year"), tag.DefaultEncoding(), track.Data.)
	//tag.AddTextFrame(tag.CommonID("BPM"), tag.DefaultEncoding(), track.Data.BPM)
	//tag.AddTextFrame(tag.CommonID("Composer"), tag.DefaultEncoding(), track.Data.)
	if err = tag.Save(); err != nil {
		log.Fatal("error while saving tag", err)
	}
	return nil
}
