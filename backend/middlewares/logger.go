package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		time := time.Now()
		url := r.URL
		method := r.Method

		log.Printf("%v, %v, %v, %v", time, method, header, url)

		next.ServeHTTP(w, r)
	})
}
