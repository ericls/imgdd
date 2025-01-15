#!/usr/bin/env bash

(cd web_client && rm -rf dist/web && npm run build -- --no-cache);

targets=("linux_amd64" "linux_arm64" "darwin_arm64" "windows_amd64")
set -f
for target in ${targets[@]}
do
  IFS='_' read -r -a array <<< "$target"
  os="${array[0]}"
  arch="${array[1]}"
  echo building ${target}
  env GOOS=${os} GOARCH=${arch} \
  go build \
  -o dist/imgdd_${target} \
  -ldflags "-s -w
    -X 'github.com/ericls/imgdd/buildflag.Debug=false'
    -X 'github.com/ericls/imgdd/buildflag.Dev=false'
    -X github.com/ericls/imgdd/buildflag.VersionHash=`git rev-parse HEAD`
  " \
  .
done
