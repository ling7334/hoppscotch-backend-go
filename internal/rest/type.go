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

type resp struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}
