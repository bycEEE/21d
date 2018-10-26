package main

// PrivateBody contains the response body to a private API call.
type PrivateBody struct {
	// Error   []string       `json:"error"`
	Results PrivateResults `json:"results"`
}

// PrivateResults contains the results from the private API.
type PrivateResults struct {
	User           PrivateUser `json:"USER"`
	CheckFormLogin string      `json:"checkFormLogin"`
	CheckForm      string      `json:"checkForm"`
}

// PrivateUser contains information about the Deezer user retrieved from the
// private API
type PrivateUser struct {
	ID          int    `json:"USER_ID"`
	SessionID   string `json:"SESSION_ID"`
	UserToken   string `json:"USER_TOKEN"`
	Country     string `json:"COUNTRY"`
	PlayerToken string `json:"PLAYER_TOKEN"`
}
