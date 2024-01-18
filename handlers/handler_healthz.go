package handlers

import (
	"net/http"
)

func (ctx *RouterCtx) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "ok")
}
