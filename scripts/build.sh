#!/bin/bash
# -----------------------------
# Build & Push Docker Images to Docker Hub
# -----------------------------

# Dừng script nếu có lỗi
set -e

# Cấu hình cơ bản
DOCKER_USERNAME="megumikatou"          # ← đổi thành username Docker Hub của bạn
VERSION="v1.0.1"                     # ← hoặc dùng $(date +%Y%m%d) để tạo version theo ngày
SERVICES=("user-service" "driver-service" "trip-service")  # Danh sách các service

## Đăng nhập Docker Hub (yêu cầu bạn đã có token hoặc sẵn sàng nhập password)
#echo "🔐 Logging in to Docker Hub..."
#docker login -u "$DOCKER_USERNAME"

# Lặp qua từng service để build và push
for SERVICE in "${SERVICES[@]}"
do
  echo "🚧 Building image for $SERVICE ..."
  docker build -t "$DOCKER_USERNAME/$SERVICE:$VERSION" "services/$SERVICE"

  echo "🏷️ Tagging latest version ..."
  docker tag "$DOCKER_USERNAME/$SERVICE:$VERSION" "$DOCKER_USERNAME/$SERVICE:latest"

  echo "📤 Pushing $SERVICE to Docker Hub..."
  docker push "$DOCKER_USERNAME/$SERVICE:$VERSION"
  docker push "$DOCKER_USERNAME/$SERVICE:latest"

  echo "✅ Done: $SERVICE"
  echo "---------------------------------------"
done

echo "🎉 All images have been built and pushed successfully!"
