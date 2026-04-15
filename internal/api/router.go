package api

import (
	"net/http"
	"todoshnik/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func (api *APIHandler) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logging(api.logger))

	// Task
	r.Route("/task", func(r chi.Router) {
		r.Get("/", api.taskHandler.List)
		r.Post("/", api.taskHandler.Create)
	})

	r.Get("/ping", api.pingHandler)

	return r
}
