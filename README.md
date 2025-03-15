#### Bài 1: API Rate Limiting

Mô tả bài toán

Bạn cần thiết kế một hệ thống Rate Limiter giúp giới hạn số lượng request từ một user trong một khoảng thời gian nhất định, giống như hệ thống của các API lớn (ví dụ: Twitter API, GitHub API).

Yêu cầu:

1. Hỗ trợ các loại giới hạn:
    - Fixed Window Rate Limiting:
    - Sliding Window Rate Limiting: Linh hoạt hơn fixed window.
    - Token Bucket Algorithm: Điều chỉnh số lượng request linh động.
2. Hỗ trợ nhiều user cùng lúc mà không bị chậm.
3. Cần đảm bảo hiệu suất. Có thể sử dụng thêm tài nguyên bên ngoài như DB, queue...

#### Bài 2: Tìm Đường Trong Đồ Thị Lớn

Mô tả bài toán

Cho một đồ thị vô hướng có trọng số rất lớn, yêu cầu:
1. Tìm đường đi ngắn nhất giữa hai điểm bất kỳ.
2. Hỗ trợ tìm kiếm trong thời gian thực
3. Hệ thống cần có khả năng cập nhật dữ liệu đồ thị khi có thêm đỉnh hoặc cạnh mới.
4. Tối ưu giải thuật
