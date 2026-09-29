package traffic

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var flagPattern = regexp.MustCompile(`(?i)[a-z][a-z0-9_-]{0,31}\{[^{}\r\n\x00-\x1f]{1,256}\}`)
var base64Pattern = regexp.MustCompile(`[A-Za-z0-9_+/-]{12,1000}={0,2}`)
var hexPattern = regexp.MustCompile(`\b[0-9a-fA-F]{16,1000}\b`)
var jsonStringPattern = regexp.MustCompile(`"(?:[^"\\\r\n]|\\.){1,512}"`)

func (a *analyzer) scanFlags(data []byte, e Evidence, encoding string) {
	a.scanAttackContent(data, e)
	values := flagPattern.FindAllString(string(data), maxFindings+1)
	for _, known := range a.known {
		if known != "" && bytes.Contains(data, []byte(known)) {
			values = append(values, known)
		}
	}
	for _, value := range values {
		verified := false
		for _, known := range a.known {
			if known == value {
				verified = true
				break
			}
		}
		key := value + "|" + evidenceKey(e)
		if a.seenFlags[key] {
			continue
		}
		if len(a.report.Flags) >= maxFindings {
			a.report.Truncated = true
			a.warn("flag_limit")
			return
		}
		a.seenFlags[key] = true
		a.report.Flags = append(a.report.Flags, FlagFinding{Value: value, Encoding: encoding, Verified: verified, Evidence: e})
	}
}

func (a *analyzer) scanEncoded(data []byte, e Evidence, encoding string) {
	// Bounded breadth/depth prevents adversarial nested encodings from expanding
	// without limit. Base64 is scanned as tokens, not guessed across packet gaps.
	type item struct {
		data     []byte
		encoding string
		depth    int
	}
	queue := []item{{data, encoding, 0}}
	seen := make(map[string]bool)
	budget := 4 * maxStreamBytes
	for len(queue) > 0 && budget > 0 {
		if a.ctx.Err() != nil {
			return
		}
		current := queue[0]
		queue = queue[1:]
		if len(current.data) > budget {
			a.report.Truncated = true
			a.warn("decode_budget")
			return
		}
		budget -= len(current.data)
		a.scanFlags(current.data, e, current.encoding)
		if current.depth >= 2 {
			continue
		}
		push := func(decoded []byte, codec string) {
			if len(decoded) == 0 || bytes.Equal(decoded, current.data) {
				return
			}
			key := packetFingerprint(decoded)
			if seen[key] {
				return
			}
			seen[key] = true
			if len(queue) >= 128 {
				a.report.Truncated = true
				a.warn("decode_budget")
				return
			}
			queue = append(queue, item{decoded, current.encoding + "/" + codec, current.depth + 1})
		}
		text := string(current.data)
		if decoded, err := url.QueryUnescape(text); err == nil {
			push([]byte(decoded), "url")
		}
		push([]byte(html.UnescapeString(text)), "html")
		for _, token := range base64Pattern.FindAll(current.data, 128) {
			for _, codec := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
				if decoded, err := codec.DecodeString(string(token)); err == nil {
					push(decoded, "base64")
					break
				}
			}
		}
		for _, token := range hexPattern.FindAllString(text, 128) {
			if decoded, err := hex.DecodeString(token); err == nil {
				push(decoded, "hex")
			}
		}
		for _, token := range jsonStringPattern.FindAll(current.data, 128) {
			var decoded string
			if json.Unmarshal(token, &decoded) == nil {
				push([]byte(decoded), "json")
			}
		}
	}
	if len(queue) > 0 {
		a.report.Truncated = true
		a.warn("decode_budget")
	}
}

func (a *analyzer) analyzeContent(data []byte, e Evidence) {
	a.scanEncoded(data, e, "raw")
	a.scanArchive(data, e)
	if len(data) >= 3 && data[0] == 0x16 && data[1] == 3 {
		a.warn("tls_payload_encrypted")
		return
	}
	reader := bufio.NewReader(bytes.NewReader(data))
	for {
		if a.ctx.Err() != nil {
			return
		}
		prefix, err := reader.Peek(5)
		if err != nil {
			return
		}
		var body io.ReadCloser
		var header http.Header
		record := HTTPRecord{Evidence: e}
		if string(prefix) == "HTTP/" {
			response, err := http.ReadResponse(reader, nil)
			if err != nil {
				return
			}
			body, header, record.Status = response.Body, response.Header, response.StatusCode
		} else {
			request, err := http.ReadRequest(reader)
			if err != nil {
				return
			}
			body, header = request.Body, request.Header
			record.Method, record.Host, record.URI = request.Method, request.Host, request.RequestURI
		}
		record.ContentType = header.Get("Content-Type")
		if len(a.report.HTTP) < maxFindings {
			a.report.HTTP = append(a.report.HTTP, record)
		} else {
			a.report.Truncated = true
			a.warn("http_limit")
			return
		}
		// ReadRequest/ReadResponse dechunk the body before content decoding.
		payload, err := io.ReadAll(io.LimitReader(body, maxStreamBytes+1))
		_ = body.Close()
		if len(payload) > maxStreamBytes {
			a.report.Truncated = true
			a.warn("http_body_limit")
			return
		}
		if err != nil {
			a.warn("http_body_incomplete")
		}
		a.scanEncoded(payload, e, "http_body")
		a.scanArchive(payload, e)
		codec := strings.ToLower(strings.TrimSpace(header.Get("Content-Encoding")))
		if codec != "" && codec != "identity" {
			var decoder io.ReadCloser
			switch codec {
			case "gzip":
				decoder, err = gzip.NewReader(bytes.NewReader(payload))
			case "deflate":
				decoder, err = zlib.NewReader(bytes.NewReader(payload))
				if err != nil {
					decoder = flate.NewReader(bytes.NewReader(payload))
					err = nil
				}
			default:
				a.warn("unsupported_content_encoding")
				continue
			}
			if err != nil {
				a.warn("content_decode_error")
				continue
			}
			decoded, readErr := io.ReadAll(io.LimitReader(decoder, maxStreamBytes+1))
			_ = decoder.Close()
			if len(decoded) > maxStreamBytes {
				a.report.Truncated = true
				a.warn("decompression_limit")
				continue
			}
			if readErr != nil {
				a.warn("content_decode_error")
				continue
			}
			a.scanEncoded(decoded, e, "http_body/"+codec)
		}
	}
}

func (a *analyzer) scanArchive(data []byte, e Evidence) {
	read := func(r io.ReadCloser, encoding string, budget *int) {
		defer r.Close()
		decoded, err := io.ReadAll(io.LimitReader(r, int64(*budget)+1))
		if len(decoded) > *budget {
			a.report.Truncated = true
			a.warn("decompression_limit")
			*budget = 0
			return
		}
		*budget -= len(decoded)
		if err != nil {
			a.warn("content_decode_error")
			return
		}
		a.scanEncoded(decoded, e, encoding)
	}
	budget := maxStreamBytes
	if bytes.HasPrefix(data, []byte{0x1f, 0x8b}) {
		if r, err := gzip.NewReader(bytes.NewReader(data)); err == nil {
			read(r, "gzip", &budget)
		}
	}
	if bytes.HasPrefix(data, []byte("PK\x03\x04")) {
		archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			a.warn("zip_incomplete")
			return
		}
		for i, file := range archive.File {
			if i >= 32 || budget <= 0 {
				a.report.Truncated = true
				a.warn("archive_limit")
				break
			}
			if a.ctx.Err() != nil {
				return
			}
			if file.FileInfo().IsDir() {
				continue
			}
			if r, err := file.Open(); err == nil {
				read(r, "zip", &budget)
			} else {
				a.warn("archive_entry_unreadable")
			}
		}
	}
}
