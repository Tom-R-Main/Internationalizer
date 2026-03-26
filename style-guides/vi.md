# Vietnamese (vi) — Translation Style Guide

## Language Profile
- Locale: `vi`
- Script: Latin, LTR
- CLDR plural forms: other
- Text expansion vs English: 15-20%

## Tone & Formality
Giọng văn cần chuyên nghiệp, lịch sự, rõ ràng và thân thiện. Trong giao diện người dùng (UI), ưu tiên cách diễn đạt ngắn gọn, trực tiếp để tiết kiệm không gian. 

Hạn chế tối đa việc sử dụng đại từ nhân xưng. Nếu bắt buộc phải xưng hô với người dùng trong các câu thông báo dài hoặc email, hãy dùng "Bạn" (viết hoa chữ B) và xưng "Chúng tôi" (đại diện cho ứng dụng/hệ thống). Tuyệt đối không dùng tiếng lóng, từ ngữ địa phương hoặc cách nói quá suồng sã.

## Grammar
- **Nút bấm (Buttons) & Mệnh lệnh:** Không dùng các từ như "Hãy" hoặc "Vui lòng" trên các nút bấm. Chỉ dùng động từ gốc để giữ UI ngắn gọn. (VD: Dùng "Lưu" thay vì "Hãy lưu", "Xóa" thay vì "Vui lòng xóa").
- **Trật tự từ trong cụm danh từ:** Tính từ và từ bổ nghĩa luôn đứng sau danh từ chính. LLM thường dịch ngược theo cấu trúc tiếng Anh. (VD: "New folder" phải dịch là "Thư mục mới", không phải "Mới thư mục").
- **Câu bị động (Passive Voice):** Tiếng Việt ít dùng câu bị động hơn tiếng Anh. Hãy chuyển sang chủ động nếu có thể. Nếu bắt buộc dùng bị động, phải phân biệt rõ: dùng "được" cho ý nghĩa tích cực/trung tính và "bị" cho ý nghĩa tiêu cực. (VD: "File was saved" -> "Tệp đã được lưu"; "Account was banned" -> "Tài khoản đã bị khóa").
- **Viết hoa (Capitalization):** Sử dụng "Sentence case" (chỉ viết hoa chữ cái đầu tiên của câu hoặc cụm từ) cho hầu hết các thành phần UI (tiêu đề, nút bấm, nhãn). Không viết hoa từng từ (Title Case) như tiếng Anh trừ khi đó là danh từ riêng. (VD: "Cài đặt tài khoản" thay vì "Cài Đặt Tài Khoản").

## Pluralization
Tiếng Việt chỉ có một dạng số nhiều theo chuẩn CLDR là `other`. Danh từ trong tiếng Việt không thay đổi hình thức (không thêm hậu tố) dù số lượng là 0, 1 hay nhiều hơn. 

Chỉ cần đặt số đếm trực tiếp trước danh từ.
- VD: 0 tệp, 1 tệp, 2 tệp, 100 tệp. (Không dịch là "1 tệp, 2 các tệp").
Nếu không có số đếm cụ thể, có thể dùng các từ chỉ số lượng như "các" hoặc "những" (VD: "Xóa các tệp đã chọn").

## Punctuation & Typography
- **Dấu ngoặc kép:** Sử dụng ngoặc kép kép ("...") thay vì ngoặc đơn ('...') cho các trích dẫn hoặc làm nổi bật từ.
- **Số thập phân & Hàng nghìn:** Dùng dấu phẩy (`,`) cho số thập phân và dấu chấm (`.`) để phân cách hàng nghìn. (VD: `1.234.567,89`).
- **Định dạng ngày tháng:** Luôn sử dụng định dạng DD/MM/YYYY (VD: `31/12/2023`).
- **Định dạng thời gian:** Ưu tiên sử dụng định dạng 24 giờ (VD: `14:30`). Nếu hệ thống yêu cầu định dạng 12 giờ, hãy dùng "SA" (Sáng - AM) và "CH" (Chiều - PM).

## Terminology

| English | Vietnamese | Notes |
|---------|------------|-------|
| Save | Lưu | Dùng cho nút lưu dữ liệu hoặc thay đổi. |
| Cancel | Hủy | Dùng để hủy bỏ một thao tác hoặc đóng hộp thoại. |
| Delete | Xóa | Dùng khi xóa vĩnh viễn một mục hoặc dữ liệu. |
| Settings | Cài đặt | Dùng cho menu cấu hình hệ thống hoặc ứng dụng. |
| Search | Tìm kiếm | Dùng cho thanh tìm kiếm hoặc nút thực thi tìm kiếm. |
| Error | Lỗi | Dùng trong các thông báo lỗi hệ thống. |
| Loading | Đang tải | Dùng cho trạng thái chờ hệ thống xử lý dữ liệu. |
| Dashboard | Bảng điều khiển | Trang tổng quan chính của ứng dụng. |
| Notifications | Thông báo | Dùng cho trung tâm cảnh báo hoặc tin báo mới. |
| Sign in | Đăng nhập | Hành động truy cập vào tài khoản người dùng. |
| Sign out | Đăng xuất | Hành động thoát khỏi tài khoản người dùng. |
| Submit | Gửi | Dùng cho nút xác nhận biểu mẫu (form). |
| Profile | Hồ sơ | Trang chứa thông tin cá nhân của người dùng. |
| Help | Trợ giúp | Liên kết đến tài liệu hướng dẫn hoặc hỗ trợ. |
| Close | Đóng | Dùng để đóng cửa sổ, tab hoặc hộp thoại. |
