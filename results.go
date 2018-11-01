package main

// PrivateResponse contains the response body to a private API call.
type PrivateResponse struct {
	//Error   PrivateError     `json:"error"`
	Results PrivateResults `json:"results"`
}

// PrivateResults contains the results from the private API.
type PrivateResults struct {
	User           PrivateUser `json:"USER"`
	CheckFormLogin string      `json:"checkFormLogin"`
	CheckForm      string      `json:"checkForm"`
	SessionID   string `json:"SESSION_ID"`
	UserToken   string `json:"USER_TOKEN"`
	Country     string `json:"COUNTRY"`
	PlayerToken string `json:"PLAYER_TOKEN"`
}

// PrivateUser contains information about the Deezer user retrieved from the private API.
type PrivateUser struct {
	ID          int    `json:"USER_ID"`
	UserPicture string `json:"USER_PICTURE"`
}

// PrivateError is a mapping of error messages.
type PrivateError struct {
	GatewayError string `json:"GATEWAY_ERROR"`
}
