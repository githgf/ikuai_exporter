#!/bin/sh

version=$(git describe --tags `git rev-list --tags --max-count=1`)
flags="-X main.buildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X main.version=$version"

export CGO_ENABLED=0 GOOS=linux

echo "start to build version:$version"

mkdir -p ./output/linux/

for i in "arm64" "amd64" ; do
  echo "building for $i..."
  GOARCH="$i" go build -ldflags "$flags" -o ./output/linux/$i/app main.go
  chmod +x ./output/linux/$i/app
done

image=ccr.ccs.tencentyun.com/imoe-tech/go-playground:ikuai-exporter-"$version"
official_img=jakes/ikuai-exporter:latest
official_img_versioned=jakes/ikuai-exporter:"$version"
echo "packaging docker multiplatform image: $image"
echo "packaging docker multiplatform image: $official_img"
echo "packaging docker multiplatform image: $official_img_versioned"

docker buildx build --push \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION="$version" \
  -t "$image" -t "$official_img" -t "$official_img_versioned" .

echo "finished: $image"

