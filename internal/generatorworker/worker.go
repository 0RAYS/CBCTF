// Package generatorworker is the small, statically linked runtime injected into
// generator images. It has no Kubernetes, database or Redis credentials.
package generatorworker

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const Port = 8080

type Request struct {
	TeamID         uint     `json:"team_id"`
	Flags          []string `json:"flags"`
	SourceRevision string   `json:"source_revision"`
}

func SourceRevision(path string) string {
	if info, err := os.Stat(path); err == nil {
		return fmt.Sprintf("%d:%d", info.Size(), info.ModTime().UnixNano())
	}
	return ""
}

func EncodeFlags(flags []string) string {
	values := make([]string, len(flags))
	for i, flag := range flags {
		values[i] = base64.StdEncoding.EncodeToString([]byte(flag))
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(values, ",")))
}

func Install(destination string) error {
	source, err := os.Executable()
	if err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func Serve(token string) error {
	if token == "" {
		return fmt.Errorf("worker token is required")
	}
	if err := os.MkdirAll("/root/mnt/attachments", 0700); err != nil {
		return err
	}
	if _, err := os.Stat("/root/mnt/generator.zip"); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		err = exec.CommandContext(ctx, "unzip", "-o", "/root/mnt/generator.zip", "-d", "/root").Run()
		cancel()
		if err != nil {
			return fmt.Errorf("unpack generator: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	var mu sync.Mutex
	loadedRevision := SourceRevision("/root/mnt/generator.zip")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("POST /generate", func(w http.ResponseWriter, r *http.Request) {
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+token)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !mu.TryLock() {
			http.Error(w, "busy", http.StatusConflict)
			return
		}
		defer mu.Unlock()
		var request Request
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&request); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
		archiveRevision := SourceRevision("/root/mnt/generator.zip")
		if request.SourceRevision != archiveRevision {
			http.Error(w, "generator source not yet visible", http.StatusConflict)
			return
		}
		if archiveRevision != loadedRevision {
			if archiveRevision == "" {
				http.Error(w, "restart worker after removing generator archive", http.StatusConflict)
				return
			}
			if err := exec.CommandContext(ctx, "unzip", "-o", "/root/mnt/generator.zip", "-d", "/root").Run(); err != nil {
				http.Error(w, "unpack generator failed", 500)
				return
			}
			loadedRevision = archiveRevision
		}
		team := strconv.FormatUint(uint64(request.TeamID), 10)
		output := filepath.Join("/root/mnt/attachments", team+".zip")
		if err := os.Remove(output); err != nil && !os.IsNotExist(err) {
			http.Error(w, "remove previous output", 500)
			return
		}
		defer os.Remove(output)
		command := exec.CommandContext(ctx, "/root/run.sh", team, EncodeFlags(request.Flags))
		command.Dir = "/root"
		// Discard script output, which can contain flags. Only the completed ZIP
		// is streamed, with Content-Length detecting incomplete transfers.
		if err := runCommand(command); err != nil {
			http.Error(w, "generator execution failed", 500)
			return
		}
		file, err := os.Open(output)
		if err != nil {
			http.Error(w, "generator did not produce an attachment", 500)
			return
		}
		defer file.Close()
		info, err := file.Stat()
		if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			http.Error(w, "invalid attachment", 500)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
		_, _ = io.Copy(w, file)
	})
	server := &http.Server{Addr: fmt.Sprintf(":%d", Port), Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 2 * time.Minute, IdleTimeout: 30 * time.Second}
	return server.ListenAndServe()
}
