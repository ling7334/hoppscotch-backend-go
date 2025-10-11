package strategies

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
)

type GithubPlan struct {
	Name          string `json:"name"`
	Space         int64  `json:"space"`
	Collaborators int64  `json:"collaborators"`
	PrivateRepos  int64  `json:"private_repos"`
}
type GithubUserInfo struct {
	Login                   string     `json:"login"`
	ID                      int64      `json:"id"`
	NodeID                  string     `json:"node_id"`
	AvatarURL               string     `json:"avatar_url"`
	GravatarID              string     `json:"gravatar_id"`
	URL                     string     `json:"url"`
	HTMLURL                 string     `json:"html_url"`
	FollowersURL            string     `json:"followers_url"`
	FollowingURL            string     `json:"following_url"`
	GistsURL                string     `json:"gists_url"`
	StarredURL              string     `json:"starred_url"`
	SubscriptionsURL        string     `json:"subscriptions_url"`
	OrganizationsURL        string     `json:"organizations_url"`
	ReposURL                string     `json:"repos_url"`
	EventsURL               string     `json:"events_url"`
	ReceivedEventsURL       string     `json:"received_events_url"`
	Type                    string     `json:"type"`
	SiteAdmin               bool       `json:"site_admin"`
	Name                    string     `json:"name"`
	Company                 string     `json:"company"`
	Blog                    string     `json:"blog"`
	Location                string     `json:"location"`
	Email                   string     `json:"email"`
	Hireable                bool       `json:"hireable"`
	Bio                     string     `json:"bio"`
	TwitterUsername         string     `json:"twitter_username"`
	PublicRepos             int64      `json:"public_repos"`
	PublicGists             int64      `json:"public_gists"`
	Followers               int64      `json:"followers"`
	Following               int64      `json:"following"`
	CreatedAt               string     `json:"created_at"`
	UpdatedAt               string     `json:"updated_at"`
	PrivateGists            int64      `json:"private_gists"`
	TotalPrivateRepos       int64      `json:"total_private_repos"`
	OwnedPrivateRepos       int64      `json:"owned_private_repos"`
	DiskUsage               int64      `json:"disk_usage"`
	Collaborators           int64      `json:"collaborators"`
	TwoFactorAuthentication bool       `json:"two_factor_authentication"`
	Plan                    GithubPlan `json:"plan"`
}

type GithubEmails struct {
	Email      string  `json:"email"`
	Primary    bool    `json:"primary"`
	Verified   bool    `json:"verified"`
	Visibility *string `json:"visibility"`
}

func getGithubUserInfo(token string) (GithubUserInfo, error) {
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return GithubUserInfo{}, err
	}

	req.Header.Set("Authorization", "token "+token)
	res, err := httpClient.Do(req)
	if err != nil {
		return GithubUserInfo{}, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return GithubUserInfo{}, err
	}
	slog.Info("Github userinfo", "userinfo", bytes)
	var userInfo GithubUserInfo
	err = json.Unmarshal(bytes, &userInfo)
	return userInfo, err
}

func getGithubEmails(token string) ([]GithubEmails, error) {
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return []GithubEmails{}, err
	}

	req.Header.Set("Authorization", "token "+token)
	res, err := httpClient.Do(req)
	if err != nil {
		return []GithubEmails{}, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []GithubEmails{}, err
	}
	slog.Info("Github userinfo", "userinfo", bytes)
	var emails []GithubEmails
	err = json.Unmarshal(bytes, &emails)
	return emails, err
}

func GithubCallback(config *oauth2.Config) http.Handler {
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

		oauth2Token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		userInfo, err := getGithubUserInfo(oauth2Token.AccessToken)
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var email_addrs []string
		if emails, err := getGithubEmails(oauth2Token.AccessToken); err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			for _, email := range emails {
				email_addrs = append(email_addrs, email.Email)
			}
		}
		profile := profile{
			ID:           strconv.FormatInt(userInfo.ID, 10),
			Provider:     "github",
			Email:        email_addrs,
			DisplayName:  []string{userInfo.Name},
			Photos:       []string{userInfo.AvatarURL},
			AccessToken:  oauth2Token.AccessToken,
			RefreshToken: oauth2Token.RefreshToken,
		}

		SignOrLogin(tx, w, r, profile)

	})
}
