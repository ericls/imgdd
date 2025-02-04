#!/usr/bin/env bash
git_rev=$(git describe --tags --always --dirty)
git_hash=$(git rev-parse HEAD)
git_tag_name=$(git describe --tags --always)
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg GIT_REV=$git_rev \
  --build-arg GIT_HASH=$git_hash \
  -t imgdd:$git_tag_name \
  -f Dockerfile \
  .
