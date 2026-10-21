package mcpserver

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type imageFixture struct {
	ctx   context.Context
	svc   *service.Service
	agent core.AgentContext
	other core.Worker
	task  core.Task
	cs    *mcp.ClientSession
}

func newImageFixture(t *testing.T) imageFixture {
	t.Helper()
	ctx := context.Background()
	st, err := persistence.Open(filepath.Join(t.TempDir(), "images.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	p, err := st.EnsureProject(ctx, "images", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	w, err := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Codex", Slug: "codex", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	other, err := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Other", Slug: "other", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	task, err := st.CreateTask(ctx, p.ID, core.TaskInput{Title: "Screenshot task", State: "features", Assignee: core.Assignee{Type: "worker", WorkerID: &w.ID}, TestingMode: "ai"})
	if err != nil {
		t.Fatal(err)
	}
	sess, err := st.StartSession(ctx, p, w, "image-test")
	if err != nil {
		t.Fatal(err)
	}
	a := core.AgentContext{Project: p, Worker: w, Session: sess}
	svc := service.New(st)
	srv := New(svc, a, "test")
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	ss, err := srv.MCP.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ss.Close() })
	client := mcp.NewClient(&mcp.Implementation{Name: "image-test", Version: "test"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	f := imageFixture{ctx: ctx, svc: svc, agent: a, other: other, task: task, cs: cs}
	f.call(t, "get_my_board", map[string]any{}, false)
	return f
}

func (f imageFixture) call(t *testing.T, name string, args map[string]any, wantError bool) *mcp.CallToolResult {
	t.Helper()
	r, err := f.cs.CallTool(f.ctx, &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if r.IsError != wantError {
		content, _ := json.Marshal(r.Content)
		t.Fatalf("%s IsError=%v, want %v: %s", name, r.IsError, wantError, content)
	}
	return r
}

func screenshotPNG(t *testing.T, shade uint8) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 12, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 12; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(x * 20), G: shade, B: uint8(y * 20), A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

type uploadResult struct {
	Image    core.Image `json:"image"`
	Markdown string     `json:"markdown"`
}

func decodeUpload(t *testing.T, result *mcp.CallToolResult) uploadResult {
	t.Helper()
	var upload uploadResult
	if err := json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &upload); err != nil {
		t.Fatal(err)
	}
	if upload.Image.ID == "" || !strings.Contains(upload.Markdown, upload.Image.URL) {
		t.Fatalf("invalid upload response: %+v", upload)
	}
	return upload
}

func nativeImages(t *testing.T, result *mcp.CallToolResult) int {
	t.Helper()
	count := 0
	for _, block := range result.Content {
		if img, ok := block.(*mcp.ImageContent); ok {
			if img.MIMEType != "image/webp" || len(img.Data) < 12 || string(img.Data[0:4]) != "RIFF" || string(img.Data[8:12]) != "WEBP" {
				t.Fatal("MCP did not return valid native WebP image content")
			}
			count++
		}
	}
	return count
}

func TestMCPScreenshotUploadCommentAndNativeViewing(t *testing.T) {
	f := newImageFixture(t)
	if instructions := f.cs.InitializeResult().Instructions; !strings.Contains(instructions, "upload_image") || !strings.Contains(instructions, "get_image") {
		t.Fatal("missing screenshot instructions")
	}
	upload := decodeUpload(t, f.call(t, "upload_image", map[string]any{
		"task_id": f.task.ID, "name": "Screenshot 1.png", "base64": base64.StdEncoding.EncodeToString(screenshotPNG(t, 20)),
	}, false))
	if upload.Image.Width != 12 || upload.Image.Height != 10 || filepath.Ext(upload.Image.Path) != ".webp" {
		t.Fatalf("unexpected image metadata: %+v", upload.Image)
	}
	if _, err := os.Stat(upload.Image.Path); err != nil {
		t.Fatalf("converted screenshot was not stored on disk: %v", err)
	}
	// Uploading alone does not grant read access until the link is attached to the task.
	f.call(t, "get_image", map[string]any{"task_id": f.task.ID, "image_id": upload.Image.ID}, true)
	comment := f.call(t, "add_history", map[string]any{"task_id": f.task.ID, "content": "Screenshot 1 shows the bug: " + upload.Markdown}, false)
	var entry core.HistoryEntry
	if err := json.Unmarshal([]byte(comment.Content[0].(*mcp.TextContent).Text), &entry); err != nil {
		t.Fatal(err)
	}
	if entry.ActorType != "agent" || entry.WorkerID == nil || *entry.WorkerID != f.agent.Worker.ID {
		t.Fatalf("wrong screenshot comment attribution: %+v", entry)
	}
	description := "Screenshot 1: fix this area. " + upload.Markdown
	if _, err := f.svc.HumanUpdateTask(f.ctx, f.task.ID, core.TaskUpdate{Description: &description}); err != nil {
		t.Fatal(err)
	}
	result := f.call(t, "get_task", map[string]any{"task_id": f.task.ID}, false)
	if count := nativeImages(t, result); count != 1 {
		t.Fatalf("get_task returned %d native images, want one deduplicated image", count)
	}
	var details core.TaskDetails
	if err := json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &details); err != nil {
		t.Fatal(err)
	}
	commentImages := 0
	for _, h := range details.History {
		if h.ID == entry.ID {
			commentImages = len(h.Images)
		}
	}
	if len(details.Images) != 1 || commentImages != 1 {
		t.Fatalf("missing structured image references: %+v", details)
	}
	if count := nativeImages(t, f.call(t, "get_image", map[string]any{"task_id": f.task.ID, "image_id": upload.Image.ID}, false)); count != 1 {
		t.Fatalf("get_image returned %d native images", count)
	}
	// Human uploads must also become native image content for the assigned AI worker.
	human, err := f.svc.HumanUploadImage(f.ctx, f.agent.Project.ID, "Human screenshot.png", screenshotPNG(t, 160))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.HumanComment(f.ctx, f.task.ID, "Screenshot 3: human explanation. "+imageMarkdown(human)); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "Screenshot 2.png")
	if err := os.WriteFile(path, screenshotPNG(t, 80), 0600); err != nil {
		t.Fatal(err)
	}
	local := decodeUpload(t, f.call(t, "upload_image", map[string]any{"task_id": f.task.ID, "path": path}, false))
	f.call(t, "add_history", map[string]any{"task_id": f.task.ID, "content": "Screenshot 2: fixed result. " + local.Markdown}, false)
	if count := nativeImages(t, f.call(t, "get_task", map[string]any{"task_id": f.task.ID}, false)); count != 3 {
		t.Fatalf("get_task returned %d native images, want human and AI screenshots", count)
	}
}

func TestMCPImageToolsRespectTaskAndReferencePermissions(t *testing.T) {
	f := newImageFixture(t)
	upload := decodeUpload(t, f.call(t, "upload_image", map[string]any{
		"task_id": f.task.ID, "base64": base64.StdEncoding.EncodeToString(screenshotPNG(t, 40)),
	}, false))
	f.call(t, "add_history", map[string]any{"task_id": f.task.ID, "content": upload.Markdown}, false)
	otherTask, err := f.svc.Store.CreateTask(f.ctx, f.agent.Project.ID, core.TaskInput{Title: "Foreign", State: "features", Assignee: core.Assignee{Type: "worker", WorkerID: &f.other.ID}})
	if err != nil {
		t.Fatal(err)
	}
	for _, check := range []struct {
		tool string
		args map[string]any
	}{
		{tool: "get_task", args: map[string]any{"task_id": otherTask.ID}},
		{tool: "get_image", args: map[string]any{"task_id": otherTask.ID, "image_id": upload.Image.ID}},
		{tool: "upload_image", args: map[string]any{"task_id": otherTask.ID, "path": filepath.Join(t.TempDir(), "does-not-exist.png")}},
	} {
		denied := f.call(t, check.tool, check.args, true)
		if text := denied.Content[0].(*mcp.TextContent).Text; !strings.Contains(text, core.ErrTaskNotAccessible.Error()) {
			t.Fatalf("%s did not deny task access before processing screenshot: %s", check.tool, text)
		}
	}
	// An image attached to one task cannot be fetched through another visible task.
	ownSecond, err := f.svc.Store.CreateTask(f.ctx, f.agent.Project.ID, core.TaskInput{Title: "Second", State: "features", Assignee: core.Assignee{Type: "worker", WorkerID: &f.agent.Worker.ID}})
	if err != nil {
		t.Fatal(err)
	}
	f.call(t, "get_image", map[string]any{"task_id": ownSecond.ID, "image_id": upload.Image.ID}, true)

	otherSession, err := f.svc.Store.StartSession(f.ctx, f.agent.Project, f.other, "other-image-test")
	if err != nil {
		t.Fatal(err)
	}
	otherAgent := core.AgentContext{Project: f.agent.Project, Worker: f.other, Session: otherSession}
	private, err := f.svc.AgentUploadImage(f.ctx, otherAgent, otherTask.ID, "private.png", screenshotPNG(t, 120))
	if err != nil {
		t.Fatal(err)
	}
	privateLink := imageMarkdown(private)
	if _, err := f.svc.HumanComment(f.ctx, otherTask.ID, privateLink); err != nil {
		t.Fatal(err)
	}
	f.call(t, "get_image", map[string]any{"task_id": f.task.ID, "image_id": private.ID}, true)
	f.call(t, "add_history", map[string]any{"task_id": f.task.ID, "content": privateLink}, true)

	for _, state := range []string{"backlog", "complete"} {
		if _, err := f.svc.HumanMove(f.ctx, f.task.ID, state); err != nil {
			t.Fatal(err)
		}
		f.call(t, "get_task", map[string]any{"task_id": f.task.ID}, true)
		f.call(t, "get_image", map[string]any{"task_id": f.task.ID, "image_id": upload.Image.ID}, true)
		f.call(t, "upload_image", map[string]any{"task_id": f.task.ID, "base64": "AA=="}, true)
	}
}

func TestMCPTaskImagePreviewLimitOffersGetImageFallback(t *testing.T) {
	f := newImageFixture(t)
	var links []string
	for i := 0; i < maxTaskImages+1; i++ {
		img, err := f.svc.AgentUploadImage(f.ctx, f.agent, f.task.ID, "Screenshot.png", screenshotPNG(t, uint8(i*20)))
		if err != nil {
			t.Fatal(err)
		}
		links = append(links, imageMarkdown(img))
	}
	f.call(t, "add_history", map[string]any{"task_id": f.task.ID, "content": strings.Join(links, "\n")}, false)
	result := f.call(t, "get_task", map[string]any{"task_id": f.task.ID}, false)
	if count := nativeImages(t, result); count != maxTaskImages {
		t.Fatalf("returned %d previews, want %d", count, maxTaskImages)
	}
	last := result.Content[len(result.Content)-1].(*mcp.TextContent).Text
	if !strings.Contains(last, "omitted") || !strings.Contains(last, "get_image") || !strings.Contains(last, f.task.ID) {
		t.Fatalf("missing image fallback instructions: %q", last)
	}
}

func TestMCPImageInputRejectsAmbiguousAndOversizedFiles(t *testing.T) {
	for _, in := range []uploadImageInput{{}, {Path: "image.png", Base64: "AA=="}, {Base64: "not base64"}} {
		if _, _, err := imageUploadData(in); err == nil {
			t.Fatalf("accepted invalid input: %+v", in)
		}
	}
	path := filepath.Join(t.TempDir(), "too-large.png")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(service.MaxImageBytes + 1); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	if _, _, err := imageUploadData(uploadImageInput{Path: path}); err == nil {
		t.Fatal("accepted oversized screenshot")
	}
}
