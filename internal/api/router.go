package api

import (
	"net/http"
	"todoshnik/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func (api *APIHandler) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logging(api.logger))

	// Tasks
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", api.taskHandler.List)
		r.Post("/", api.taskHandler.Create)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", api.taskHandler.View)      // GET //tasks/123
			r.Put("/", api.taskHandler.Update)    // PUT /tasks/123
			r.Post("/done", api.taskHandler.Done) // POST /tasks/123/done
			r.Delete("/", api.taskHandler.Delete) // DELETE /tasks/123
		})
	})

	r.Get("/ping", api.pingHandler)

	return r
}
