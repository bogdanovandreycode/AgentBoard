# Cài đặt và khởi chạy trên Windows

## Trình cài đặt (được khuyến nghị)

Tải xuống `agentboard-VERSION-windows-amd64-setup.exe` từ trang [Bản phát hành](https://github.com/bogdanovandreycode/AgentBoard/releases). Trình hướng dẫn cài đặt sẽ đề xuất một thư mục; mặc định là `C:\AI\AgentBoard`. Nó sẽ sao chép `agentboard.exe` và tài liệu, tạo lối tắt và thêm thư mục đã chọn vào hệ thống `PATH`. Sau khi cài đặt, hãy mở một thiết bị đầu cuối mới để có lệnh `agentboard`.

Trong thư mục dự án của bạn chạy:

```powershell
agentboard init
agentboard open
```

Giao diện được tích hợp trong `agentboard.exe`; Không cần cài đặt riêng Go hoặc Node.js. Dữ liệu được lưu trữ trong `%AppData%\AgentBoard` và được giữ lại khi chương trình được cập nhật hoặc gỡ cài đặt. Gỡ cài đặt thông qua “Ứng dụng đã cài đặt” sẽ xóa các phím tắt và mục nhập khỏi `PATH`.

## Cài đặt Scoop

Nếu Scoop chưa được cài đặt, hãy mở PowerShell với tư cách người dùng thông thường của bạn và làm theo [hướng dẫn chính thức của Scoop](https://scoop.sh/). Nếu có những hạn chế trên máy tính công ty của bạn, hãy liên hệ với quản trị viên của bạn; AgentBoard cũng có thể được khởi chạy từ ZIP mà không cần Scoop.

Hoặc cài đặt AgentBoard qua Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Bạn có thể kiểm tra tính khả dụng của bảng kê khai trong Bản phát hành GitHub trên [trang Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Sau khi thêm tệp kê khai vào Scoop, nhóm có thể được cài đặt theo tên nhóm và cập nhật bằng lệnh `scoop update agentboard`.

## ZIP không có Scoop

Tải xuống `agentboard-VERSION-windows-amd64.zip` từ Bản phát hành, giải nén, chẳng hạn như vào `C:\Tools\AgentBoard`. Trong PowerShell trong thư mục dự án:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Đối với máy khách AI, hãy chỉ định đường dẫn đầy đủ đến `agentboard.exe` trong cấu hình MCP của nó nếu chương trình không nằm trong `PATH`.

## Xây dựng từ nguồn

Cài đặt phiên bản Go từ `go.mod` và Node.js 22 trở lên. Trong PowerShell ở thư mục gốc của kho lưu trữ:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` lắp ráp giao diện người dùng thành `internal/webui/dist`, chạy thử nghiệm Go và lắp ráp một `agentboard.exe`. Thứ tự rất quan trọng: giao diện được tích hợp vào hệ nhị phân khi Go được xây dựng. Nếu `agentboard.exe open` cũ đang chạy, hãy dừng nó trước khi xây dựng lại (Ctrl+C), nếu không Windows sẽ không cho phép bạn thay thế tệp.

## Dữ liệu ở đâu

- `%AppData%\AgentBoard\agentboard.db` - nhiệm vụ, dự án, công nhân và cài đặt. Bạn có thể chỉ định một tệp khác bằng cờ `--db`, nhưng `init`, `open`/`serve` và `mcp` phải có **cùng một đường dẫn**.
- `<ваш проект>\.agentboard\project.json` — mã định danh dự án. Tập tin này không chứa các nhiệm vụ.
- Máy chủ mặc định chỉ nghe `127.0.0.1:7337`. Chỉ định một địa chỉ khác `--addr` trước đường dẫn dự án: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` mở trang dự án và sử dụng lại máy chủ đang chạy tại địa chỉ đó. Nếu một số dự án được mở trong một trình duyệt, hãy chọn chúng qua danh sách bên trái.
