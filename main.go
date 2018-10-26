package main

import (
	"fmt"
	"log"
	"net/http"
)

// PrivateClient is a client that connects to the private API.
type PrivateClient struct {
	// client provides access to the original http.client functions
	httpClient *http.Client
}

// NewPrivateClient returns a new client which talks to the private Deezer API.
func NewPrivateClient() *PrivateClient {
	c := &PrivateClient{
		httpClient: &http.Client{},
	}
	return c
}

func main() {
	c := NewPrivateClient()
	checkFormLogin, err := c.GetCheckFormLogin()
	if err != nil {
		log.Fatalf("Error getting checkFormLogin value: %+v", err)
	}
	fmt.Print(checkFormLogin)
}
