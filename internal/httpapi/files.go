package httpapi

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/go-chi/chi/v5"
)

func (a *API) uploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, service.MaxFileBytes+(64<<10))
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			write(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "File exceeds the 50 MiB limit."})
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
	data, err := io.ReadAll(io.LimitReader(file, service.MaxFileBytes+1))
	if err != nil {
		writeError(w, err)
		return
	}
	if len(data) > service.MaxFileBytes {
		write(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "File exceeds the 50 MiB limit."})
		return
	}
	stored, err := a.Service.HumanUploadFile(r.Context(), chi.URLParam(r, "projectID"), header.Filename, data)
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, http.StatusCreated, stored)
}

func (a *API) file(w http.ResponseWriter, r *http.Request) {
	stored, data, err := a.Service.HumanFile(r.Context(), chi.URLParam(r, "fileID"))
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": stored.Name}))
	w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	modified, _ := time.Parse(time.RFC3339Nano, stored.CreatedAt)
	http.ServeContent(w, r, stored.Name, modified, bytes.NewReader(data))
}

func (a *API) revealFile(w http.ResponseWriter, r *http.Request) {
	if !allowFileReveal(w, r) {
		return
	}
	stored, _, err := a.Service.HumanFile(r.Context(), chi.URLParam(r, "fileID"))
	if err != nil {
		writeError(w, err)
		return
	}
	if err := revealStoredImage(stored.Path); err != nil {
		writeError(w, err)
		return
	}
	write(w, http.StatusOK, map[string]bool{"ok": true})
}
