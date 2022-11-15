package middleware

import (
	"log"
	"net/http"
)

func MiddlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Print("Executing middlewareOne again")
	})
}

func MiddlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareTwo")
		if r.URL.Path == "/foo" {
			return
		}

		next.ServeHTTP(w, r)
		log.Print("Executing middlewareTwo again")
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type Middleware func(http.Handler) http.Handler

func ChainingMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) < 1 {
		return h
	}

	wrappedHandler := h
	for i := len(m) - 1; i >= 0; i-- {
		wrappedHandler = m[i](wrappedHandler)
	}

	return wrappedHandler
}
