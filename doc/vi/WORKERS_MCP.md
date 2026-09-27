# Công nhân và kết nối MCP

## Bước 1. Tạo một công nhân

Trong AgentBoard, mở **Công nhân → Thêm công nhân**. Chọn hồ sơ khách hàng: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama qua OpenCode, VS Code Copilot hoặc ứng dụng khách MCP khác. Hồ sơ điền vào các tính năng điển hình (mã, kiểm tra, Git); bạn có thể thay đổi chúng. Tên hiển thị trên bảng và `Slug` là mã định danh ngắn không có dấu cách cho lệnh MCP. Nhấp vào **Lưu**.

Một công nhân tương ứng với một tính cách AI. Tạo các nhân viên khác nhau cho các khách hàng hoặc nhóm khác nhau. Hồ sơ năng lực mô tả chuyên môn hóa nhưng không mở rộng quyền của AI đối với các giai đoạn nhiệm vụ.

## Bước 2. Copy cấu hình

Mở công nhân đã tạo. Khối **MCP Diagnostics** hiển thị đoạn cấu hình và tệp nơi cần thêm nó. Nhấp vào **Sao chép cấu hình MCP**. Nếu tệp đã tồn tại, hãy thêm máy chủ được đề xuất vào đối tượng `mcpServers`/`servers`/`mcp` hiện có mà không xóa các máy chủ khác.

Lệnh chính trông như thế này:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` sẽ trỏ đến thư mục bạn đã đăng ký qua `agentboard init`. `--worker` — `Slug` của nhân viên đã tạo. Máy khách MCP tự chạy lệnh này khi nó cần các công cụ. Trong trình duyệt AgentBoard, máy chủ web có thể chạy riêng.

Sau khi kiểm tra, khối chẩn đoán hiển thị đường dẫn tuyệt đối đến `agentboard.exe` đang chạy. Nó đặc biệt hữu ích khi cài đặt từ ZIP. Khi cài đặt qua Scoop, bạn có thể sử dụng lệnh `agentboard` nếu khách hàng thấy `PATH` tương tự.

## Bước 3: Thêm máy chủ vào máy khách của bạn

Giao diện công nhân đã có sẵn một đoạn. Dưới đây là lời giải thích về nơi nó được sử dụng:

| Khách hàng | Chèn vào đâu | Cách kiểm tra phía khách hàng |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, phần `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Mã Claude | `.mcp.json` trong thư mục dự án | `claude mcp list` |
| Song Tử CLI | `%USERPROFILE%\.gemini\settings.json`, đối tượng `mcpServers` | `/mcp list` trong Gemini CLI |
| Con trỏ | Dự án `.cursor\mcp.json` | danh sách máy chủ MCP trong cài đặt Con trỏ |
| Mã mở | Dự án `opencode.json`, đối tượng `mcp` | danh sách các công cụ MCP trong OpenCode |
| Phi công phụ mã VS | Dự án `.vscode\mcp.json`, đối tượng `servers` | lệnh **MCP: Máy chủ danh sách** |

Đối với Codex và Claude Code, một ví dụ với dự án `C:\Projects\MyFirstProject` và công nhân `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

Trong JSON, dấu gạch chéo ngược của Windows được nhân đôi; một đoạn được tạo sẵn từ giao diện sẽ tự động thực hiện việc này. Nếu bạn đang sử dụng `--db` với đế tùy chỉnh, hãy thêm nó vào cấu hình `args` MCP và chỉ định đường dẫn giống như khi khởi động `open`.

**Ollama** cung cấp mô hình cục bộ nhưng không thay thế ứng dụng khách MCP. Cấu hình **Ollama qua OpenCode** tạo cấu hình MCP cho OpenCode; cấu hình riêng OpenCode trên mô hình Ollama. Một ứng dụng khách tương thích MCP khác với Ollama cũng phù hợp.

## Bước 4: Kiểm tra kết nối của bạn

1. Mở thẻ công nhân và nhấp vào **Kiểm tra lại**. **Kiểm tra máy chủ** sẽ hiển thị số lượng công cụ MCP. Đây là công cụ kiểm tra và phát hiện giao thức nội bộ của các công cụ máy chủ.
2. Khởi động hoặc khởi động lại máy khách AI sau khi thêm tệp cấu hình. Yêu cầu anh ấy gọi `get_my_board`.
3. **Máy khách được kết nối** sẽ xuất hiện trong AgentBoard và một phiên mới sẽ xuất hiện trong danh sách. Chỉ điều này mới xác nhận kết nối của khách hàng của bạn. Nếu có công cụ nhưng máy khách chưa được kết nối, hãy kiểm tra đường dẫn đến chương trình, tên tệp cấu hình và cú pháp JSON/TOML của nó.

Bắt đầu phiên làm việc của bạn với `get_my_board`. AI không nhận nhiệm vụ từ `Backlog` và `Complete`, ngay cả khi nó biết ID của họ. AI không thể giả vờ là con người và không có lệnh chung “di chuyển đi đâu”. Làm việc với hệ thống tệp bên ngoài AgentBoard tùy thuộc vào khả năng và hộp cát của máy khách đã chọn.

Hướng dẫn chính thức dành cho khách hàng: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
