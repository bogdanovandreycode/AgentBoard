package webui

import (
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
)

func Handler() http.Handler {
	dist, err := fs.Sub(FS, "dist")
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "." || name == "" {
			name = "index.html"
		}
		b, err := fs.ReadFile(dist, name)
		if err != nil {
			b, err = fs.ReadFile(dist, "index.html")
			name = "index.html"
		}
		if err != nil {
			http.Error(w, "embedded UI unavailable", 500)
			return
		}
		if typ := mime.TypeByExtension(path.Ext(name)); typ != "" {
			w.Header().Set("Content-Type", typ)
		}
		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
}
