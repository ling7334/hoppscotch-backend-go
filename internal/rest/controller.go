package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	ex "exception"
	"mail"
	mw "middleware"
	"model"

	"github.com/coreos/go-oidc/v3/oidc"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func ServeMux(path string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(path+"providers", Providers)
	mux.HandleFunc(path+"signin", SignInMagicLink)
	mux.HandleFunc(path+"verify", Verify)
	mux.HandleFunc(path+"refresh", Refersh)
	mux.HandleFunc(path+"logout", Logout)
	mux.Handle(path+"google", Redirect(GoogleConfig))
	mux.Handle(path+"google/callback", Callback(GoogleConfig, GoogleProvider))
	mux.Handle(path+"github", Redirect(GithubConfig))
	mux.Handle(path+"github/callback", Callback(GithubConfig, GithubProvider))
	mux.Handle(path+"microsoft", Redirect(MicrosoftConfig))
	mux.Handle(path+"microsoft/callback", Callback(MicrosoftConfig, MicrosoftProvider))

	return mux
}

// Ping is healthcheck endpoint
func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func Providers(w http.ResponseWriter, r *http.Request) {
	p := provider{
		strings.Split(os.Getenv("VITE_ALLOWED_AUTH_PROVIDERS"), ","),
	}
	if len(p.Providers) == 0 {
		responder(w, resp{
			ex.ErrEnvEmptyAuthProviders.Error(),
			ex.ErrEnvEmptyAuthProviders.Error(),
			http.StatusInternalServerError,
		}, http.StatusInternalServerError)
		return
	}
	responder(w, p, http.StatusOK)
}

// Route to initiate magic-link auth for a users email
func SignInMagicLink(w http.ResponseWriter, r *http.Request) {
	db := getDB(w, r)
	if db == nil {
		return
	}
	authData := signInMagicDto{}
	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		responder(w, resp{
			err.Error(),
			ex.ErrJSONInvalid.Error(),
			http.StatusBadRequest,
		}, http.StatusBadRequest)
		return
	}
	user := &model.User{}
	token := &model.VerificationToken{}
	tx := db.Begin()

	if err := tx.First(user, "email=?", authData.Email).Error; err == nil {
		var err error
		if token, err = generateMagicLinkTokens(user, tx); err != nil {
			tx.Rollback()
			responder(w, resp{
				err.Error(),
				ex.ErrJSONInvalid.Error(),
				http.StatusInternalServerError,
			}, http.StatusInternalServerError)
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		if user, err = createUserViaMagicLink(authData.Email, tx); err != nil {
			tx.Rollback()
			responder(w, resp{
				err.Error(),
				ex.ErrJSONInvalid.Error(),
				http.StatusInternalServerError,
			}, http.StatusInternalServerError)
			return
		}
		if token, err = generateMagicLinkTokens(user, tx); err != nil {
			tx.Rollback()
			responder(w, resp{
				err.Error(),
				ex.ErrJSONInvalid.Error(),
				http.StatusInternalServerError,
			}, http.StatusInternalServerError)
			return
		}
	} else {
		responder(w, resp{
			err.Error(),
			ex.ErrJSONInvalid.Error(),
			http.StatusInternalServerError,
		}, http.StatusInternalServerError)
		return
	}
	var url string
	query := r.URL.Query()
	origin := query.Get("origin")
	switch origin {
	case "admin":
		url = os.Getenv("VITE_ADMIN_URL")
	case "app":
		url = os.Getenv("VITE_BASE_URL")
	default:
		// if origin is invalid by default set URL to Hoppscotch-App
		url = os.Getenv("VITE_BASE_URL")
	}
	magicLink := fmt.Sprintf("%s/enter?token=%s", url, token.Token)
	if err := mail.SendUserInvitation(*user.Email, magicLink); err != nil {
		tx.Rollback()
		responder(w, resp{
			err.Error(),
			err.Error(),
			http.StatusInternalServerError,
		}, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"deviceIdentifier\":\"%s\"}", token.DeviceIdentifier)))
}

// Route to verify and sign in a valid user via magic-link
func Verify(w http.ResponseWriter, r *http.Request) {
	db := getDB(w, r)
	if db == nil {
		return
	}
	authData := verifyMagicDto{}
	err := json.NewDecoder(r.Body).Decode(&authData)
	if err != nil {
		responder(w, resp{
			err.Error(),
			ex.ErrJSONInvalid.Error(),
			http.StatusBadRequest,
		}, http.StatusBadRequest)
		return
	}
	authTokens, err := verifyMagicLinkTokens(&authData, db)
	if err != nil {
		switch err {
		case ex.ErrInvalidMagicLinkData:
			responder(w, resp{
				err.Error(),
				err.Error(),
				http.StatusNotFound,
			}, http.StatusNotFound)
		default:
			responder(w, resp{
				err.Error(),
				err.Error(),
				http.StatusBadRequest,
			}, http.StatusBadRequest)

		}
		return
	}
	log.Info().Msgf("ccess_token: %s \n refresh_token: %s", authTokens.Access_token, authTokens.Refresh_token)
	authCookieHandler(w, r, authTokens)
	w.WriteHeader(http.StatusOK)
}

/*
Route to refresh auth tokens with Refresh Token Rotation

@see https://auth0.com/docs/secure/tokens/refresh-tokens/refresh-token-rotation
*/
func Refersh(w http.ResponseWriter, r *http.Request) {
	rt, err := r.Cookie("refresh_token")
	if err != nil || rt == nil {
		responder(w, resp{
			ex.ErrCookiesNotFound.Error(),
			ex.ErrCookiesNotFound.Error(),
			http.StatusNotFound,
		}, http.StatusNotFound)
		return
	}
	rToken, err := jwt.Parse(rt.Value, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		responder(w, resp{
			ex.ErrInvalidRefreshToken.Error(),
			ex.ErrInvalidRefreshToken.Error(),
			http.StatusBadRequest,
		}, http.StatusBadRequest)
		return
	}
	if claims, ok := rToken.Claims.(jwt.MapClaims); ok {
		uid, err := claims.GetSubject()
		if err != nil {
			responder(w, resp{
				ex.ErrInvalidRefreshToken.Error(),
				ex.ErrInvalidRefreshToken.Error(),
				http.StatusBadRequest,
			}, http.StatusBadRequest)
			return
		}
		tk, err := mw.NewToken(uid, false)
		if err != nil {
			responder(w, resp{
				err.Error(),
				ex.ErrInvalidRefreshToken.Error(),
				http.StatusInternalServerError,
			}, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tk,
			Expires:  time.Now().Add(refreshExpires),
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusOK)
	}

}

// Log user out by clearing cookies containing auth tokens
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: "",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: "",
	})
	w.WriteHeader(http.StatusOK)
}

func Redirect(config *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		state, err := randString(16)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		c := &http.Cookie{
			Name:     "state",
			Value:    state,
			MaxAge:   int(time.Hour.Seconds()),
			Secure:   r.TLS != nil,
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
	})
}

func Callback(config *oauth2.Config, provider *oidc.Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := getDB(w, r)
		if db == nil {
			return
		}
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

		userInfo, err := provider.UserInfo(r.Context(), oauth2.StaticTokenSource(oauth2Token))
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var profile profile
		if err := json.Unmarshal([]byte(userInfo.Profile), &profile); err != nil {
			http.Error(w, "Failed to parse profile: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user := &model.User{}
		err = db.Where("email =?", userInfo.Email).First(user).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
				return
			}
			user, err = createSSOUser(db, profile)
			if err != gorm.ErrRecordNotFound {
				http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if user.DisplayName == nil {
			user.DisplayName = &userInfo.Email
			err = db.Save(user).Error
			if err != nil {
				http.Error(w, "Failed to save user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		at, err := mw.NewToken(user.UID, true)
		if err != nil {
			http.Error(w, "Failed to generate access token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rt, err := mw.NewToken(user.UID, false)
		if err != nil {
			http.Error(w, "Failed to generate refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		user.RefreshToken = &rt
		err = db.Save(&user).Error
		if err != nil {
			http.Error(w, "Failed to set refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		account := &model.Account{}
		err = db.Where("provider=? AND providerAccountId=?", "magic", user.Email).First(account).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				http.Error(w, "Failed to get account: "+err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = createSSOAccount(db, profile, user.UID, at, rt)
			if err != nil {
				http.Error(w, "Failed to create account: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		authCookieHandler(w, r, &authTokens{at, rt})
		http.Redirect(w, r, config.RedirectURL, http.StatusOK)
	})
}
