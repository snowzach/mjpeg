#/bin/sh

echo "Building amd64"
docker build -t snowzach/mjproxy:amd64 .
docker push snowzach/mjproxy:amd64

echo "Building arm32v7"
docker buildx build --platform linux/arm/v7 -t snowzach/mjproxy:arm32v7 --push -f Dockerfile .

echo "Building arm64"
docker buildx build --platform linux/arm64 -t snowzach/mjproxy:arm64 --push -f Dockerfile .

echo "Creating latest manifest"
docker manifest push --purge snowzach/mjproxy:latest
docker manifest create snowzach/mjproxy:latest snowzach/mjproxy:amd64 snowzach/mjproxy:arm32v7 snowzach/mjproxy:arm64
docker manifest push --purge snowzach/mjproxy:latest
