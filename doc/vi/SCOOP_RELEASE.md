# Chuẩn bị phát hành cho Scoop

## Những gì đã được tự động hóa

`scripts/package-scoop.ps1 -Version 0.2.0` thực thi `npm ci`, bản dựng giao diện người dùng, `go test ./...`, bản dựng Windows `agentboard.exe` với số phiên bản, tạo ZIP và tính toán SHA-256 của ZIP đó. Từ **cùng một** tệp, nó tạo ra `release/agentboard.json` với URL, hàm băm, CLI-shim, phím tắt, `checkver` và `autoupdate`. Cơ sở dữ liệu nằm ngoài thư mục cài đặt, vì vậy `persist` không cần thiết trong tệp kê khai.

`.github/workflows/release.yml` trên thẻ `vX.Y.Z` chạy cùng một gói trong trình chạy Windows, xây dựng trình cài đặt Inno Setup và đính kèm ZIP, trình cài đặt và tệp kê khai vào Bản phát hành GitHub. Việc chạy quy trình làm việc theo cách thủ công chỉ tạo ra một tạo phẩm để thử nghiệm mà không xuất bản bản phát hành.

## Ghi lệnh cho người bảo trì

1. Kiểm tra xem mã, tài liệu và số phiên bản đã sẵn sàng chưa. Xác định giấy phép của dự án: bảng kê khai hiện cho biết `Unknown` vì tệp giấy phép không được xác định trong kho lưu trữ. Nếu bạn chọn giấy phép, vui lòng thêm `LICENSE` và cập nhật `license` trong tập lệnh trước khi phát hành.
2. Chạy `./scripts/package-scoop.ps1 -Version X.Y.Z` cục bộ. Kiểm tra đầu ra `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` và `agentboard version` sau khi giải nén. Không thay đổi ZIP sau khi đã tính toán hàm băm.
3. Tạo và gửi thẻ `vX.Y.Z`. GitHub Actions sẽ xuất bản Bản phát hành có ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` và tệp kê khai. Kiểm tra tất cả ba tệp trên trang Phát hành và ZIP SHA-256 từ bảng kê khai.
4. Trên máy Windows sạch có Scoop, hãy chạy `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, sau đó là `agentboard version`, `agentboard init` và `agentboard open` trong thư mục kiểm tra.
5. Để có nguồn cấp dữ liệu cập nhật vĩnh viễn, hãy đặt `agentboard.json` đã tạo vào nhóm Scoop của riêng bạn hoặc cung cấp nó trong nhóm công khai phù hợp. Hãy kiểm tra lại `scoop update agentboard` sau lần phát hành tiếp theo. Lệnh cài đặt từ URL phù hợp cho người mới làm quen lần đầu, trong khi nhóm sẽ thuận tiện hơn cho việc cập nhật.

Không thay thế giá trị ngẫu nhiên `hash` theo cách thủ công: Scoop kiểm tra nội dung của ZIP đã tải xuống.

Để kiểm tra tệp kê khai, hãy tập trung vào [định dạng Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [tạo tệp kê khai](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) và [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
