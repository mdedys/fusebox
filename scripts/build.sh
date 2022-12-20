#!/bin/sh
#
# Build all executables.

echo "Building all executables..." >&2
for d in cmd/*/ ; do
  GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o bin/ ./"$d"
done

