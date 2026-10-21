package httpapi

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/go-chi/chi/v5"
)

func (a *API) uploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, service.MaxImageBytes+(64<<10))
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			write(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "Image exceeds the 20 MB limit."})
		} else {
			writeError(w, core.ErrInvalidInput)
		}
		return
	}
	defer r.MultipartForm.RemoveAll()
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, core.ErrInvalidInput)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, service.MaxImageBytes+1))
	if err != nil {
		writeError(w, err)
		return
	}
	if len(data) > service.MaxImageBytes {
		write(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "Image exceeds the 20 MB limit."})
		return
	}
	stored, err := a.Service.HumanUploadImage(r.Context(), chi.URLParam(r, "projectID"), header.Filename, data)
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, http.StatusCreated, stored)
}

func (a *API) image(w http.ResponseWriter, r *http.Request) {
	stored, data, err := a.Service.HumanImage(r.Context(), chi.URLParam(r, "imageID"))
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	if r.URL.Query().Get("download") == "1" {
		name := strings.TrimSuffix(stored.Name, filepath.Ext(stored.Name)) + ".webp"
		w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
	}
	modified, _ := time.Parse(time.RFC3339Nano, stored.CreatedAt)
	http.ServeContent(w, r, stored.Name, modified, bytes.NewReader(data))
}

func (a *API) revealImage(w http.ResponseWriter, r *http.Request) {
	// Revealing a file is an explicit action from this application's viewer.
	if !allowFileReveal(w, r) {
		return
	}
	stored, _, err := a.Service.HumanImage(r.Context(), chi.URLParam(r, "imageID"))
	if err != nil {
		writeError(w, err)
		return
	}
	if err = revealStoredImage(stored.Path); err != nil {
		writeError(w, err)
		return
	}
	write(w, http.StatusOK, map[string]bool{"ok": true})
}

func allowFileReveal(w http.ResponseWriter, r *http.Request) bool {
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		if err != nil || u.Host != r.Host {
			write(w, http.StatusForbidden, map[string]string{"error": "Origin is not allowed."})
			return false
		}
	}
	if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		write(w, http.StatusForbidden, map[string]string{"error": "Origin is not allowed."})
		return false
	}
	return true
}

func revealStoredImage(path string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("explorer.exe", "/select,"+path)
	case "darwin":
		command = exec.Command("open", "-R", path)
	default:
		command = exec.Command("xdg-open", filepath.Dir(path))
	}
	if err := command.Start(); err != nil {
		return err
	}
	go command.Wait()
	return nil
}
