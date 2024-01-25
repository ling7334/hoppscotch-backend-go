package rest

type signInMagicDto struct {
	Email string `json:"email"`
}

type verifyMagicDto struct {
	DeviceIdentifier string `json:"deviceIdentifier"`
	Token            string `json:"token"`
}

type authTokens struct {
	Access_token  string
	Refresh_token string
}

type provider struct {
	Providers []string `json:"providers"`
}

type profile struct {
	ID          string   `json:"id"`
	Provider    string   `json:"provider"`
	Email       []string `json:"email"`
	DisplayName []string `json:"displayName"`
	Photos      []string `json:"photos"`
}

type resp struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}
