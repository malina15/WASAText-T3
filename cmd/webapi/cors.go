package main

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         string
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		MaxAge:         "1",
	}
}

// CORSHandler creates a CORS middleware handler
func CORSHandler(config *CORSConfig) func(http.HandlerFunc) httprouter.Handle {
	return func(handler http.HandlerFunc) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Set CORS headers
			origin := r.Header.Get("Origin")
			if origin == "" || isOriginAllowed(origin, config.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(config.AllowedOrigins) > 0 && config.AllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", config.MaxAge)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Call the actual handler
			handler.ServeHTTP(w, r)
		}
	}
}

// isOriginAllowed checks if an origin is allowed
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

// SimpleCORSHandler creates a simple CORS handler for development
func SimpleCORSHandler(handler http.HandlerFunc) httprouter.Handle {
	config := DefaultCORSConfig()
	return CORSHandler(config)(handler)
}
