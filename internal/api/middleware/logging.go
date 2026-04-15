package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func Logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println(2, r.URL.Path)
			next.ServeHTTP(w, r)
			fmt.Println("Конец обработки")
		})
	}
}
