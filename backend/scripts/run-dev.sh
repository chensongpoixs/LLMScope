#!/bin/bash

# 开发环境启动脚本

# 加载环境变量
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# 启动服务
echo "Starting LLMScope Backend..."
go run main.go
