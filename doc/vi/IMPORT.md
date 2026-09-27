# Nhập tác vụ từ JSON

Tab **Nhập tác vụ** chấp nhận tệp JSON ở dạng mã hóa UTF-8. Hãy sử dụng nó nếu bạn đang chuyển nhiệm vụ từ một dịch vụ khác. Nhấp vào **Tải xuống mẫu JSON**: Mẫu bao gồm tên thuộc tính tùy chỉnh hiện tại và `Slug` của nhân viên có sẵn cho dự án này.

## Từng bước một

1. Tạo các công nhân cần thiết (**Workers**) và các thuộc tính (**Properties**) cần thiết trước khi nhập.
2. Tải xuống mẫu, mở nó trong trình soạn thảo văn bản, thay thế các ví dụ bằng các tác vụ của riêng bạn và lưu tệp bằng tiện ích mở rộng `.json`.
3. Trên tab **Nhập tác vụ**, chọn một tệp. Giao diện sẽ hiển thị số lượng nhiệm vụ.
4. Nhấp vào **Nhập nhiệm vụ**. Máy chủ sẽ kiểm tra toàn bộ tệp: nếu phát hiện thấy lỗi, sẽ không có một tác vụ nào từ đó được ghi. Hãy sửa thông báo lỗi và thử lại.

Tệp tối thiểu:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Ví dụ với công nhân, tài sản và sự phụ thuộc:

```json
{
  "version": 1,
  "tasks": [
    {
      "key": "design",
      "title": "Prepare the design",
      "description": "## Goal\nPrepare the home page mockup.",
      "state": "features",
      "priority": "high",
      "testing_mode": "hybrid",
      "assignee": { "type": "worker", "worker": "codex" },
      "ai_test_instructions": "Check the build.",
      "human_test_instructions": "Review the page in a browser.",
      "properties": { "Department": "Design" }
    },
    {
      "title": "Approve the design",
      "depends_on": ["design"],
      "assignee": { "type": "human" }
    }
  ]
}
```

`version` phải là `1`, mảng `tasks` - từ 1 đến 1000 nhiệm vụ. `title` là bắt buộc. Các giai đoạn hợp lệ: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Ưu tiên: `critical`, `high`, `medium`, `low`; chế độ kiểm tra: `ai`, `human`, `hybrid`. Chịu trách nhiệm: `unassigned`, `human` hoặc `worker` với `worker` (Slug hoặc ID) hiện có. `properties` sử dụng tên hoặc ID của các thuộc tính đã được tạo. `key` là duy nhất trong tệp; `depends_on` đề cập đến các phím như vậy. Máy chủ tự tạo ID nhiệm vụ.

Nhập lại cùng một tệp sẽ tạo ra sự cố mới, vì vậy hãy kiểm tra bảng trước khi nhấp lại. Khi được nhập, các tác vụ được ghi là do con người tạo ra; Một mục hệ thống xuất hiện trong lịch sử.
