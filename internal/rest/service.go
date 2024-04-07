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

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func searchByTitle(db *gorm.DB, teamID, term string) (res []searchResult, err error) {
	coll, err := searchCollections(db, teamID, term)
	if err != nil {
		return nil, err
	}
	res = append(res, coll...)
	req, err := searchRequests(db, teamID, term)
	if err != nil {
		return nil, err
	}
	res = append(res, req...)
	return
}

func searchCollections(db *gorm.DB, teamID, term string) (res []searchResult, err error) {
	coll := []model.TeamCollection{}
	if err := db.Where(`"teamID"=? AND title @@ ?`, teamID, term).Find(&coll).Error; err != nil {
		return nil, err
	}
	for _, c := range coll {
		parent := []searchResult{}
		if c.ParentID != nil {
			if parent, err = findParentCollection(db, *c.ParentID); err != nil {
				return nil, err
			}
		}
		res = append(res, searchResult{
			ID:    c.ID,
			Title: c.Title,
			Type:  "collection",
			Path:  parent,
		})
	}
	return
}

func searchRequests(db *gorm.DB, teamID, term string) (res []searchResult, err error) {
	req := []model.TeamRequest{}
	if err := db.Where(`"teamID"=? AND title @@ ?`, teamID, term).Find(&req).Error; err != nil {
		return nil, err
	}
	for _, c := range req {
		if parent, err := findParentCollection(db, c.CollectionID); err != nil {
			return nil, err
		} else {
			res = append(res, searchResult{
				ID:     c.ID,
				Title:  c.Title,
				Type:   "request",
				Method: c.Request.Method,
				Path:   parent,
			})
		}
	}
	return
}

func findParentCollection(db *gorm.DB, id string) (res []searchResult, err error) {
	coll := &model.TeamCollection{}
	if err := db.Where("id=?", id).First(coll).Error; err != nil {
		return nil, err
	} else {
		parent := []searchResult{}
		if coll.ParentID != nil {
			if parent, err = findParentCollection(db, *coll.ParentID); err != nil {
				return nil, err
			}
		}
		res = append(res, searchResult{
			ID:    coll.ID,
			Title: coll.Title,
			Type:  "collection",
			Path:  parent,
		})
	}
	return
}
