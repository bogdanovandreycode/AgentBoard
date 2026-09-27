# Kiến trúc và API

AgentBoard là một tiến trình Go cục bộ chạy SQLite. API HTTP cho con người, bộ điều hợp MCP cho AI và giao diện web chia sẻ logic dịch vụ chung. SQLite - nguồn trạng thái; giao diện người dùng không bỏ qua máy chủ. Các công cụ AI có bề mặt quyền riêng biệt và trạng thái `Backlog`/`Complete` trong MCP không khả dụng. Bản ghi lịch sử từ `System` được tạo bởi chính ứng dụng.

## Thành phần

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - mô hình miền và lỗi.
- `internal/service` - chuyển đổi tác vụ, quyền, nhập và cài đặt.
- `internal/persistence` - SQLite và di chuyển.
- `internal/httpapi` - API HTTP của con người.
- `internal/mcpserver` - Công cụ MCP dành cho một công nhân riêng biệt.
- `web` - Giao diện người dùng React/TypeScript; phần lắp ráp kết thúc bằng `internal/webui/dist` và được bao gồm trong EXE.

## Các tuyến HTTP cơ bản

| Phương pháp và đường dẫn | Điểm đến |
| --- | --- |
| `GET /api/health` | Kiểm tra máy chủ. |
| `GET /api/projects` | Các dự án đã đăng ký |
| `GET /api/projects/{id}/board` | Ban dự án. |
| `GET/PUT /api/projects/{id}/settings` | Cài đặt, thứ tự và tên của các cột. |
| `POST /api/projects/{id}/tasks` | Tạo một nhiệm vụ. |
| `POST /api/projects/{id}/tasks/import` | Nhập nguyên tử phiên bản JSON 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Thẻ, thay đổi, xóa. |
| `POST /api/tasks/{id}/move` | Di chuyển bởi một người. |
| `GET/POST /api/projects/{id}/workers` | Danh sách và tạo một công nhân. |
| `GET /api/workers/{id}/mcp/check` | Thử nghiệm nội bộ về bắt tay và các công cụ MCP. |
| `GET /api/workers/{id}/sessions` | Chẩn đoán phiên. |
| `GET/POST /api/projects/{id}/properties` | Thuộc tính tùy chỉnh. |

HTTP dành cho người dùng đáng tin cậy cục bộ. Không xuất bản một cổng web lên Internet mà không có xác thực riêng, các hạn chế về mạng và HTTPS. Máy chủ MCP được khởi chạy bằng `stdio` cho một dự án và nhân viên cụ thể; bắt đầu với `get_my_board`. Các quyền của nó hạn chế các hành động trong AgentBoard nhưng không thay thế hộp cát hệ thống tệp của máy khách AI.

Các cột người dùng được lưu trữ riêng biệt với `tasks.state`: lõi giữ nhiệm vụ trong cột như vậy trong `backlog` và `tasks.board_column` xác định vị trí trên bảng Con người. Điều này duy trì mô hình chuyển đổi AI tương tự. Việc xóa một cột sẽ xóa các tác vụ `board_column`.
