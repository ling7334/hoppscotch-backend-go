package strategies

import (
	"encoding/json"
	ex "exception"
	"fmt"
	mw "middleware"
	"model"
	"net/http"
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

func authCookieHandler(w http.ResponseWriter, r *http.Request, token *authTokens) {
	now := time.Now()
	aexp := now.Add(time.Duration(accessExpires))
	rexp := now.Add(time.Duration(refreshExpires))
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.Access_token,
		Expires:  aexp,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token.Refresh_token,
		Expires:  rexp,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

func createSSOUser(db *gorm.DB, profile profile) (*model.User, error) {
	user := &model.User{
		UID:         cuid.New(),
		Email:       &profile.Email[0],
		DisplayName: &profile.DisplayName[0],
		PhotoURL:    &profile.Photos[0],
	}
	err := db.Create(user).Error
	return user, err
}

func createSSOAccount(db *gorm.DB, profile profile, uid string) (*model.Account, error) {
	account := &model.Account{
		ID:                   cuid.New(),
		Provider:             profile.Provider,
		ProviderAccountID:    profile.ID,
		ProviderAccessToken:  &profile.AccessToken,
		ProviderRefreshToken: &profile.RefreshToken,
		UserID:               uid,
	}
	err := db.Create(account).Error
	return account, err
}

func responder(w http.ResponseWriter, r any, status int) {
	if res, err := json.Marshal(r); err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%s","error":"%s","statusCode":%d}`, err, ex.ErrJSONInvalid, http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(status)
		w.Write(res)
		return
	}
}

func getDB(w http.ResponseWriter, r *http.Request) *gorm.DB {
	ctx := r.Context()
	db, ok := ctx.Value(mw.ContextKey("DB")).(*gorm.DB)
	if !ok {
		responder(w, resp{
			"database not exist",
			ex.ErrBugAuthNoUserCtx.Error(),
			http.StatusInternalServerError,
		}, http.StatusInternalServerError)
		return nil
	}
	return db
}

func SignOrLogin(tx *gorm.DB, w http.ResponseWriter, r *http.Request, profile profile, RedirectURL string) {
	user := &model.User{}
	err := tx.Where("email =?", profile.Email[0]).First(user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		user, err = createSSOUser(tx, profile)
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if user.DisplayName == nil {
		user.DisplayName = &profile.DisplayName[0]
		err = tx.Save(user).Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to save user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	at, err := mw.NewToken(user.UID, true)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to generate access token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rt, err := mw.NewToken(user.UID, false)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to generate refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user.RefreshToken = &rt
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to set refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	account := &model.Account{}
	err = tx.Where(`provider=? AND "providerAccountId"=?`, "github", profile.ID).First(account).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			http.Error(w, "Failed to get account: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = createSSOAccount(tx, profile, user.UID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to create account: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	authCookieHandler(w, r, &authTokens{at, rt})
	http.Redirect(w, r, RedirectURL, http.StatusOK)
}
