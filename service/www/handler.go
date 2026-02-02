package www

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// Handler returns an http.Handler that serves static files from the embedded
// filesystem with SPA fallback support. If a requested file does not exist,
// it serves index.html to support client-side routing.
func Handler() http.Handler {
	// Strip the "dist" prefix from the embedded filesystem so files are
	// served from the root path.
	distFS, err := fs.Sub(DistFS, "dist")
	if err != nil {
		// This should never happen since dist is embedded at compile time.
		panic("www: failed to create sub filesystem: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean the path and remove leading slash for fs.Open.
		urlPath := path.Clean(r.URL.Path)
		if urlPath == "/" {
			urlPath = "index.html"
		} else {
			urlPath = strings.TrimPrefix(urlPath, "/")
		}

		// Try to open the requested file.
		f, err := distFS.Open(urlPath)
		if err == nil {
			f.Close()
			// File exists, serve it normally.
			fileServer.ServeHTTP(w, r)
			return
		}

		// File doesn't exist. For SPA support, serve index.html.
		// This allows client-side routing to handle the path.
		indexFile, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		// Read index.html content.
		content, err := io.ReadAll(indexFile)
		if err != nil {
			http.Error(w, "failed to read index.html", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
	})
}
