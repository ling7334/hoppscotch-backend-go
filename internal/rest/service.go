package rest

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	ex "exception"
	mw "middleware"
	"model"

	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func generateMagicLinkTokens(user *model.User, db *gorm.DB) (*model.VerificationToken, error) {
	saltComplexityStr := os.Getenv("TOKEN_SALT_COMPLEXITY")
	tokenValidityHoursStr := os.Getenv("MAGIC_LINK_TOKEN_VALIDITY")
	tokenValidityHours, err := strconv.Atoi(tokenValidityHoursStr)
	if err != nil {
		// Handle the error
		return nil, err
	}
	// Generate a salt using bcrypt
	salt, err := bcrypt.GenerateFromPassword([]byte(saltComplexityStr), bcrypt.DefaultCost)
	if err != nil {
		// Handle the error
		return nil, err
	}

	// Calculate the expiration time
	expiresOn := time.Now().Add(time.Hour * time.Duration(tokenValidityHours))

	idToken := &model.VerificationToken{
		DeviceIdentifier: string(salt),
		Token:            cuid.New(),
		UserUID:          user.UID,
		ExpiresOn:        expiresOn,
	}
	err = db.Create(idToken).Error
	return idToken, err
}

func createUserViaMagicLink(email string, db *gorm.DB) (*model.User, error) {
	user := &model.User{
		UID:         cuid.New(),
		DisplayName: &email,
		Email:       &email,
	}
	err := db.Create(user).Error
	return user, err
}

func createMagicAccount(db *gorm.DB, user *model.User) (*model.Account, error) {
	account := &model.Account{
		ID:                cuid.New(),
		UserID:            user.UID,
		Provider:          "magic",
		ProviderAccountID: *user.Email,
	}
	err := db.Create(account).Error
	return account, err
}

func verifyMagicLinkTokens(dto *verifyMagicDto, db *gorm.DB) (*authTokens, error) {
	tk := &model.VerificationToken{}
	err := db.Preload("User").Where("\"deviceIdentifier\"=?", dto.DeviceIdentifier).Where("token=?", dto.Token).First(tk).Error
	if err != nil {
		return nil, ex.ErrInvalidMagicLinkData
		// switch err {
		// case gorm.ErrRecordNotFound:
		// 	return nil, ex.ErrInvalidMagicLinkData
		// default:
		// 	return nil, err
		// }
	}
	if tk.ExpiresOn.Before(time.Now()) {
		return nil, ex.ErrTokenExpired
	}
	user := tk.User
	account := &model.Account{}
	err = db.Where("provider=? AND \"providerAccountId\"=?", "magic", user.Email).First(account).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, ex.ErrInvalidMagicLinkData
		}
		_, err = createMagicAccount(db, &user)
		if err != nil {
			return nil, ex.ErrInvalidMagicLinkData
		}
	}

	at, err := mw.NewToken(user.UID, true)
	if err != nil {
		return nil, ex.ErrInvalidMagicLinkData
	}
	rt, err := mw.NewToken(user.UID, false)
	if err != nil {
		return nil, ex.ErrInvalidMagicLinkData
	}
	user.RefreshToken = &rt
	err = db.Save(&user).Error
	if err != nil {
		return nil, ex.ErrInvalidMagicLinkData
	}
	// delete token
	err = db.Delete(tk, "\"deviceIdentifier\"=? AND token=?", tk.DeviceIdentifier, tk.Token).Error
	if err != nil {
		return nil, ex.ErrInvalidMagicLinkData
	}
	return &authTokens{at, rt}, nil
}

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

func createSSOAccount(db *gorm.DB, profile profile, uid, access, refresh string) (*model.Account, error) {
	account := &model.Account{
		ID:                   cuid.New(),
		Provider:             profile.Provider,
		ProviderAccountID:    profile.ID,
		ProviderAccessToken:  &access,
		ProviderRefreshToken: &refresh,
		UserID:               uid,
	}
	err := db.Create(account).Error
	return account, err
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
