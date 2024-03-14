package strategies

type authTokens struct {
	Access_token  string
	Refresh_token string
}

type resp struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

type profile struct {
	ID           string   `json:"id"`
	Provider     string   `json:"provider"`
	Email        []string `json:"email"`
	DisplayName  []string `json:"displayName"`
	Photos       []string `json:"photos"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}
