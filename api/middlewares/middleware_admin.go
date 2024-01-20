package middlewares

import (
	"net/http"

	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/models"
)

func (ctx *MiddlewareContext) WithAdmin(next http.Handler) http.Handler {
	return ctx.WithAuth(func(w http.ResponseWriter, r *http.Request, u models.User) {
		if !u.IsAdmin {
			api.RespondWithError(w, http.StatusUnauthorized, "Only admin can perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}
