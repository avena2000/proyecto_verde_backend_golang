package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Registra la petición entrante
		log.Printf("Iniciando petición %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		// Registra el tiempo de respuesta
		log.Printf("Completada petición %s %s en %v", r.Method, r.URL.Path, time.Since(start))
	})
}
