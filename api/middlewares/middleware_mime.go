package middlewares

import (
	"log"
	"net/http"
	"path/filepath"
)

func WithMime(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("URL:", r.URL.Path)

			contentType := getContentType(r.URL.Path)

			log.Println("MIME Type:", contentType)

			w.Header().Set("Content-Type", contentType)

			handler.ServeHTTP(w, r)
		})
}

func getContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "text/javascript"
	case ".json":
		return "application/json"
	default:
		return "text/plain"
	}
}
