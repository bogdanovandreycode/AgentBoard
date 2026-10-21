package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/HugoSmits86/nativewebp"
	"github.com/bogdanovandreycode/agentboard/internal/core"
	_ "golang.org/x/image/bmp"
)

const (
	MaxImageBytes  = 20 << 20
	maxImageSide   = 16384
	maxImagePixels = 20_000_000
)

var imageReferencePattern = regexp.MustCompile(`!?\[((?:\\.|[^\]\\])*)\]\(\s*<?/api/images/([A-Za-z0-9_-]+)>?(?:\s+["'][^"']*["'])?\s*\)`)

// ImageReferences resolves the stable Markdown links without embedding image data in SQLite.
// References from another project and missing images are rejected before content is saved.
func (s *Service) ImageReferences(ctx context.Context, projectID, content string) ([]core.Image, error) {
	images := []core.Image{}
	for _, imageID := range imageReferenceIDs(content) {
		img, err := s.Store.GetImage(ctx, imageID)
		if err != nil {
			return nil, fmt.Errorf("%w: image reference %s", core.ErrInvalidInput, imageID)
		}
		if img.ProjectID != projectID {
			return nil, fmt.Errorf("%w: image belongs to another project", core.ErrInvalidInput)
		}
		images = append(images, img)
	}
	return images, nil
}

// ImageReferencesInTexts keeps Markdown parsing scoped to each independent field.
// An unmatched code fence in one field must not hide another field's attachments.
func (s *Service) ImageReferencesInTexts(ctx context.Context, projectID string, texts ...string) ([]core.Image, error) {
	images := []core.Image{}
	seen := map[string]bool{}
	for _, text := range texts {
		refs, err := s.ImageReferences(ctx, projectID, text)
		if err != nil {
			return nil, err
		}
		for _, img := range refs {
			if !seen[img.ID] {
				seen[img.ID] = true
				images = append(images, img)
			}
		}
	}
	return images, nil
}

func imageReferenceIDs(content string) []string {
	seen := map[string]bool{}
	var ids []string
	for _, match := range imageReferencePattern.FindAllStringSubmatch(withoutImageCode(content), -1) {
		imageID := strings.ToLower(match[2])
		if !seen[imageID] {
			seen[imageID] = true
			ids = append(ids, imageID)
		}
	}
	return ids
}

// Ignore fenced blocks and inline code spans, where a link is an example rather
// than an attachment. Fence closing runs may be longer than their opening run.
func withoutImageCode(content string) string {
	data := []byte(content)
	mask := func(start, end int) {
		for i := start; i < end; i++ {
			if data[i] != '\n' {
				data[i] = ' '
			}
		}
	}
	var fence byte
	fenceLength, offset := 0, 0
	for _, line := range strings.SplitAfter(content, "\n") {
		trimmed := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmed)
		marker, run := byte(0), 0
		if indent <= 3 && len(trimmed) > 0 && (trimmed[0] == '`' || trimmed[0] == '~') {
			marker = trimmed[0]
			for run < len(trimmed) && trimmed[run] == marker {
				run++
			}
		}
		if fence != 0 {
			mask(offset, offset+len(line))
			if marker == fence && run >= fenceLength && strings.TrimSpace(trimmed[run:]) == "" {
				fence = 0
			}
		} else if run >= 3 {
			fence, fenceLength = marker, run
			mask(offset, offset+len(line))
		}
		offset += len(line)
	}
	for start := 0; start < len(data); {
		if data[start] != '`' {
			start++
			continue
		}
		end := start
		for end < len(data) && data[end] == '`' {
			end++
		}
		length, next := end-start, end
		closed := false
		for next < len(data) {
			if data[next] != '`' {
				next++
				continue
			}
			closeEnd := next
			for closeEnd < len(data) && data[closeEnd] == '`' {
				closeEnd++
			}
			if closeEnd-next == length {
				mask(start, closeEnd)
				start, closed = closeEnd, true
				break
			}
			next = closeEnd
		}
		if !closed {
			start = end
		}
	}
	return string(data)
}

func (s *Service) HumanUploadImage(ctx context.Context, projectID, name string, data []byte) (core.Image, error) {
	if _, err := s.Store.GetProject(ctx, projectID); err != nil {
		return core.Image{}, err
	}
	return s.uploadImage(ctx, projectID, "", name, data)
}

func (s *Service) AgentUploadImage(ctx context.Context, a core.AgentContext, taskID, name string, data []byte) (core.Image, error) {
	task, err := s.accessibleTask(ctx, a, taskID)
	if err != nil {
		return core.Image{}, err
	}
	return s.uploadImage(ctx, task.ProjectID, task.ID, name, data)
}

func (s *Service) uploadImage(ctx context.Context, projectID, taskID, name string, data []byte) (core.Image, error) {
	if len(data) == 0 || len(data) > MaxImageBytes {
		return core.Image{}, fmt.Errorf("%w: image must be between 1 byte and 20 MiB", core.ErrInvalidInput)
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width <= 0 || config.Height <= 0 || config.Width > maxImageSide || config.Height > maxImageSide || int64(config.Width)*int64(config.Height) > maxImagePixels {
		return core.Image{}, fmt.Errorf("%w: unsupported image or excessive dimensions", core.ErrInvalidInput)
	}
	switch format {
	case "png", "jpeg", "gif", "bmp", "webp":
	default:
		return core.Image{}, fmt.Errorf("%w: unsupported image format", core.ErrInvalidInput)
	}
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return core.Image{}, fmt.Errorf("%w: invalid image data", core.ErrInvalidInput)
	}
	if err := ctx.Err(); err != nil {
		return core.Image{}, err
	}
	var encoded bytes.Buffer
	if err := nativewebp.Encode(&encoded, decoded, nil); err != nil {
		return core.Image{}, fmt.Errorf("encode image: %w", err)
	}
	if encoded.Len() > MaxImageBytes {
		return core.Image{}, fmt.Errorf("%w: converted image exceeds 20 MiB", core.ErrInvalidInput)
	}
	img := core.Image{ID: id(), ProjectID: projectID, UploadedTaskID: taskID, Name: imageName(name), Width: config.Width, Height: config.Height, CreatedAt: core.Now()}
	imageDir := s.Store.ImageDirectory()
	if imageDir == "" {
		return core.Image{}, fmt.Errorf("%w: image storage requires a file-backed database", core.ErrInvalidInput)
	}
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		return core.Image{}, err
	}
	img.Path = filepath.Join(imageDir, img.ID+".webp")
	img.URL = "/api/images/" + img.ID
	// The row is inserted only after closing the file, so adapters cannot read a
	// partial upload. O_EXCL also prevents overwriting an existing file.
	file, err := os.OpenFile(img.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return core.Image{}, err
	}
	persisted := false
	defer func() {
		if !persisted {
			os.Remove(img.Path)
		}
	}()
	_, writeErr := file.Write(encoded.Bytes())
	closeErr := file.Close()
	if writeErr != nil {
		return core.Image{}, writeErr
	}
	if closeErr != nil {
		return core.Image{}, closeErr
	}
	if err := ctx.Err(); err != nil {
		return core.Image{}, err
	}
	if err := s.Store.CreateImage(ctx, img); err != nil {
		return core.Image{}, err
	}
	persisted = true
	return img, nil
}

func imageName(name string) string {
	name = filepath.Base(strings.ReplaceAll(name, `\`, "/"))
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name)
	name = strings.TrimSpace(strings.TrimSuffix(name, filepath.Ext(name)))
	if name == "" || name == "." || name == ".." {
		name = "image"
	}
	if runes := []rune(name); len(runes) > 180 {
		name = string(runes[:180])
	}
	return name + ".webp"
}

func (s *Service) HumanImage(ctx context.Context, imageID string) (core.Image, []byte, error) {
	img, err := s.Store.GetImage(ctx, imageID)
	if err != nil {
		return core.Image{}, nil, err
	}
	data, err := readImageFile(img.Path)
	return img, data, err
}

func readImageFile(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > MaxImageBytes {
		return nil, fmt.Errorf("%w: invalid stored image file", core.ErrInvalidInput)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, MaxImageBytes+1))
	if err == nil && len(data) > MaxImageBytes {
		err = fmt.Errorf("%w: stored image is too large", core.ErrInvalidInput)
	}
	return data, err
}

func (s *Service) AgentImage(ctx context.Context, a core.AgentContext, taskID, imageID string) (core.Image, []byte, error) {
	task, err := s.accessibleTask(ctx, a, taskID)
	if err != nil {
		return core.Image{}, nil, err
	}
	img, err := s.Store.GetImage(ctx, imageID)
	if err != nil || img.ProjectID != task.ProjectID {
		return core.Image{}, nil, core.ErrTaskNotAccessible
	}
	referenced, err := s.taskReferencesImage(ctx, task, imageID)
	if err != nil {
		return core.Image{}, nil, err
	}
	if !referenced {
		return core.Image{}, nil, core.ErrTaskNotAccessible
	}
	data, err := readImageFile(img.Path)
	return img, data, err
}

// AgentValidateImageReferences prevents linking a guessed image from a hidden task
// to make it accessible. A newly uploaded image can only be attached to its task.
func (s *Service) AgentValidateImageReferences(ctx context.Context, a core.AgentContext, taskID, content string) error {
	images, err := s.ImageReferences(ctx, a.Project.ID, content)
	if err != nil || len(images) == 0 {
		return err
	}
	var tasks []core.Task
	if taskID != "" {
		task, err := s.accessibleTask(ctx, a, taskID)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	} else {
		board, err := s.AgentBoard(ctx, a)
		if err != nil {
			return err
		}
		for _, column := range board {
			tasks = append(tasks, column...)
		}
	}
	for _, img := range images {
		if taskID != "" && img.UploadedTaskID == taskID {
			continue
		}
		allowed := false
		for _, task := range tasks {
			referenced, err := s.taskReferencesImage(ctx, task, img.ID)
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

func (s *Service) taskReferencesImage(ctx context.Context, task core.Task, imageID string) (bool, error) {
	contains := func(content string) bool {
		for _, ref := range imageReferenceIDs(content) {
			if ref == imageID {
				return true
			}
		}
		return false
	}
	if contains(task.Description) || contains(task.AITestInstructions) || contains(task.HumanTestInstructions) {
		return true, nil
	}
	// Only explicit comments attach images. System history can interpolate names
	// supplied by agents (for example an artifact name) and cannot grant access.
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
