# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Xem trước AgentBoard](../../assets/social-preview.png)

**[Tải xuống cho Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Trang tài liệu](https://bogdanovandreycode.github.io/AgentBoard/) · [Giấy phép MIT](../../LICENSE)**

AgentBoard là một ban nhiệm vụ địa phương nơi nhân viên con người và AI làm việc trên các nhiệm vụ chung nhưng có các quyền khác nhau. Ứng dụng được khởi chạy với một tệp `agentboard.exe`, mở giao diện web trong trình duyệt và cung cấp cho nhân viên một máy chủ MCP riêng thông qua `stdio`. Dữ liệu vẫn còn trên máy tính của bạn.

**[Bắt đầu từ đầu](START_HERE.md) · [Làm việc với các tác vụ](TASKS.md) · [Kết nối AI qua MCP](WORKERS_MCP.md) · [Nhập JSON](IMPORT.md) · [Cài đặt](SETTINGS.md) · [Giải quyết vấn đề](TROUBLESHOOTING.md)**

## Tính năng MVP

- Bảng dự án cục bộ với tìm kiếm tác vụ, nhập JSON, cột và thuộc tính tùy chỉnh.
- Quyền truy cập MCP cho một nhân viên cụ thể với khả năng chuyển đổi AI hạn chế và sự chấp nhận nhiệm vụ cuối cùng của con người.
- Lịch sử nhiệm vụ chung, hướng dẫn kiểm tra, tạo tác, chi phí AI và chẩn đoán kết nối nhân viên.
- Trình cài đặt Windows, tệp kê khai ZIP và Scoop; giao diện web được tích hợp vào tệp thực thi.

AgentBoard được thiết kế cho người dùng địa phương đáng tin cậy. Ứng dụng không lưu trữ các dự án trên đám mây và không tự khởi chạy ứng dụng khách AI; nếu cần, hãy kết nối máy khách tương thích MCP với nhân viên.

## Trong năm phút nữa

1. Tải xuống trình cài đặt `agentboard-VERSION-windows-amd64-setup.exe` từ [Release](https://github.com/bogdanovandreycode/AgentBoard/releases). Nó sẽ đề xuất một thư mục (theo mặc định là `C:\AI\AgentBoard`) và thêm nó vào `PATH`. Scoop và ZIP cũng có sẵn.
2. Mở PowerShell trong thư mục dự án của bạn, ví dụ `C:\Projects\MyApp`.
3. Thực thi `agentboard init` (đối với ZIP: đường dẫn đầy đủ đến `agentboard.exe` và `init`).
4. Thực thi `agentboard open`. `http://127.0.0.1:7337` sẽ mở ra.
5. Thêm tác vụ bằng nút **Tác vụ mới**. Đối với nhân viên AI, hãy mở **Công nhân → Thêm nhân viên**, chọn hồ sơ khách hàng và sao chép cấu hình MCP.

Nếu bạn chưa có thư mục dự án, hãy tạo một thư mục trong Windows Explorer. Một dự án có thể là bất kỳ thư mục nào, ngay cả khi không có Git và mã.

## Cài đặt qua Scoop

Trong PowerShell với [Scoop](https://scoop.sh/)] đã được cài đặt sau khi phát hành:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Đối với các nhà phát triển, có [bản dựng từ nguồn ](INSTALL.md). Quy trình phát hành tạo trình cài đặt, tệp kê khai ZIP và Scoop có SHA-256 từ cùng một cấu phần phần mềm. Cập nhật phiên bản được cài đặt qua Scoop: `scoop update agentboard` sau khi thêm tệp kê khai vào nhóm; chi tiết - [chuẩn bị phát hành](SCOOP_RELEASE.md).

## Cấu trúc của bảng như thế nào

`Backlog → Features → In progress → Testing → Verification → Complete`

Một người có thể di chuyển các nhiệm vụ xung quanh bảng. AI chỉ có thể di chuyển `Features → In progress → Testing → Verification`; Việc chấp nhận cuối cùng vào `Complete` được thực hiện bởi con người. Các cột người dùng dành cho con người: tác vụ trong đó vẫn ở trạng thái `Backlog` cho MCP. Chế độ thử nghiệm: AI, Human và Hybrid. Lịch sử, bài kiểm tra, hiện vật và chi phí AI được gắn liền với nhiệm vụ.

## Đội

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` đăng ký thư mục và chỉ ghi `.agentboard/project.json` vào đó. Dữ liệu sản xuất SQLite nằm trong thư mục cấu hình người dùng Windows (`%AppData%\AgentBoard\agentboard.db`), bên ngoài dự án và bên ngoài bản cài đặt Scoop. Việc xóa hoặc cập nhật gói sẽ không xóa dữ liệu này. Trước khi chuyển sang máy tính khác, hãy tạo một bản sao của cơ sở dữ liệu trong khi AgentBoard bị dừng.

Công nhân - tài khoản logic; Bản thân AgentBoard không chạy Codex, Claude hoặc bất kỳ ứng dụng khách AI nào khác. Máy khách bắt đầu quy trình MCP cục bộ cho một nhân viên cụ thể. MCP là ranh giới cấp phép của ứng dụng và để cách ly tệp, hãy sử dụng hộp cát máy khách AI.

## Dành cho nhà phát triển

Ngăn xếp: Go, SQLite, MCP Go SDK chính thức, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Xây dựng giao diện người dùng trước, sau đó đi: `./scripts/build.ps1`. Các tệp web được bao gồm trong tệp nhị phân thông qua `go:embed`. Kiến trúc và API được mô tả trong [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Giấy phép

AgentBoard là một dự án nguồn mở và miễn phí theo [giấy phép MIT](../../LICENSE). Được phép sử dụng, sửa đổi, phân tách và phân phối lại vì mục đích thương mại với điều kiện duy trì thông báo bản quyền và giấy phép.
