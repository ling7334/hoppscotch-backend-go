package strategies

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type MicrosoftUserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
}

func getMicrosoftUserInfo(config *oauth2.Config, ctx context.Context, token *oauth2.Token) (MicrosoftUserInfo, error) {

	response, err := config.Client(ctx, token).Get("https://graph.microsoft.com/oidc/userinfo")
	if err != nil {
		return MicrosoftUserInfo{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return MicrosoftUserInfo{}, fmt.Errorf("failed read response: %s", err.Error())
	}
	var userInfo MicrosoftUserInfo
	err = json.Unmarshal(contents, &userInfo)
	return userInfo, err
}

func MicrosoftCallback(config *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := getDB(w, r)
		if db == nil {
			return
		}
		tx := db.Begin()
		state, err := r.Cookie("state")
		if err != nil {
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("state") != state.Value {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		userInfo, err := getMicrosoftUserInfo(config, r.Context(), token)
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		profile := profile{
			ID:           userInfo.Sub,
			Provider:     "microsoft",
			Email:        []string{userInfo.Email},
			DisplayName:  []string{userInfo.Name},
			Photos:       []string{userInfo.Picture},
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		SignOrLogin(tx, w, r, profile)
	})
}
