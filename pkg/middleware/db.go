package middleware

import (
	"context"
	"net/http"

	"model"

	"gorm.io/gorm"
)

// DBMiddleware add gorm DB instance to the context
func DBMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			next.ServeHTTP(w, r)
		}()
		ctx := context.WithValue(r.Context(), ContextKey("DB"), db)
		r = r.WithContext(ctx)
	})
}

func OperatorMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			next.ServeHTTP(w, r)
		}()
		sess, ok := r.Context().Value(ContextKey("session")).(Session)
		if !ok {
			return
		}
		user := &model.User{}
		if db.First(user, "uid=?", sess.Uid).Error != nil {
			return
		}
		ctx := context.WithValue(r.Context(), ContextKey("operator"), user)
		r = r.WithContext(ctx)
	})
}
