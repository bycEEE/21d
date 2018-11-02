package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net/url"
	"os"
	"strings"
	"syscall"
)

func login(cmd *cobra.Command, args []string) error {
	// remove existing token and cookies
	os.Remove(".token")
	os.Remove(".cookie")

	// create private client and remove cookies from cookie jar, though none should be loaded
	privateClient, err := NewPrivateClient()
	if err != nil {
		log.Fatalf("Error establishing connection to the private Deezer API: %+v", err)
	}
	privateClient.jar.RemoveAll()
	// get required checkFormLogin value to send along with username/password
	v := url.Values{}
	v.Set("method", "deezer.getUserData")
	resp, err := privateClient.GetPrivateResponse(v)
	if err != nil {
		log.Fatalf("Error retrieving user data: %+v", err)
	}
	if resp.Results.CheckFormLogin == "" {
		log.Fatal("checkFormLogin value is empty")
	}

	// prompt for credentials and save cookie if non-existent
	// username
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Deezer username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())
	// password
	fmt.Print("Deezer password: ")
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(passwordBytes[:])
	fmt.Println("")
	if err != nil {
		return err
	}

	// verify login succeeded
	_, err = privateClient.Login(username, password, resp.Results.CheckFormLogin)
	if err != nil {
		return err
	}

	return nil
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into Deezer",
	RunE:  login,
}

var rootCmd = &cobra.Command{
	Use:   "21d",
	Short: "21d is a tool to search and download tracks from Deezer",
}

// Execute executes commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Global variables for flags.
var downloadQuality string

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getTrackCmd)
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.AddCommand(downloadTrackCmd)
	downloadCmd.PersistentFlags().StringVarP(&downloadQuality, "quality", "q", "MP3_320",
		"Select quality of downloads (default MP3_320). Valid values: MP3_128, MP3_256, MP3_320, FLAC")
}
