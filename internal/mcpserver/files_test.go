package mcpserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPFileUploadCommentAndChunkReading(t *testing.T) {
	f := newImageFixture(t)
	data := []byte("Report from AI: тест\n")
	result := f.call(t, "upload_file", map[string]any{"task_id": f.task.ID, "name": "report.txt", "base64": base64.StdEncoding.EncodeToString(data)}, false)
	var upload struct {
		File     core.FileAttachment `json:"file"`
		Markdown string              `json:"markdown"`
	}
	if err := json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &upload); err != nil {
		t.Fatal(err)
	}
	if upload.File.Name != "report.txt" || !strings.Contains(upload.Markdown, upload.File.URL) {
		t.Fatalf("file upload: %+v", upload)
	}
	f.call(t, "get_file", map[string]any{"task_id": f.task.ID, "file_id": upload.File.ID}, true)
	f.call(t, "add_history", map[string]any{"task_id": f.task.ID, "content": upload.Markdown}, false)
	var reconstructed []byte
	for offset := 0; ; {
		result := f.call(t, "get_file", map[string]any{"task_id": f.task.ID, "file_id": upload.File.ID, "offset": offset, "length": 7}, false)
		var chunk fileChunk
		if err := json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &chunk); err != nil {
			t.Fatal(err)
		}
		piece := []byte(chunk.Content)
		if chunk.Encoding == "base64" {
			var err error
			piece, err = base64.StdEncoding.DecodeString(chunk.Content)
			if err != nil {
				t.Fatal(err)
			}
		}
		reconstructed = append(reconstructed, piece...)
		if chunk.EOF {
			break
		}
		if chunk.NextOffset <= offset {
			t.Fatal("chunk did not advance")
		}
		offset = chunk.NextOffset
	}
	if !bytes.Equal(reconstructed, data) {
		t.Fatalf("chunk reconstruction changed bytes: %q", reconstructed)
	}
	path := filepath.Join(t.TempDir(), "binary.zip")
	if err := os.WriteFile(path, []byte{0, 255, 1}, 0o600); err != nil {
		t.Fatal(err)
	}
	f.call(t, "upload_file", map[string]any{"task_id": f.task.ID, "path": path}, false)
	f.call(t, "get_file", map[string]any{"task_id": f.task.ID, "file_id": upload.File.ID, "offset": -1}, true)
	f.call(t, "get_file", map[string]any{"task_id": f.task.ID, "file_id": upload.File.ID, "length": 1048577}, true)
	other, err := f.svc.Store.CreateTask(f.ctx, f.agent.Project.ID, core.TaskInput{Title: "Other worker", State: "features", Assignee: core.Assignee{Type: "worker", WorkerID: &f.other.ID}})
	if err != nil {
		t.Fatal(err)
	}
	denied := f.call(t, "upload_file", map[string]any{"task_id": other.ID, "path": "missing-file"}, true)
	if !strings.Contains(denied.Content[0].(*mcp.TextContent).Text, "TASK_NOT_ACCESSIBLE") {
		t.Fatal("upload read file before checking task permissions")
	}
	f.call(t, "get_file", map[string]any{"task_id": other.ID, "file_id": upload.File.ID}, true)
	task, err := f.svc.AgentTask(f.ctx, f.agent, f.task.ID)
	if err != nil || len(task.History[len(task.History)-1].Files) != 1 {
		t.Fatalf("file metadata not included in task comments: %+v %v", task.History, err)
	}
}
