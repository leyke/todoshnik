package api

import (
	"net/http"
	"strconv"

	"todoshnik/internal/infrastructure/identity"

	"github.com/go-chi/chi/v5"
)

func getTaskID(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}

func getScope(userID int) identity.AccessScope {
	return identity.AccessScope{
		IsAdmin: false,
		UserID:  userID,
	}
}
