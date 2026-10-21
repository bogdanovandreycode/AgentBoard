package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/google/uuid"
)

// ImageDirectory follows the selected database, so HTTP and MCP share image files.
func (s *Store) ImageDirectory() string { return s.imageDir }

func (s *Store) CreateImage(ctx context.Context, img core.Image) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO images(id,project_id,name,path,width,height,uploaded_task_id,created_at) VALUES(?,?,?,?,?,?,NULLIF(?,''),?)`,
		img.ID, img.ProjectID, img.Name, img.ID+".webp", img.Width, img.Height, img.UploadedTaskID, img.CreatedAt)
	return err
}

func (s *Store) GetImage(ctx context.Context, imageID string) (core.Image, error) {
	var img core.Image
	var filename string
	err := s.DB.QueryRowContext(ctx, `SELECT id,project_id,name,path,width,height,COALESCE(uploaded_task_id,''),created_at FROM images WHERE id=?`, imageID).
		Scan(&img.ID, &img.ProjectID, &img.Name, &filename, &img.Width, &img.Height, &img.UploadedTaskID, &img.CreatedAt)
	if err == sql.ErrNoRows {
		return img, core.ErrNotFound
	}
	if err != nil {
		return img, err
	}
	if _, err := uuid.Parse(img.ID); err != nil || filename != img.ID+".webp" || s.imageDir == "" {
		return core.Image{}, fmt.Errorf("%w: invalid stored image path", core.ErrInvalidInput)
	}
	img.Path = filepath.Join(s.imageDir, filename)
	img.URL = "/api/images/" + img.ID
	return img, nil
}
