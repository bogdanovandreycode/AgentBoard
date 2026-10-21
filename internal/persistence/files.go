package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/google/uuid"
)

func (s *Store) FileDirectory() string {
	if s.imageDir == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(s.imageDir), "files")
}

func (s *Store) CreateFile(ctx context.Context, file core.FileAttachment) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO files(id,project_id,name,path,mime_type,size,uploaded_task_id,created_at) VALUES(?,?,?,?,?,?,NULLIF(?,''),?)`,
		file.ID, file.ProjectID, file.Name, file.ID, file.MIMEType, file.Size, file.UploadedTaskID, file.CreatedAt)
	return err
}

func (s *Store) GetFile(ctx context.Context, fileID string) (core.FileAttachment, error) {
	var file core.FileAttachment
	var filename string
	err := s.DB.QueryRowContext(ctx, `SELECT id,project_id,name,path,mime_type,size,COALESCE(uploaded_task_id,''),created_at FROM files WHERE id=?`, fileID).
		Scan(&file.ID, &file.ProjectID, &file.Name, &filename, &file.MIMEType, &file.Size, &file.UploadedTaskID, &file.CreatedAt)
	if err == sql.ErrNoRows {
		return file, core.ErrNotFound
	}
	if err != nil {
		return file, err
	}
	if _, err := uuid.Parse(file.ID); err != nil || filename != file.ID || s.FileDirectory() == "" {
		return core.FileAttachment{}, fmt.Errorf("%w: invalid stored file path", core.ErrInvalidInput)
	}
	file.Path = filepath.Join(s.FileDirectory(), filename)
	file.URL = "/api/files/" + file.ID
	return file, nil
}
