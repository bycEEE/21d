package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func configure(cmd *cobra.Command, args []string) error {
	// prompt for credentials and save if non-existent
	// configure username
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Deezer username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())
	viper.Set("deezer.username", username)

	// configure password
	fmt.Print("Deezer password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	encryptedPassword, err := encryptCredentials(password, localKey)
	if err != nil {
		return err
	}
	// encode as base64 to easily store into config
	viper.Set("deezer.password", base64.StdEncoding.EncodeToString(encryptedPassword))

	// save credentials
	viper.WriteConfig()
	return nil
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure 21d",
	RunE:  configure,
}

var rootCmd = &cobra.Command{
	Use:   "21d",
	Short: "21d is a tool to search and download tracks from Deezer",
}

// Execute executes commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
