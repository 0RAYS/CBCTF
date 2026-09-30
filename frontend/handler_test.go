package frontend

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestGzipNegotiation(t *testing.T) {
	for header, want := range map[string]bool{
		"": false, "br": false, "gzip": true, "br, gzip;q=0.5": true,
		"gzip;q=0": false, "*;q=1, gzip;q=0": false, "gzip;q=0, *": false,
		"*": true, "*;q=0": false, "GZIP; Q=0.8": true,
		"gzip;q=invalid": false, "gzip;q=NaN": false, "gzip;q=2": false,
	} {
		t.Run(header, func(t *testing.T) {
			if got := acceptsGzip(header); got != want {
				t.Fatalf("acceptsGzip(%q) = %v, want %v", header, got, want)
			}
		})
	}
}

func TestStaticRepresentations(t *testing.T) {
	plain := []byte(strings.Repeat("export const text = 'CBCTF';\n", 100))
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write(plain); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	files := fstest.MapFS{
		"assets/app.js":    {Data: plain},
		"assets/app.js.gz": {Data: compressed.Bytes()},
		"assets/small.js":  {Data: plain},
		"index.html":       {Data: plain},
		"index.html.gz":    {Data: compressed.Bytes()},
	}
	handler := staticFileServer(files)
	for _, name := range []string{"/assets/app.js", "/"} {
		for _, encoding := range []string{"gzip", "identity", "gzip;q=0, *"} {
			t.Run(name+encoding, func(t *testing.T) {
				request := httptest.NewRequest(http.MethodGet, name, nil)
				request.Header.Set("Accept-Encoding", encoding)
				response := httptest.NewRecorder()
				handler.ServeHTTP(response, request)
				if response.Code != http.StatusOK || response.Header().Get("Vary") != "Accept-Encoding" {
					t.Fatalf("unexpected response: %d %v", response.Code, response.Header())
				}
				body := response.Body.Bytes()
				if encoding == "gzip" {
					if response.Header().Get("Content-Encoding") != "gzip" {
						t.Fatal("missing gzip content encoding")
					}
					reader, err := gzip.NewReader(bytes.NewReader(body))
					if err != nil {
						t.Fatal(err)
					}
					defer reader.Close()
					body, err = io.ReadAll(reader)
					if err != nil {
						t.Fatal(err)
					}
				} else if response.Header().Get("Content-Encoding") != "" {
					t.Fatal("identity response must not be marked gzip")
				}
				if !bytes.Equal(body, plain) || strings.Contains(response.Header().Get("Content-Type"), "gzip") {
					t.Fatal("representation changed the original content or MIME type")
				}
				if request.URL.Path != name {
					t.Fatal("handler mutated the caller's request")
				}
			})
		}
	}

	for _, scenario := range []struct {
		name, method, path, byteRange string
		status                        int
		encoding                      string
	}{
		{"head", "HEAD", "/assets/app.js", "", 200, "gzip"},
		{"uncompressed fallback", "GET", "/assets/small.js", "", 200, ""},
		{"range", "GET", "/assets/app.js", "bytes=0-5", 206, ""},
		{"missing", "GET", "/missing.js", "", 404, ""},
		{"index redirect", "GET", "/index.html", "", 301, ""},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			request := httptest.NewRequest(scenario.method, scenario.path, nil)
			request.Header.Set("Accept-Encoding", "gzip")
			if scenario.byteRange != "" {
				request.Header.Set("Range", scenario.byteRange)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != scenario.status || response.Header().Get("Content-Encoding") != scenario.encoding {
				t.Fatalf("unexpected response: %d %v", response.Code, response.Header())
			}
			if scenario.method == "HEAD" && response.Body.Len() != 0 {
				t.Fatal("HEAD response has a body")
			}
			if scenario.byteRange != "" && !bytes.Equal(response.Body.Bytes(), plain[:6]) {
				t.Fatal("range response changed")
			}
		})
	}
}
