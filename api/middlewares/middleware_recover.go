package middlewares

import (
	"net/http"

	"github.com/mladenovic-13/bank-api/api"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
