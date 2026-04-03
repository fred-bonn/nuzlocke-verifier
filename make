#!/usr/bin/env bash

OUTPUT_NAME="myapp"
BUILD_DIR="."
BIN_DIR="./bin"

if go build -o "$OUTPUT_NAME" "$BUILD_DIR"; then
    mkdir -p "$BIN_DIR"
    mv "$OUTPUT_NAME" "$BIN_DIR/"
else
    exit 1
fi