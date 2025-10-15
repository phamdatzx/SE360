# =========================
# 1. Build stage
# =========================
FROM golang:alpine AS builder

# Cài đặt các công cụ cần thiết
RUN apk add --no-cache git

# Tạo thư mục làm việc trong container
WORKDIR /app

# Sao chép go.mod và go.sum trước (để cache dependency)
COPY go.mod go.sum ./

# Tải module
RUN go mod download

# Sao chép toàn bộ source code vào container
COPY . .

# Build binary (static binary để dễ chạy ở Alpine)
RUN go build -o server .

# =========================
# 2. Run stage
# =========================
FROM alpine:latest

# Cài chứng chỉ CA (để HTTPS hoạt động)
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Sao chép binary từ stage trước
COPY --from=builder /app/server .

# Expose port (ví dụ Gin chạy cổng 8080)
EXPOSE 8080

# Lệnh chạy ứng dụng
CMD ["./server"]
