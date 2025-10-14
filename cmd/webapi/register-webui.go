//go:build webui

package main

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

//go:embed ../../webui/dist
var webuiFiles embed.FS

func registerWebUI() {
	// Create a filesystem from the embedded webui files
	webuiFS, err := fs.Sub(webuiFiles, "webui/dist")
	if err != nil {
		log.Printf("Failed to create webui filesystem: %v", err)
		return
	}

	// Serve static files
	fileServer := http.FileServer(http.FS(webuiFS))

	// Handle all routes that don't match API endpoints
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file, if not found serve index.html for SPA routing
		file, err := webuiFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			// File not found, serve index.html for client-side routing
			indexFile, err := webuiFS.Open("index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			defer indexFile.Close()
			w.Header().Set("Content-Type", "text/html")
			http.ServeContent(w, r, "index.html", time.Time{}, indexFile.(io.ReadSeeker))
			return
		}
		defer file.Close()
		fileServer.ServeHTTP(w, r)
	})
}
