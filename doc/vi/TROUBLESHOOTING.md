# Giải quyết vấn đề

| Triệu chứng | Kiểm tra những gì |
| --- | --- |
| Không tìm thấy `agentboard` | Khởi động lại PowerShell sau Scoop. Khi cài đặt từ ZIP, hãy sử dụng đường dẫn đầy đủ tới `agentboard.exe`. |
| Trang web hiển thị giao diện cũ sau khi xây dựng | Dừng máy chủ đang chạy Ctrl+C. Thực thi `./scripts/build.ps1`, khởi chạy một tệp nhị phân mới. Vite phải build trước Go vì giao diện được tích hợp sẵn trong EXE. Làm mới trang Ctrl+F5. |
| Cổng 7337 bận | AgentBoard có thể đã chạy. Mở `http://127.0.0.1:7337` hoặc kết thúc quá trình cũ. Đối với cổng khác, hãy sử dụng `--addr`. |
| Dự án không tìm thấy | Trong thư mục mong muốn, thực thi `agentboard init`. Sau đó là `agentboard open` từ nó hoặc `agentboard open C:\путь\к\проекту`. |
| Người lao động không thấy nhiệm vụ | Nhiệm vụ phải được giao cho nhân viên cụ thể này và được đặt trong `Features`, `In progress`, `Testing` hoặc `Verification`. AI không thấy `Backlog`, `Complete` và các cột tùy chỉnh. |
| Kiểm tra máy chủ MCP tồn tại nhưng máy khách không được kết nối | Kiểm tra máy chủ không kiểm tra cài đặt máy khách bên ngoài. Khởi động lại máy khách, kiểm tra tệp cấu hình của nó, đường dẫn đến `agentboard.exe`, `--project`, `--worker` và `--db` chung. Yêu cầu gọi `get_my_board`. |
| Công Nhân Ngoại Tuyến | Máy khách có thể đã chấm dứt hoặc có thể chưa bắt đầu MCP. Sau 90 giây không có nhịp tim, phiên được coi là bị ngắt kết nối. |
| Tệp JSON không được nhập | Kiểm tra `version: 1`, `title` được yêu cầu, tên thuộc tính và nhân viên `Slug` hiện có. JSON không cho phép nhận xét hoặc dấu phẩy ở cuối. |
| Không thể xây dựng lại `agentboard.exe` | Windows không thể thay thế EXE đang chạy. Dừng máy chủ Ctrl+C và thử xây dựng lại. |
| Nhiệm vụ biến mất sau khi cập nhật | Kiểm tra xem `--db` không trỏ đến tệp khác và bạn đã đăng nhập với cùng một người dùng Windows. Cơ sở mặc định là `%AppData%\AgentBoard`. |

Nếu lỗi không được mô tả, hãy thu thập nội dung chính xác của thông báo, phiên bản `agentboard version` và Windows cũng như các bước để thử lại. Không xuất bản dữ liệu dự án riêng tư hoặc nội dung cơ sở dữ liệu trong một số báo mở.
