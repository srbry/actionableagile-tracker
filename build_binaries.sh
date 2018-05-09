#!/bin/bash -e

go_dists=$(go tool dist list)

mkdir -p releases
for dist in $go_dists; do
  os=$(awk -F/ '{print $1}' <<< "$dist")
  arch=$(awk -F/ '{print $2}' <<< "$dist")
  [ "${os}" = "android" ] || [ "${os}" = "nacl" ] && continue
  echo "Building $dist"
  GOOS="${os}" ARCH="${arch}" go build -o "releases/${os}_${arch}_actionableagile-tracker" main.go
done
