# Bắt đầu tại đây

## AgentBoard là gì

Hãy tưởng tượng một bảng thông thường với các thẻ nhiệm vụ. Bạn tạo nhiệm vụ, phân công người chịu trách nhiệm và giám sát công việc. Nhân viên AI chỉ nhận nhiệm vụ được giao thông qua MCP và báo cáo tiến độ. Bạn quyết định khi nào nhiệm vụ cuối cùng đã sẵn sàng.

**Dự án** - một thư mục trên máy tính và một bảng riêng. **Nhiệm vụ** - một thẻ có mô tả, người chịu trách nhiệm và giai đoạn. **Worker** là tên logic của ứng dụng khách AI. **MCP** là cách máy khách AI kết nối với AgentBoard. **Scoop** là trình quản lý cài đặt phần mềm cho Windows.

Không cần mã. Bạn cần Windows, trình duyệt, PowerShell và để AI hoạt động, cần có ứng dụng khách AI được cài đặt có hỗ trợ cho máy chủ MCP cục bộ.

## Ra mắt lần đầu

1. Cài đặt ứng dụng theo [instructions](INSTALL.md).
2. Tạo thư mục dự án trong Explorer, ví dụ `C:\Projects\MyFirstProject`.
3. Mở thư mục này trong Explorer. Nhấp vào thanh địa chỉ, nhập `powershell` và nhấn Enter.
4. Trong cửa sổ mở ra, hãy thực hiện:

```powershell
agentboard init
agentboard open
```

5. Một trình duyệt sẽ mở ra với địa chỉ `http://127.0.0.1:7337`. Để cửa sổ PowerShell mở trong khi bạn sử dụng bảng trắng. Đóng cửa sổ sẽ dừng máy chủ cục bộ nhưng các tác vụ sẽ vẫn còn.

Nếu không tìm thấy lệnh `agentboard`, hãy đóng PowerShell và mở lại sau khi cài đặt Scoop. Khi cài đặt từ ZIP, hãy sử dụng đường dẫn đầy đủ tới `agentboard.exe`.

## Nhiệm vụ đầu tiên

Nhấp vào **Nhiệm vụ mới**, điền **Tiêu đề**, nếu cần **Mô tả**, sau đó **Lưu**. Một nhiệm vụ mới trong `Backlog` có sẵn cho con người. Để AI bắt đầu hoạt động, hãy phân công một công nhân và chuyển nhiệm vụ cho `Features`. Mô tả từng bước của các trường - [TASKS.md](TASKS.md).

## Công nhân đầu tiên

Mở **Công nhân → Thêm công nhân**. Chọn client AI bạn đang sử dụng (ví dụ Codex hoặc Claude Code), kiểm tra tên và ID rút gọn `Slug`, nhấp vào **Save**. Mở trình chạy đã tạo: có cấu hình MCP, nút sao chép và kiểm tra máy chủ. Copy cấu hình vào AI client theo [WORKERS_MCP.md](WORKERS_MCP.md). Sau khi kết nối, hãy yêu cầu khách hàng gọi `get_my_board`.

## Đọc gì tiếp theo

- [Nhiệm vụ, bài kiểm tra, câu chuyện và chuyên mục](TASKS.md)
- [Kết nối công nhân và kiểm tra MCP](WORKERS_MCP.md)
- [Nhập tác vụ từ JSON](IMPORT.md)
- [Ngôn ngữ, chủ đề, múi giờ và cập nhật](SETTINGS.md)
- [Các vấn đề điển hình](TROUBLESHOOTING.md)
