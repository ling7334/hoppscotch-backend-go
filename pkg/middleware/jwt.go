package middleware

import (
	"context"
	"exception"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

const defaultSecret = "secret123"
const defaultAccessExpires = 24 * time.Hour
const defaultRefreshExpires = 7 * 24 * time.Hour

var secret string
var baseURL string
var accessExpires time.Duration
var refreshExpires time.Duration

func init() {
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env")
	}
	baseURL = os.Getenv("VITE_BASE_URL")
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
}

type Session struct {
	Uid string
	Sid string
}

func NewToken(uid string, access bool) (string, error) {
	now := time.Now().UnixMilli()
	var exp int64
	if access {
		exp = accessExpires.Milliseconds() + now
	} else {
		exp = refreshExpires.Milliseconds() + now
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": baseURL,
		"sub": uid,
		"aud": []string{baseURL},
		"iat": now,
		"exp": exp,
	}).SignedString([]byte(secret))
}

// JwtMiddleware decodes the share session cookie and packs the logged user into context
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			next.ServeHTTP(w, r)
		}()
		at, err := r.Cookie("access_token")
		if err != nil || at == nil {
			slog.Error("access_token not found", "error", err)
			return
		}
		rt, err := r.Cookie("refresh_token")
		if err != nil || rt == nil {
			slog.Error("refresh_token not found", "error", err)
			return
		}
		var sidStr string
		sid, err := r.Cookie("connect.sid")
		if err == nil && sid != nil {
			sidStr = sid.Value
		}
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// validate token
		access, err := jwt.Parse(at.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			if err == jwt.ErrTokenExpired {
				// access token expired, check Refersh token
				_, err := jwt.Parse(rt.Value, func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					return []byte(secret), nil
				})
				if err != nil {
					switch err {
					case jwt.ErrTokenExpired:
						slog.Error("refresh token expired", "error", err)
						ctx := context.WithValue(ctx, ContextKey("error"), exception.ErrTokenExpired)
						r = r.WithContext(ctx)
						return
					default:
						slog.Error("refresh token parse failed", "error", err)
						ctx := context.WithValue(ctx, ContextKey("error"), exception.ErrInvalidRefreshToken)
						r = r.WithContext(ctx)
						return
					}
				}
				if claims, ok := access.Claims.(jwt.MapClaims); ok {
					// all good
					uid, err := claims.GetSubject()
					if err != nil {
						slog.Error("get uid failed", "error", err)
						ctx = context.WithValue(ctx, ContextKey("error"), err)
						r = r.WithContext(ctx)
						return
					}
					// refresh token is valid, get new access token
					token, err := NewToken(uid, true)
					if err != nil {
						slog.Error("refresh generate token failed", "error", err)
						ctx = context.WithValue(ctx, ContextKey("error"), err)
						r = r.WithContext(ctx)
						return
					}
					access_cookie := http.Cookie{
						Name:     "access_token",
						Value:    token,
						Expires:  time.Now().Add(accessExpires),
						HttpOnly: true,
						Secure:   true,
						SameSite: http.SameSiteLaxMode,
					}
					http.SetCookie(w, &access_cookie)
				}
			}
		}
		if claims, ok := access.Claims.(jwt.MapClaims); ok {
			// all good
			uid, err := claims.GetSubject()
			if err != nil {
				slog.Error("get uid failed", "error", err)
				ctx = context.WithValue(ctx, ContextKey("error"), err)
			} else {
				slog.Debug("user", "uid", uid)
				ctx = context.WithValue(ctx, ContextKey("session"), Session{
					Uid: uid,
					Sid: sidStr,
				})
			}
			r = r.WithContext(ctx)
			return
		}
		// access token parse failed
		slog.Error("invaild access token", "error", err)
		ctx = context.WithValue(ctx, ContextKey("error"), exception.ErrInvalidAccessToken)
		r = r.WithContext(ctx)
	})
}
