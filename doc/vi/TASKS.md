# Làm việc với các nhiệm vụ

## Thêm nhiệm vụ

Trên tab **Bảng**, hãy nhấp vào **Nhiệm vụ mới**. Điền tiêu đề và mô tả. Phần mô tả và hướng dẫn kiểm tra hỗ trợ Markdown: đánh dấu văn bản và sử dụng thanh định dạng để in đậm, in nghiêng, liên kết và danh sách. Nhấp vào **Lưu**.

Lĩnh vực:

| Lĩnh vực | |
| --- | --- |
| Tiêu đề | Tên viết tắt của nhiệm vụ. |
| Mô tả | Những gì cần phải làm và làm thế nào để hiểu rằng công việc đã sẵn sàng. |
| Tiểu bang | Giai đoạn làm việc; một nhiệm vụ mới thường bắt đầu tại `Backlog`. |
| Ưu tiên | `critical`, `high`, `medium` hoặc `low`. |
| Chịu trách nhiệm | Một người, một công nhân cụ thể hoặc không có hẹn. |
| Chế độ thử nghiệm | `AI` - kiểm tra AI; `Human` - kiểm tra của con người; `Hybrid` - cả hai. |
| Phụ thuộc | Những công việc cần hoàn thành sớm hơn. |
| Hướng dẫn kiểm tra AI/Con người | Hướng dẫn cho người đánh giá thích hợp. |
| Thuộc tính tùy chỉnh | Các trường bổ sung được tạo trong tab **Thuộc tính**. |

Đối với tác vụ AI, trước tiên hãy tạo một nhân viên, chọn nhân viên đó trong Chịu trách nhiệm, sau đó chuyển thẻ sang `Features`. AI chỉ nhìn thấy các nhiệm vụ được giao trong bốn giai đoạn từ `Features` đến `Verification`.

## Di chuyển nhiệm vụ

Kéo thẻ giữa các cột. Bạn có thể chuyển đổi giữa chế độ xem tất cả các cột và cột rộng bằng cách cuộn ngang. `Backlog` và `Complete` do con người điều khiển. AI chỉ có thể nâng cao nhiệm vụ `Features → In progress → Testing → Verification` thông qua các công cụ MCP đặc biệt. Nhiệm vụ có thử nghiệm thủ công chưa hoàn thành sẽ không vượt qua quá trình xác minh AI.

## Thẻ nhiệm vụ

Nhấp vào thẻ để xem mô tả, chủ sở hữu, hướng dẫn kiểm tra và các tab về lịch sử, bài kiểm tra, hiện vật và chi phí AI. **Sửa tác vụ** thay đổi nội dung. Nhận xét của người đó được thêm vào câu chuyện tổng thể. Tìm kiếm trong các bộ lọc hàng đầu tìm kiếm theo tiêu đề, mô tả, ID và nhân viên; một nút riêng biệt sẽ mở ra một tìm kiếm lớn.

## Cột và thuộc tính

Trong tab **Cài đặt**, bạn có thể thay đổi thứ tự các cột có sẵn cho một người và thêm cột của riêng bạn. Bốn giai đoạn AI được cố định và tiến hành theo cùng một thứ tự. Cột người dùng là nơi dành cho các nhiệm vụ do một người trì hoãn: đối với AI, nhiệm vụ đó có trạng thái `Backlog`. Khi một cột bị xóa, các tác vụ của nó sẽ trở lại `Backlog` bình thường.

Trong tab **Thuộc tính** bạn có thể thêm các trường như văn bản, số, cờ, ngày, lựa chọn và URL. `Human only` ẩn trường khỏi AI; `Agent read` cho phép đọc, `Agent read/write` còn cho phép ghi thông qua các công cụ hỗ trợ. Điều này không thay đổi quyền của AI đối với các bước nhiệm vụ.
