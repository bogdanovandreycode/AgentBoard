package service

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

func TestFileAttachmentsHumanAgentFlowAndPermissions(t *testing.T) {
	f := setup(t)
	data := []byte{0, 1, 2, 255, 80, 75, 3, 4}
	file, err := f.svc.HumanUploadFile(f.ctx, f.project.ID, `C:\docs\report.zip`, data)
	if err != nil {
		t.Fatal(err)
	}
	if file.Name != "report.zip" || file.Size != int64(len(data)) || filepath.Dir(file.Path) != f.store.FileDirectory() {
		t.Fatalf("metadata: %+v", file)
	}
	stored, err := os.ReadFile(file.Path)
	if err != nil || !bytes.Equal(stored, data) {
		t.Fatalf("original file bytes changed: %v", err)
	}
	link := "[📄 report.zip](" + file.URL + ")"
	task, err := f.svc.HumanCreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Read attachment", Description: link, State: "features", Assignee: workerAssignee(f.a.ID)})
	if err != nil {
		t.Fatal(err)
	}
	details, err := f.svc.AgentTask(f.ctx, f.ac, task.ID)
	if err != nil || len(details.Files) != 1 || details.Files[0].ID != file.ID {
		t.Fatalf("agent metadata: %+v %v", details.Files, err)
	}
	_, read, err := f.svc.AgentFile(f.ctx, f.ac, task.ID, file.ID)
	if err != nil || !bytes.Equal(read, data) {
		t.Fatalf("agent read: %v", err)
	}
	if _, _, err := f.svc.AgentFile(f.ctx, f.bc, task.ID, file.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("foreign worker read: %v", err)
	}
	comment, err := f.svc.HumanComment(f.ctx, task.ID, link)
	if err != nil || len(comment.Files) != 1 {
		t.Fatalf("human comment: %+v %v", comment, err)
	}
	fromAgent, err := f.svc.AgentUploadFile(f.ctx, f.ac, task.ID, "answer.txt", []byte("AI output"))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.AgentFile(f.ctx, f.ac, task.ID, fromAgent.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("unattached file readable: %v", err)
	}
	agentLink := "[📄 answer.txt](" + fromAgent.URL + ")"
	agentComment, err := f.svc.AgentAddHistory(f.ctx, f.ac, task.ID, agentLink)
	if err != nil || len(agentComment.Files) != 1 {
		t.Fatalf("agent comment: %+v %v", agentComment, err)
	}
	human, err := f.svc.HumanTask(f.ctx, task.ID)
	if err != nil || len(human.History[len(human.History)-1].Files) != 1 {
		t.Fatalf("human cannot see AI attachment: %+v %v", human.History, err)
	}
	other := createTask(t, f, f.b, "features", "ai")
	secret, err := f.svc.AgentUploadFile(f.ctx, f.bc, other.ID, "private.txt", []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	secretLink := "[📄 private](" + secret.URL + ")"
	if _, err := f.svc.AgentAddHistory(f.ctx, f.ac, task.ID, secretLink); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("guessed file attached: %v", err)
	}
	if _, err := f.svc.CreateFeature(f.ctx, f.ac, core.FeatureInput{Title: "Leak", Reason: secretLink, Assignee: workerAssignee(f.a.ID)}); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("reason bypass: %v", err)
	}
	if _, err := f.svc.AgentAddArtifact(f.ctx, f.ac, task.ID, core.ArtifactInput{Name: secretLink, Path: "result.txt"}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.AgentFile(f.ctx, f.ac, task.ID, secret.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("system event bypass: %v", err)
	}
	for _, state := range []string{"backlog", "complete"} {
		if _, err := f.svc.HumanMove(f.ctx, task.ID, state); err != nil {
			t.Fatal(err)
		}
		if _, _, err := f.svc.AgentFile(f.ctx, f.ac, task.ID, file.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
			t.Fatalf("%s file read: %v", state, err)
		}
		if _, err := f.svc.AgentUploadFile(f.ctx, f.ac, task.ID, "hidden", nil); !errors.Is(err, core.ErrTaskNotAccessible) {
			t.Fatalf("%s upload: %v", state, err)
		}
	}
}

func TestFileAttachmentsValidationAndIndependentFields(t *testing.T) {
	f := setup(t)
	file, err := f.svc.HumanUploadFile(f.ctx, f.project.ID, "empty.txt", nil)
	if err != nil || file.Size != 0 {
		t.Fatalf("empty file: %+v %v", file, err)
	}
	link := "[📄 empty](" + file.URL + ")"
	task, err := f.svc.HumanCreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Independent fields", Description: "```", AITestInstructions: link, State: "features", Assignee: workerAssignee(f.a.ID)})
	if err != nil {
		t.Fatal(err)
	}
	details, err := f.svc.AgentTask(f.ctx, f.ac, task.ID)
	if err != nil || len(details.Files) != 1 {
		t.Fatalf("field parsing: %+v %v", details.Files, err)
	}
	if _, _, err := f.svc.AgentFile(f.ctx, f.ac, task.ID, file.ID); err != nil {
		t.Fatalf("test instruction file: %v", err)
	}
	refs, err := f.svc.FileReferences(f.ctx, f.project.ID, "`"+link+"`", "```text\n"+link+"\n```", "https://example.com"+file.URL)
	if err != nil || len(refs) != 0 {
		t.Fatalf("code/remote link: %+v %v", refs, err)
	}
	other, err := f.store.EnsureProject(f.ctx, "other", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.HumanCreateTask(f.ctx, other.ID, core.TaskInput{Title: "Foreign", Description: link}); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("foreign project reference: %v", err)
	}
	if _, err := f.svc.HumanUploadFile(f.ctx, f.project.ID, "large.bin", make([]byte, MaxFileBytes+1)); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("oversized file: %v", err)
	}
	if _, err := f.store.DB.Exec(`UPDATE files SET path='../escape' WHERE id=?`, file.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.HumanFile(f.ctx, file.ID); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("tampered path: %v", err)
	}
}
