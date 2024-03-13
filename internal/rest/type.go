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

// // UserInfo represents the OpenID Connect userinfo claims.
// type UserInfo struct {
// 	Subject       string `json:"sub"`
// 	Profile       string `json:"profile"`
// 	Email         string `json:"email"`
// 	EmailVerified bool   `json:"email_verified"`

// 	claims []byte
// }

type UserInfo struct {
	AvatarURL     string `json:"avatar_url"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
