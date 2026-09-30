package frontend

import (
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// FileServer serves embedded assets, preferring their build-time gzip representation.
// Callers strip /platform and set the entry/hashed-asset cache policy.
func FileServer() http.Handler {
	return staticFileServer(SubFS)
}

func staticFileServer(files fs.FS) http.Handler {
	fallback := http.FileServer(http.FS(files))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = "index.html"
		}
		switch path.Ext(name) {
		case ".js", ".css", ".html", ".svg":
			w.Header().Add("Vary", "Accept-Encoding")
			// Leave range requests and FileServer's canonical index redirects intact.
			if r.Header.Get("Range") == "" && !strings.HasSuffix(r.URL.Path, "/index.html") &&
				acceptsGzip(strings.Join(r.Header.Values("Accept-Encoding"), ",")) {
				if info, err := fs.Stat(files, name+".gz"); err == nil && !info.IsDir() {
					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
					w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
					compressed := r.Clone(r.Context())
					compressed.URL.Path = "/" + name + ".gz"
					fallback.ServeHTTP(w, compressed)
					return
				}
			}
		}
		fallback.ServeHTTP(w, r)
	})
}

func acceptsGzip(header string) bool {
	wildcard := false
	for encoding := range strings.SplitSeq(header, ",") {
		parts := strings.Split(encoding, ";")
		name := strings.ToLower(strings.TrimSpace(parts[0]))
		if name != "gzip" && name != "*" {
			continue
		}
		quality := 1.0
		for _, parameter := range parts[1:] {
			key, value, ok := strings.Cut(parameter, "=")
			if ok && strings.EqualFold(strings.TrimSpace(key), "q") {
				parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
				if err != nil || !(parsed >= 0 && parsed <= 1) {
					quality = 0
				} else {
					quality = parsed
				}
			}
		}
		if name == "gzip" {
			return quality > 0
		}
		wildcard = quality > 0
	}
	return wildcard
}
