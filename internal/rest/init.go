package rest

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const defaultSecret = "secret123"
const defaultAccessExpires = 24 * time.Hour
const defaultRefreshExpires = 7 * 24 * time.Hour

const GoogleURL = "https://accounts.google.com"
const GithubURL = "https://github.com"
const MicrosoftURL = "https://login.microsoftonline.com/%s"

var secret string
var accessExpires time.Duration
var refreshExpires time.Duration
var GoogleProvider *oidc.Provider
var GithubProvider *oidc.Provider
var MicrosoftProvider *oidc.Provider
var GoogleConfig *oauth2.Config
var GithubConfig *oauth2.Config
var MicrosoftConfig *oauth2.Config

func init() {
	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = defaultSecret
	}
	sec, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_VALIDITY"), 10, 0)
	if err != nil {
		accessExpires = defaultAccessExpires
	} else {
		accessExpires = time.Duration(sec) * time.Millisecond
	}
	sec, err = strconv.ParseInt(os.Getenv("REFRESH_TOKEN_VALIDITY"), 10, 0)
	if err != nil {
		refreshExpires = defaultRefreshExpires
	} else {
		refreshExpires = time.Duration(sec) * time.Millisecond
	}
	google_client_id := os.Getenv("GOOGLE_CLIENT_ID")
	google_client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	google_callback_url := os.Getenv("GOOGLE_CALLBACK_URL")
	google_scope := os.Getenv("GOOGLE_SCOPE")
	if google_client_id != "" && google_client_secret != "" && google_callback_url != "" && google_scope != "" {
		GoogleProvider, err = oidc.NewProvider(context.Background(), GoogleURL)
		if err == nil {
			scope := strings.Split(google_scope, ",")
			scope = append(scope, oidc.ScopeOpenID)
			GoogleConfig = &oauth2.Config{
				ClientID:     google_client_id,
				ClientSecret: google_client_secret,
				Endpoint:     GoogleProvider.Endpoint(),
				RedirectURL:  google_callback_url,
				Scopes:       scope,
			}
		}
	}
	github_client_id := os.Getenv("GITHUB_CLIENT_ID")
	github_client_secret := os.Getenv("GITHUB_CLIENT_SECRET")
	github_callback_url := os.Getenv("GITHUB_CALLBACK_URL")
	github_scope := os.Getenv("GITHUB_SCOPE")
	if github_client_id != "" && github_client_secret != "" && github_callback_url != "" && github_scope != "" {
		GithubProvider, err := oidc.NewProvider(context.Background(), GithubURL)
		if err == nil {
			scope := strings.Split(github_scope, ",")
			scope = append(scope, oidc.ScopeOpenID)
			GithubConfig = &oauth2.Config{
				ClientID:     github_client_id,
				ClientSecret: github_client_secret,
				Endpoint:     GithubProvider.Endpoint(),
				RedirectURL:  github_callback_url,
				Scopes:       scope,
			}
		}
	}
	microsoft_client_id := os.Getenv("MICROSOFT_CLIENT_ID")
	microsoft_client_secret := os.Getenv("MICROSOFT_CLIENT_SECRET")
	microsoft_callback_url := os.Getenv("MICROSOFT_CALLBACK_URL")
	microsoft_scope := os.Getenv("MICROSOFT_SCOPE")
	microsoft_tenant := os.Getenv("MICROSOFT_TENANT")
	if microsoft_client_id != "" && microsoft_client_secret != "" && microsoft_callback_url != "" && microsoft_scope != "" {
		MicrosoftProvider, err := oidc.NewProvider(context.Background(), fmt.Sprintf(MicrosoftURL, microsoft_tenant))
		if err == nil {
			scope := strings.Split(microsoft_scope, ",")
			scope = append(scope, oidc.ScopeOpenID)
			MicrosoftConfig = &oauth2.Config{
				ClientID:     microsoft_client_id,
				ClientSecret: microsoft_client_secret,
				Endpoint:     MicrosoftProvider.Endpoint(),
				RedirectURL:  microsoft_callback_url,
				Scopes:       scope,
			}
		}
	}
}
