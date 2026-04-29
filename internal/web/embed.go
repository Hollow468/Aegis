package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/frontend/*
var frontendFS embed.FS

// Handler returns an http.Handler that serves the embedded frontend.
// SPA fallback: any path that doesn't match a real file serves index.html.
func Handler() http.Handler {
	sub, err := fs.Sub(frontendFS, "dist/frontend")
	if err != nil {
		panic("failed to access embedded frontend: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Try to open the file; if it doesn't exist and looks like a page route, serve index.html
		if path != "/" && !strings.Contains(path, ".") {
			if f, err := sub.Open(strings.TrimPrefix(path, "/")); err != nil {
				r.URL.Path = "/"
			} else {
				f.Close()
			}
		}
		fileServer.ServeHTTP(w, r)
	})
}
