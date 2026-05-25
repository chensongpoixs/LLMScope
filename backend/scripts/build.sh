#!/bin/bash

# 构建脚本

echo "Building LLMScope Backend..."

# 清理旧的构建
rm -f llmscope-backend

# 构建
go build -o llmscope-backend main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Run with: ./llmscope-backend"
else
    echo "Build failed!"
    exit 1
fi
