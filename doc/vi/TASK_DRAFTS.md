# Quản lý dự án và dự thảo nhiệm vụ

Trang **Tất cả dự án** hiển thị thư mục đã đăng ký, số lượng bài viết đã xuất bản
nhiệm vụ và hoạt động mới nhất. Hoạt động tính đến những thay đổi trong nhiệm vụ, lịch sử và
hồ sơ chi phí. Nút bút chì cho phép bạn đổi tên dự án hoặc chọn
thư mục hiện có trong hộp thoại hệ thống. Thay đổi đường dẫn sẽ cập nhật đăng ký và điểm đánh dấu;
ID dự án, nhân viên, nhiệm vụ và cài đặt được lưu, tệp không bị di chuyển.
Không thể chọn các thư mục đã được đăng ký cho dự án khác.

## Cài đặt chung

Mở **Cài đặt chung** từ trang hoặc tab chọn dự án
**Cài đặt chung** trong cài đặt của bất kỳ bảng nào. Các thông số này là chung cho tất cả các dự án
một cơ sở AgentBoard:

- Tính năng tự động lưu bản nháp theo bộ hẹn giờ được bật theo mặc định, cứ sau **10 giây**.
- **Lưu khi rời khỏi trường** được bật theo mặc định và có thể được cấu hình riêng.
- **Kiểm tra bối cảnh dự án MCP** được bật và tắt trên toàn cầu. Khu vực của họ
  các hành động và cờ CLI riêng biệt được mô tả trong phần [Workers và MCP](WORKERS_MCP.md).

## Vòng đời dự thảo

Hộp thoại Tác vụ mới lưu công việc chưa hoàn thành dưới dạng bản nháp. Dòng trạng thái
phân biệt giữa các thay đổi chưa được lưu, đang lưu, bản nháp đã lưu và lỗi
bảo tồn. Bản dự thảo không xuất hiện trên bảng và không được cung cấp cho nhân viên MCP cho đến khi
được xuất bản.

Trước khi xuất bản, hãy sử dụng nút **Tạo tác vụ** để điền vào **Tiêu đề, Trạng thái,
Chịu trách nhiệm và Mô tả**. Trước đó, **Lưu bản nháp** lưu bản nháp mà không tạo
nhiệm vụ. Theo mặc định, biểu mẫu trống sẽ bị loại bỏ khi đóng. Nếu xóa
bản nháp hiện có sang trạng thái trống, nó sẽ bị xóa vào lần lưu tiếp theo.
Đóng một hộp thoại không trống cũng sẽ lưu nó.

Nút **Bản nháp** trên bảng cho phép bạn tiếp tục hoặc xóa các bản nháp đã lưu.
Hỗ trợ nhiều bản nháp. Sau khi khởi động lại, mở
dự thảo. Bản nháp được lưu trữ trong SQLite; Với mỗi thay đổi, thay đổi cục bộ sẽ được cập nhật
sao chép trong trình duyệt, để các chỉnh sửa chưa được lưu vẫn tồn tại khi khởi động lại hoặc
lỗi máy chủ tạm thời. Nếu bộ nhớ trình duyệt không có sẵn, một hộp thoại sẽ cảnh báo
rằng bản nháp phải được lưu trên máy chủ trước khi đóng.

Viết bản nháp sử dụng các bản sửa đổi để ngăn một cửa sổ ghi đè công việc
cái khác. Phản hồi lưu bị mất là nhất quán nếu đã có một phản hồi trên máy chủ
dữ liệu giống hệt nhau. Nếu nội dung khác nhau, xung đột sẽ được hiển thị và
**Khôi phục dưới dạng bản nháp mới** lưu nội dung cục bộ một cách riêng biệt. Xuất bản
nguyên tử và bình thường: xuất bản lại bản nháp tương tự sau khi bị mất
phản hồi trả về cùng một tác vụ thay vì tạo một bản sao.
