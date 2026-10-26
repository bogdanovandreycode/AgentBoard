package mcpserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	maxTaskImages     = 8
	maxTaskImageBytes = 8 << 20
)

type uploadImageInput struct {
	TaskID string `json:"task_id" jsonschema:"Assigned task that will receive the screenshot"`
	Path   string `json:"path,omitempty" jsonschema:"Local image file path; provide exactly one of path or base64"`
	Base64 string `json:"base64,omitempty" jsonschema:"Base64-encoded image bytes without a data URI prefix; provide exactly one of path or base64"`
	Name   string `json:"name,omitempty" jsonschema:"Optional screenshot name used in the Markdown link"`
}

type imageIDInput struct {
	TaskID  string `json:"task_id" jsonschema:"Assigned task referencing the screenshot"`
	ImageID string `json:"image_id" jsonschema:"Image ID returned in task image metadata"`
}

func (s *Server) registerImages() {
	addTool(s, &mcp.Tool{Name: "upload_image", Description: "Upload a screenshot for an assigned task from a local path or base64 (maximum 20 MiB). Converts it to WebP and returns a Markdown link. Insert that link in add_history to attach the screenshot to your comment."}, func(ctx context.Context, _ *mcp.CallToolRequest, in uploadImageInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		// Check task visibility before reading any local file or decoding its bytes.
		if _, err := s.Service.AgentTask(ctx, s.Agent, in.TaskID); err != nil {
			return failure(err)
		}
		data, name, err := imageUploadData(in)
		if err != nil {
			return failure(err)
		}
		img, err := s.Service.AgentUploadImage(ctx, s.Agent, in.TaskID, name, data)
		if err != nil {
			return failure(err)
		}
		return output(struct {
			Image        core.Image `json:"image"`
			Markdown     string     `json:"markdown"`
			Instructions string     `json:"instructions"`
		}{Image: img, Markdown: imageMarkdown(img), Instructions: "Include this Markdown link in add_history content to attach the screenshot. Refer to screenshots by name or number when explaining changes."})
	})
	addTool(s, &mcp.Tool{Name: "get_image", Description: "View an image referenced in an assigned task description, testing instructions or comment. Returns native WebP image content. Use this for screenshots omitted from get_task."}, func(ctx context.Context, _ *mcp.CallToolRequest, in imageIDInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		img, data, err := s.Service.AgentImage(ctx, s.Agent, in.TaskID, in.ImageID)
		if err != nil {
			return failure(err)
		}
		result, _, _ := output(img)
		result.Content = append(result.Content, &mcp.ImageContent{MIMEType: "image/webp", Data: data})
		return result, nil, nil
	})
}

func imageUploadData(in uploadImageInput) ([]byte, string, error) {
	return uploadData(in.Path, in.Base64, in.Name, service.MaxImageBytes)
}

func uploadData(path, encoded, name string, maxBytes int) ([]byte, string, error) {
	if (path == "") == (encoded == "") {
		return nil, "", fmt.Errorf("%w: provide exactly one of path or base64", core.ErrInvalidInput)
	}
	if path == "" {
		if len(encoded) > base64.StdEncoding.EncodedLen(maxBytes) {
			return nil, "", fmt.Errorf("%w: file exceeds %d MiB", core.ErrInvalidInput, maxBytes>>20)
		}
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil || len(data) > maxBytes {
			return nil, "", fmt.Errorf("%w: invalid base64 or file exceeds %d MiB", core.ErrInvalidInput, maxBytes>>20)
		}
		return data, name, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("read screenshot: %w", err)
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("read screenshot: %w", err)
	}
	if !stat.Mode().IsRegular() || stat.Size() > int64(maxBytes) {
		return nil, "", fmt.Errorf("%w: upload must be a regular file of at most %d MiB", core.ErrInvalidInput, maxBytes>>20)
	}
	data, err := io.ReadAll(io.LimitReader(f, int64(maxBytes)+1))
	if err != nil {
		return nil, "", fmt.Errorf("read screenshot: %w", err)
	}
	if len(data) > maxBytes {
		return nil, "", fmt.Errorf("%w: file exceeds %d MiB", core.ErrInvalidInput, maxBytes>>20)
	}
	if name == "" {
		name = filepath.Base(path)
	}
	return data, name, nil
}

func imageMarkdown(img core.Image) string {
	name := strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]", "\r", " ", "\n", " ").Replace(img.Name)
	return fmt.Sprintf("[🖼 %s](%s)", name, img.URL)
}

func (s *Server) taskOutput(ctx context.Context, task core.TaskDetails) (*mcp.CallToolResult, any, error) {
	result, _, _ := output(struct {
		core.TaskDetails
		Workflow service.AgentWorkflow
	}{task, s.Service.AgentWorkflow(ctx, s.Agent, task.Task)})
	type referencedImage struct {
		image  core.Image
		source string
	}
	var images []referencedImage
	seen := map[string]bool{}
	collect := func(refs []core.Image, source string) {
		for _, img := range refs {
			if !seen[img.ID] {
				seen[img.ID] = true
				images = append(images, referencedImage{image: img, source: source})
			}
		}
	}
	collect(task.Images, "task description/instructions")
	for _, entry := range task.History {
		collect(entry.Images, "comment "+entry.ID)
	}
	count, size := 0, 0
	for _, ref := range images {
		img := ref.image
		if count >= maxTaskImages || size >= maxTaskImageBytes {
			result.Content = append(result.Content, omittedImage(task.ID, img, "automatic preview limit"))
			continue
		}
		_, data, err := s.Service.AgentImage(ctx, s.Agent, task.ID, img.ID)
		if err != nil {
			result.Content = append(result.Content, omittedImage(task.ID, img, "image unavailable"))
			continue
		}
		if size+len(data) > maxTaskImageBytes {
			result.Content = append(result.Content, omittedImage(task.ID, img, "automatic preview byte limit"))
			continue
		}
		result.Content = append(result.Content,
			&mcp.TextContent{Text: fmt.Sprintf("Screenshot %q from %s (image_id: %s):", img.Name, ref.source, img.ID)},
			&mcp.ImageContent{MIMEType: "image/webp", Data: data},
		)
		count++
		size += len(data)
	}
	return result, nil, nil
}

func omittedImage(taskID string, img core.Image, reason string) *mcp.TextContent {
	return &mcp.TextContent{Text: fmt.Sprintf("Screenshot %q omitted (%s). View it with get_image using task_id=%q and image_id=%q.", img.Name, reason, taskID, img.ID)}
}
