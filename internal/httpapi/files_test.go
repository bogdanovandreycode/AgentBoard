package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
)

func TestFileHTTPUploadAndOriginalDownload(t *testing.T) {
	ctx := context.Background()
	store, err := persistence.Open(filepath.Join(t.TempDir(), "files.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	project, err := store.EnsureProject(ctx, "files", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	handler := New(service.New(store))
	data := []byte("<html><script>test</script></html>")
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "Отчёт.html")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(data)
	writer.Close()
	request := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID+"/files", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("upload: %d %s", response.Code, response.Body.String())
	}
	var file core.FileAttachment
	if err := json.Unmarshal(response.Body.Bytes(), &file); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, file.URL, nil))
	disposition, params, err := mime.ParseMediaType(response.Header().Get("Content-Disposition"))
	if response.Code != http.StatusOK || !bytes.Equal(response.Body.Bytes(), data) || err != nil || disposition != "attachment" || params["filename"] != "Отчёт.html" || response.Header().Get("Content-Type") != "application/octet-stream" {
		t.Fatalf("download: %d %v %s", response.Code, response.Header(), response.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, file.URL, nil)
	request.Header.Set("Range", "bytes=0-5")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusPartialContent || !bytes.Equal(response.Body.Bytes(), data[:6]) {
		t.Fatalf("partial download: %d %s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest(http.MethodPost, file.URL+"/reveal", nil)
	request.Header.Set("Origin", "https://unrelated.example")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("cross-origin reveal: %d", response.Code)
	}
}
