package rest

import (
	"encoding/json"
	ex "exception"
	"fmt"
	mw "middleware"
	"net/http"

	"gorm.io/gorm"
)

func responder(w http.ResponseWriter, r any, status int) {
	if res, err := json.Marshal(r); err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%s","error":"%s","statusCode":%d}`, err, ex.ErrJSONInvalid, http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
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
