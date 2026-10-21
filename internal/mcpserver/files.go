package mcpserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type uploadFileInput struct {
	TaskID string `json:"task_id" jsonschema:"Assigned task receiving the attachment"`
	Path   string `json:"path,omitempty" jsonschema:"Local file path; provide exactly one of path or base64"`
	Base64 string `json:"base64,omitempty" jsonschema:"Base64 bytes; provide exactly one of path or base64"`
	Name   string `json:"name,omitempty" jsonschema:"Original filename, including its extension"`
}

type getFileInput struct {
	TaskID string `json:"task_id"`
	FileID string `json:"file_id"`
	Offset int    `json:"offset,omitempty" jsonschema:"Byte offset, default 0"`
	Length int    `json:"length,omitempty" jsonschema:"Bytes to return, default 65536, maximum 1048576"`
}

type fileChunk struct {
	File       core.FileAttachment `json:"file"`
	Offset     int                 `json:"offset"`
	NextOffset int                 `json:"next_offset"`
	EOF        bool                `json:"eof"`
	Encoding   string              `json:"encoding"`
	Content    string              `json:"content"`
}

func (s *Server) registerFiles() {
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "upload_file", Description: "Upload any file for an assigned task from a local path or base64, up to 50 MiB. Preserves the original bytes. Insert the returned document Markdown link in add_history to attach it. For screenshots use upload_image."}, func(ctx context.Context, _ *mcp.CallToolRequest, in uploadFileInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		if _, err := s.Service.AgentTask(ctx, s.Agent, in.TaskID); err != nil {
			return failure(err)
		}
		data, name, err := uploadData(in.Path, in.Base64, in.Name, service.MaxFileBytes)
		if err != nil {
			return failure(err)
		}
		file, err := s.Service.AgentUploadFile(ctx, s.Agent, in.TaskID, name, data)
		if err != nil {
			return failure(err)
		}
		name = strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]", "\r", " ", "\n", " ").Replace(file.Name)
		return output(struct {
			File     core.FileAttachment `json:"file"`
			Markdown string              `json:"markdown"`
		}{File: file, Markdown: fmt.Sprintf("[📄 %s](%s)", name, file.URL)})
	})
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "get_file", Description: "Read a file attached to an assigned task's description, instructions or comment. Returns a bounded chunk encoded as utf-8 or base64 with next_offset and eof. Read following chunks using next_offset. Access follows task permissions."}, func(ctx context.Context, _ *mcp.CallToolRequest, in getFileInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		if in.Offset < 0 || in.Length < 0 || in.Length > 1<<20 {
			return failure(core.ErrInvalidInput)
		}
		file, data, err := s.Service.AgentFile(ctx, s.Agent, in.TaskID, in.FileID)
		if err != nil {
			return failure(err)
		}
		if in.Offset > len(data) {
			return failure(core.ErrInvalidInput)
		}
		if in.Length == 0 {
			in.Length = 65536
		}
		end := min(len(data), in.Offset+in.Length)
		chunk := data[in.Offset:end]
		encoding, content := "base64", base64.StdEncoding.EncodeToString(chunk)
		if utf8.Valid(chunk) && !strings.ContainsRune(string(chunk), '\x00') {
			encoding, content = "utf-8", string(chunk)
		}
		return output(fileChunk{File: file, Offset: in.Offset, NextOffset: end, EOF: end == len(data), Encoding: encoding, Content: content})
	})
}
