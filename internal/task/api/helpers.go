package api

import (
	"errors"
	"net/http"
	"strconv"
	"todoshnik/internal/api/contextkeys"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/identity"

	"github.com/go-chi/chi/v5"
)

func getTaskID(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}

func getUserID(r *http.Request) int {
	return r.Context().Value(contextkeys.UserIDKey).(int)
}

func getScope(r *http.Request) identity.AccessScope {
	return identity.AccessScope{
		IsAdmin: false,
		UserID:  getUserID(r),
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
