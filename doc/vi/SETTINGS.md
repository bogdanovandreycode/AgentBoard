# Cài đặt

Mở **Cài đặt** trong menu bên của dự án đã chọn. Các thay đổi được lưu bằng nút **Save settings** và được lưu trữ cho dự án trong cơ sở dữ liệu AgentBoard.

## Ngôn ngữ và ngoại hình

Trường **Ngôn ngữ** sẽ mở danh sách thả xuống tìm kiếm. Theo mặc định, tùy chọn **Theo dõi hệ thống** được chọn, sử dụng ngôn ngữ trình duyệt. Các ngôn ngữ có sẵn từ ảnh chụp màn hình: tiếng Ả Rập, tiếng Bồ Đào Nha (Brazil), tiếng Trung giản thể, tiếng Séc, tiếng Đan Mạch, tiếng Hà Lan, tiếng Anh, tiếng Phần Lan, tiếng Pháp, tiếng Đức, tiếng Ý, tiếng Nhật, tiếng Hàn, tiếng Na Uy Bokmål, tiếng Ba Lan, tiếng Nga, tiếng Tây Ban Nha, tiếng Thụy Điển, tiếng Thổ Nhĩ Kỳ, tiếng Ukraina, tiếng Việt; ngoài ra còn có tiếng Belarus, tiếng Romania và tiếng Bungari. **Hệ thống theo dõi** sử dụng ngôn ngữ trình duyệt. Các chữ ký chính được dịch thủ công, các dòng UI còn lại được dịch tự động sơ bộ. Trước khi phát hành rộng rãi, nên hiệu đính bản dịch của người bản xứ; tên khách hàng, lệnh, trường JSON và dữ liệu người dùng vẫn không được dịch.

**Bảng màu**: Tối (chủ đề gốc), Sáng, Đen, Ubuntu và Windows. **Múi giờ** kiểm soát việc hiển thị ngày tháng; dữ liệu tiếp tục được lưu trữ trong UTC. **Múi giờ hệ thống** sử dụng cài đặt máy tính.

## Cột

Bốn giai đoạn `Features`, `In progress`, `Testing`, `Verification` được cố định theo thứ tự sau: chúng không thể đổi tên hoặc xóa. Các cột khác có thể được sắp xếp lại bằng các mũi tên có sẵn. Nhập tên và nhấp vào **Thêm cột** để tạo cột tùy chỉnh. Nó được thiết kế cho các nhiệm vụ của con người: AI nhìn thấy một nhiệm vụ như `Backlog` và không nhận nó thông qua MCP. Xóa một cột trong khi lưu sẽ chuyển nhiệm vụ của cột đó sang `Backlog` thông thường.

## Web và MCP

**Làm mới bảng** đặt tốc độ làm mới bảng và thẻ tính bằng giây (1–60). **Làm mới công nhân** cập nhật trạng thái của công nhân (2–120 giây). Đây là cuộc thăm dò giao diện web, không phải tần số kích hoạt AI. Công nhân không tự động bắt đầu.

Địa chỉ máy chủ web được đặt khi khởi động CLI, ví dụ `agentboard open --addr 127.0.0.1:7444`. Để thay đổi địa chỉ, máy chủ phải được khởi động lại. Mặc định là `127.0.0.1:7337`. MCP hoạt động thông qua lệnh cục bộ `agentboard mcp --project ... --worker ...` riêng biệt và độc lập với cổng web. Đối với cơ sở không chuẩn, hãy chỉ định cùng một `--db` trong tất cả các lệnh. Sao chép cấu hình của từng khách hàng từ thẻ nhân viên.
