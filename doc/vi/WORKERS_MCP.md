# Công nhân và kết nối MCP

## Bước 1. Tạo một công nhân

Trong AgentBoard, mở **Công nhân → Thêm công nhân**. Chọn hồ sơ khách hàng: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama qua OpenCode, VS Code Copilot hoặc ứng dụng khách MCP khác. Hồ sơ điền vào các tính năng điển hình (mã, kiểm tra, Git); bạn có thể thay đổi chúng. Tên hiển thị trên bảng và `Slug` là mã định danh ngắn không có dấu cách cho lệnh MCP. Nhấp vào **Lưu**.

Một công nhân tương ứng với một tính cách AI. Tạo các nhân viên khác nhau cho các khách hàng hoặc nhóm khác nhau. Hồ sơ năng lực mô tả chuyên môn hóa nhưng không mở rộng quyền của AI đối với các giai đoạn nhiệm vụ.

## Bước 2. Copy cấu hình

Mở công nhân đã tạo. Khối **MCP Diagnostics** hiển thị đoạn cấu hình và tệp nơi cần thêm nó. Nhấp vào **Sao chép cấu hình MCP**. Nếu tệp đã tồn tại, hãy thêm máy chủ được đề xuất vào đối tượng `mcpServers`/`servers`/`mcp` hiện có mà không xóa các máy chủ khác.

Lệnh chính trông như thế này:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` của nhân viên đã tạo. Máy khách MCP tự chạy lệnh này khi nó cần các công cụ.

### MCP hoạt động mà không cần giao diện web

`agentboard mcp` là một quy trình stdio riêng biệt tự mở cơ sở dữ liệu AgentBoard. Nó không cần `agentboard open`/`serve`, nó không phụ thuộc vào thư mục mà giao diện web được khởi chạy và số lượng phiên bản của giao diện web đang chạy. Nhiều tác nhân có thể chạy đồng thời, mỗi tác nhân trong dự án riêng của chúng: mỗi máy khách chạy quy trình MCP của riêng mình và tất cả chúng đều sử dụng cùng một cơ sở dữ liệu SQLite một cách an toàn.

### Cách chọn dự án

Máy chủ MCP chạy trong bất kỳ thư mục nào. Lần đầu tiên bạn gọi công cụ này, nó sẽ bị ràng buộc với dự án:

1. `--project PATH`, nếu được chỉ định (bất kỳ thư mục nào bên trong dự án đã đăng ký đều phù hợp);
2. nếu không thì các thư mục làm việc (gốc) của khách hàng, nếu khách hàng báo cáo chúng;
3. nếu không thì thư mục mà máy khách đã khởi chạy máy chủ MCP.

`get_my_board` trả lại dự án và công nhân đã chọn cho `session` cùng với `binding_source`. Nếu thư mục làm việc của máy khách thay đổi, lệnh gọi công cụ tiếp theo sẽ thực hiện lại liên kết. Nếu không có `--worker`, nhân viên được kích hoạt duy nhất của dự án sẽ được sử dụng.

Nếu không thể liên kết, mỗi công cụ sẽ trả về một lỗi mà tác nhân phải chuyển cho người dùng; dữ liệu nhiệm vụ không được đọc trong trường hợp này:

| Mã | Ý nghĩa |
| --- | --- |
| `NO_PROJECTS` | Không có dự án nào trong cơ sở dữ liệu. Thực thi `agentboard init` trong thư mục dự án và tạo một worker. Không cần khởi động lại. |
| `PROJECT_REQUIRED` | Thư mục làm việc nằm ngoài các dự án đã đăng ký. Tin nhắn liệt kê các dự án. |
| `PROJECT_NOT_FOUND` | `--project` không nằm trong dự án đã đăng ký. |
| `PROJECT_AMBIGUOUS` | Thư mục công việc của khách hàng trải rộng trên nhiều dự án. Thêm `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Tạo, kích hoạt hoặc chọn một nhân viên thông qua `--worker`. |

Các cấu hình cho **Codex** và **Gemini CLI** là tùy chỉnh và được chia sẻ trên tất cả các thư mục, do đó chúng không có `--project`: Tạo một worker có cùng một slug trong mọi dự án mà ứng dụng khách này đang chạy. Các tệp cấp dự án (Mã Claude, Con trỏ, Mã mở, Mã VS) chứa `--project` của dự án của họ.

### Thiết kế rõ ràng và kiểm tra ngữ cảnh

Nếu `--project` được chỉ định và thư mục khởi động hoặc thư mục làm việc của khách hàng nằm trong một dự án khác, các công cụ sẽ yêu cầu xác nhận thông qua gợi ý MCP. Khách hàng không có biểu mẫu sẽ nhận được `PROJECT_CONTEXT_MISMATCH` trước khi truy cập các tác vụ. Sửa cấu hình, xóa `--project` để tự động phát hiện hoặc xác nhận rõ ràng rằng bạn có ý định làm việc với một dự án khác có cờ `--allow-project-mismatch`. Thay đổi thư mục làm việc sẽ đặt lại các xác nhận trước đó. Các dự án được xác định tự động không yêu cầu xác nhận.

**Cài đặt chung → Bối cảnh dự án MCP** kiểm soát các hoạt động kiểm tra này đối với tất cả các dự án và công nhân của cùng một cơ sở dữ liệu. Được bật theo mặc định. Nếu bị tắt, MCP sẽ sử dụng dự án được chỉ định mà không có xác nhận hoặc lỗi không khớp. Các phiên đã mở sẽ áp dụng cài đặt vào lần tiếp theo khi công cụ được gọi; khả năng hiển thị nhiệm vụ, chuyển tiếp và các quy tắc kiểm tra con người tiếp tục được tuân theo.

Sau khi kiểm tra, khối chẩn đoán hiển thị đường dẫn tuyệt đối đến `agentboard.exe` đang chạy. Nó đặc biệt hữu ích khi cài đặt từ ZIP. Khi cài đặt qua Scoop, bạn có thể sử dụng lệnh `agentboard` nếu máy khách nhìn thấy `PATH` tương tự.

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

Ví dụ về Codex (cấu hình người dùng, dự án được phát hiện tự động) và Claude Code (tệp dự án) với dự án `C:\Projects\MyFirstProject` và worker `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
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

Trong JSON, dấu gạch chéo ngược của Windows được nhân đôi; một đoạn được tạo sẵn từ giao diện sẽ tự động thực hiện việc này. Các cấu hình được tạo chứa `--db` thực tế của máy chủ web, do đó HTTP và MCP sử dụng cùng một cơ sở, ngay cả khi môi trường máy khách khác nhau.

**Ollama** cung cấp mô hình cục bộ nhưng không thay thế ứng dụng khách MCP. Cấu hình **Ollama qua OpenCode** tạo cấu hình MCP cho OpenCode; cấu hình riêng OpenCode trên mô hình Ollama. Một ứng dụng khách tương thích MCP khác với Ollama cũng phù hợp.

## Bước 4: Kiểm tra kết nối của bạn

1. Mở thẻ công nhân và nhấp vào **Kiểm tra lại**. **Kiểm tra máy chủ** chạy một tệp thực thi stdio thực sự với dự án, nhân viên và cơ sở này, kiểm tra sự bắt tay, danh sách các công cụ cũng như vị trí của dự án và nhân viên.
2. Khởi động hoặc khởi động lại máy khách AI sau khi thêm tệp cấu hình. Yêu cầu anh ấy gọi `get_my_board`.
3. **Máy khách được kết nối** sẽ xuất hiện trong AgentBoard và một phiên mới sẽ xuất hiện trong danh sách. Chỉ điều này mới xác nhận kết nối của khách hàng của bạn. Nếu có công cụ nhưng máy khách chưa được kết nối, hãy kiểm tra đường dẫn đến chương trình, tên tệp cấu hình và cú pháp JSON/TOML của nó.

Bắt đầu phiên làm việc của bạn với `get_my_board`. AI không nhận nhiệm vụ từ `Backlog` và `Complete`, ngay cả khi nó biết ID của họ. AI không thể giả vờ là con người và không có lệnh chung “di chuyển đi đâu”. Làm việc với hệ thống tệp bên ngoài AgentBoard tùy thuộc vào khả năng và hộp cát của máy khách đã chọn.

Mỗi công cụ khai báo `title` và tất cả bốn gợi ý hành vi MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) dưới dạng giá trị boolean rõ ràng. Công cụ đọc ở dạng chỉ đọc và bình thường; các công cụ ghi chỉ thêm bản ghi hoặc di chuyển các nhiệm vụ về phía trước, do đó không có công cụ nào gây gián đoạn hoặc vượt ra ngoài cơ sở AgentBoard cục bộ.

Hướng dẫn chính thức dành cho khách hàng: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
