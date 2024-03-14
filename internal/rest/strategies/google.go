package strategies

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func getUserDataFromGoogle(config *oauth2.Config, ctx context.Context, token *oauth2.Token) (GoogleUserInfo, error) {
	response, err := config.Client(ctx, token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return GoogleUserInfo{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return GoogleUserInfo{}, fmt.Errorf("failed read response: %s", err.Error())
	}
	var userInfo GoogleUserInfo
	err = json.Unmarshal(contents, &userInfo)
	return userInfo, err
}

func GoogleCallback(config *oauth2.Config) http.Handler {
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
		if r.FormValue("state") != state.Value {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userInfo, err := getUserDataFromGoogle(config, r.Context(), token)
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		profile := profile{
			ID:           userInfo.ID,
			Provider:     "google",
			Email:        []string{userInfo.Email},
			DisplayName:  []string{userInfo.Name},
			Photos:       []string{userInfo.Picture},
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		SignOrLogin(tx, w, r, profile, config.RedirectURL)
	})
}
