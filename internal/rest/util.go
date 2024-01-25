package rest

import (
	"encoding/json"
	ex "exception"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func responder(w http.ResponseWriter, r any, status int) {
	if res, err := json.Marshal(r); err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"%s","error":"%s","statusCode":%d}`, err, ex.ErrJSONInvalid, http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		w.Write(res)
		w.WriteHeader(status)
		return
	}
}

func getDB(w http.ResponseWriter, r *http.Request) *gorm.DB {
	ctx := r.Context()
	db, ok := ctx.Value("DB").(*gorm.DB)
	if !ok {
		responder(w, resp{
			"database not exist",
			ex.ErrEnvEmptyAuthProviders.Error(),
			http.StatusInternalServerError,
		}, http.StatusInternalServerError)
		return nil
	}
	return db
}
