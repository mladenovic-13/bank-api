package middlewares

import (
	"net/http"

	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *MiddlewareContext) WithAuth(handler api.ProtectedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie("token")

		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims, err := utils.ValidateJWT(jwtCookie.Value)

		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := ctx.Queries.GetUserByIDAndUsername(
			r.Context(),
			database.GetUserByIDAndUsernameParams{
				ID:       claims.ID,
				Username: claims.Username,
			},
		)

		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		handler(w, r, *models.ToUser(user))
	}
}
