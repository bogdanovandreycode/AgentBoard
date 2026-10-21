package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

const MaxFileBytes = 50 << 20

var fileReferencePattern = regexp.MustCompile(`!?\[((?:\\.|[^\]\\])*)\]\(\s*<?/api/files/([A-Za-z0-9_-]+)>?(?:\s+["'][^"']*["'])?\s*\)`)

func fileReferenceIDs(content string) []string {
	var ids []string
	seen := map[string]bool{}
	for _, match := range fileReferencePattern.FindAllStringSubmatch(withoutImageCode(content), -1) {
		if id := match[2]; !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}

func (s *Service) FileReferences(ctx context.Context, projectID string, texts ...string) ([]core.FileAttachment, error) {
	files := []core.FileAttachment{}
	seen := map[string]bool{}
	for _, text := range texts {
		for _, id := range fileReferenceIDs(text) {
			if seen[id] {
				continue
			}
			file, err := s.Store.GetFile(ctx, id)
			if err != nil || file.ProjectID != projectID {
				return nil, fmt.Errorf("%w: invalid file reference", core.ErrInvalidInput)
			}
			seen[id] = true
			files = append(files, file)
		}
	}
	return files, nil
}

func (s *Service) HumanUploadFile(ctx context.Context, projectID, name string, data []byte) (core.FileAttachment, error) {
	if _, err := s.Store.GetProject(ctx, projectID); err != nil {
		return core.FileAttachment{}, err
	}
	return s.uploadFile(ctx, projectID, "", name, data)
}

func (s *Service) AgentUploadFile(ctx context.Context, agent core.AgentContext, taskID, name string, data []byte) (core.FileAttachment, error) {
	task, err := s.accessibleTask(ctx, agent, taskID)
	if err != nil {
		return core.FileAttachment{}, err
	}
	return s.uploadFile(ctx, task.ProjectID, task.ID, name, data)
}

func (s *Service) uploadFile(ctx context.Context, projectID, taskID, name string, data []byte) (core.FileAttachment, error) {
	if len(data) > MaxFileBytes {
		return core.FileAttachment{}, fmt.Errorf("%w: file exceeds 50 MiB", core.ErrInvalidInput)
	}
	dir := s.Store.FileDirectory()
	if dir == "" {
		return core.FileAttachment{}, fmt.Errorf("%w: file storage requires a file-backed database", core.ErrInvalidInput)
	}
	if err := ctx.Err(); err != nil {
		return core.FileAttachment{}, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return core.FileAttachment{}, err
	}
	file := core.FileAttachment{ID: id(), ProjectID: projectID, UploadedTaskID: taskID, Name: attachmentName(name), Size: int64(len(data)), MIMEType: http.DetectContentType(data), CreatedAt: core.Now()}
	file.Path, file.URL = filepath.Join(dir, file.ID), "/api/files/"+file.ID
	f, err := os.OpenFile(file.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return core.FileAttachment{}, err
	}
	persisted := false
	defer func() {
		if !persisted {
			os.Remove(file.Path)
		}
	}()
	_, writeErr := f.Write(data)
	closeErr := f.Close()
	if writeErr != nil {
		return core.FileAttachment{}, writeErr
	}
	if closeErr != nil {
		return core.FileAttachment{}, closeErr
	}
	if err := ctx.Err(); err != nil {
		return core.FileAttachment{}, err
	}
	if err := s.Store.CreateFile(ctx, file); err != nil {
		return core.FileAttachment{}, err
	}
	persisted = true
	return file, nil
}

func attachmentName(name string) string {
	name = filepath.Base(strings.ReplaceAll(name, `\`, "/"))
	name = strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name))
	if name == "" || name == "." || name == ".." {
		return "file"
	}
	if runes := []rune(name); len(runes) > 240 {
		return string(runes[:240])
	}
	return name
}

func (s *Service) HumanFile(ctx context.Context, fileID string) (core.FileAttachment, []byte, error) {
	file, err := s.Store.GetFile(ctx, fileID)
	if err != nil {
		return file, nil, err
	}
	data, err := readAttachmentFile(file.Path)
	return file, data, err
}

func readAttachmentFile(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > MaxFileBytes {
		return nil, core.ErrInvalidInput
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, MaxFileBytes+1))
	if len(data) > MaxFileBytes {
		return nil, core.ErrInvalidInput
	}
	return data, err
}

func (s *Service) AgentFile(ctx context.Context, agent core.AgentContext, taskID, fileID string) (core.FileAttachment, []byte, error) {
	task, err := s.accessibleTask(ctx, agent, taskID)
	if err != nil {
		return core.FileAttachment{}, nil, err
	}
	file, err := s.Store.GetFile(ctx, fileID)
	if err != nil || file.ProjectID != task.ProjectID {
		return core.FileAttachment{}, nil, core.ErrTaskNotAccessible
	}
	referenced, err := s.taskReferencesFile(ctx, task, fileID)
	if err != nil {
		return core.FileAttachment{}, nil, err
	}
	if !referenced {
		return core.FileAttachment{}, nil, core.ErrTaskNotAccessible
	}
	data, err := readAttachmentFile(file.Path)
	return file, data, err
}

func (s *Service) AgentValidateFileReferences(ctx context.Context, agent core.AgentContext, taskID, content string) error {
	files, err := s.FileReferences(ctx, agent.Project.ID, content)
	if err != nil || len(files) == 0 {
		return err
	}
	var tasks []core.Task
	if taskID != "" {
		task, err := s.accessibleTask(ctx, agent, taskID)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	} else {
		board, err := s.AgentBoard(ctx, agent)
		if err != nil {
			return err
		}
		for _, column := range board {
			tasks = append(tasks, column...)
		}
	}
	for _, file := range files {
		if taskID != "" && file.UploadedTaskID == taskID {
			continue
		}
		allowed := false
		for _, task := range tasks {
			referenced, err := s.taskReferencesFile(ctx, task, file.ID)
			if err != nil {
				return err
			}
			if referenced {
				allowed = true
				break
			}
		}
		if !allowed {
			return core.ErrTaskNotAccessible
		}
	}
	return nil
}

func (s *Service) taskReferencesFile(ctx context.Context, task core.Task, fileID string) (bool, error) {
	contains := func(content string) bool {
		for _, id := range fileReferenceIDs(content) {
			if id == fileID {
				return true
			}
		}
		return false
	}
	if contains(task.Description) || contains(task.AITestInstructions) || contains(task.HumanTestInstructions) {
		return true, nil
	}
	rows, err := s.Store.DB.QueryContext(ctx, `SELECT content FROM history_entries WHERE task_id=? AND entry_type IN ('human_comment','agent_comment')`, task.ID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return false, err
		}
		if contains(content) {
			return true, nil
		}
	}
	return false, rows.Err()
}
