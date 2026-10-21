package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
)

func TestImageUploadDescriptionCommentAndDownload(t *testing.T) {
	ctx := context.Background()
	store, err := persistence.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	project, err := store.EnsureProject(ctx, "images", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	handler := New(service.New(store))
	var pngData bytes.Buffer
	pixels := image.NewNRGBA(image.Rect(0, 0, 40, 24))
	pixels.SetNRGBA(4, 5, color.NRGBA{R: 220, G: 50, B: 20, A: 180})
	if err = png.Encode(&pngData, pixels); err != nil {
		t.Fatal(err)
	}
	upload := func(projectID string, data []byte) *httptest.ResponseRecorder {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", "Скрин 1.png")
		if err != nil {
			t.Fatal(err)
		}
		part.Write(data)
		writer.Close()
		request := httptest.NewRequest(http.MethodPost, "/api/projects/"+projectID+"/images", &body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		return response
	}
	response := upload(project.ID, pngData.Bytes())
	if response.Code != http.StatusCreated {
		t.Fatalf("upload: %d %s", response.Code, response.Body.String())
	}
	var stored core.Image
	if err = json.Unmarshal(response.Body.Bytes(), &stored); err != nil {
		t.Fatal(err)
	}
	if stored.Width != 40 || stored.Height != 24 || !filepath.IsAbs(stored.Path) || filepath.Ext(stored.Path) != ".webp" {
		t.Fatalf("invalid metadata: %+v", stored)
	}
	data, err := os.ReadFile(stored.Path)
	if err != nil || len(data) < 12 || string(data[8:12]) != "WEBP" {
		t.Fatalf("not a stored WebP: %v", err)
	}
	call := func(method, path string, value any) *httptest.ResponseRecorder {
		body, _ := json.Marshal(value)
		request := httptest.NewRequest(method, path, bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		result := httptest.NewRecorder()
		handler.ServeHTTP(result, request)
		return result
	}
	link := "[🖼 Скрин 1](" + stored.URL + ") Тут баг"
	created := call(http.MethodPost, "/api/projects/"+project.ID+"/tasks", core.TaskInput{Title: "With screenshot", Description: link})
	if created.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", created.Code, created.Body.String())
	}
	var task core.Task
	json.Unmarshal(created.Body.Bytes(), &task)
	comment := call(http.MethodPost, "/api/tasks/"+task.ID+"/comments", map[string]string{"content": link})
	if comment.Code != http.StatusCreated {
		t.Fatalf("comment: %d %s", comment.Code, comment.Body.String())
	}
	details := call(http.MethodGet, "/api/tasks/"+task.ID, nil)
	var got core.TaskDetails
	if err = json.Unmarshal(details.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if details.Code != http.StatusOK || got.Description != link || len(got.Images) != 1 {
		t.Fatalf("details: %d %+v", details.Code, got)
	}
	foundComment := false
	for _, entry := range got.History {
		if entry.Content == link && len(entry.Images) == 1 && entry.Images[0].ID == stored.ID {
			foundComment = true
		}
	}
	if !foundComment {
		t.Fatalf("missing comment image metadata: %+v", got.History)
	}
	download := call(http.MethodGet, stored.URL+"?download=1", nil)
	if download.Code != http.StatusOK || download.Header().Get("Content-Type") != "image/webp" || !strings.Contains(download.Header().Get("Content-Disposition"), "attachment;") || !bytes.Equal(download.Body.Bytes(), data) {
		t.Fatalf("download: %d %v", download.Code, download.Header())
	}
	if bad := upload(project.ID, []byte("<svg>not a raster image</svg>")); bad.Code != http.StatusBadRequest {
		t.Fatalf("invalid upload: %d %s", bad.Code, bad.Body.String())
	}
	if missing := upload("missing", pngData.Bytes()); missing.Code != http.StatusNotFound {
		t.Fatalf("missing project: %d %s", missing.Code, missing.Body.String())
	}
	other, err := store.EnsureProject(ctx, "other", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	foreign := call(http.MethodPost, "/api/projects/"+other.ID+"/tasks", core.TaskInput{Title: "Foreign screenshot", Description: link})
	if foreign.Code < 400 {
		t.Fatalf("cross-project image accepted: %d", foreign.Code)
	}
	request := httptest.NewRequest(http.MethodPost, stored.URL+"/reveal", nil)
	request.Header.Set("Origin", "https://unrelated.example")
	reveal := httptest.NewRecorder()
	handler.ServeHTTP(reveal, request)
	if reveal.Code != http.StatusForbidden {
		t.Fatalf("cross-origin reveal: %d", reveal.Code)
	}
}
