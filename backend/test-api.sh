#!/bin/bash

# 测试 API 脚本

echo "Testing LLMScope Backend API..."

# 测试健康检查
echo -e "\n1. Health Check:"
curl -s http://localhost:8080/health | jq .

# 测试获取模型列表
echo -e "\n2. Get Models:"
curl -s http://localhost:8080/api/models | jq .

# 测试获取模型结构
echo -e "\n3. Get Model Structure:"
curl -s http://localhost:8080/api/models/llama-2-7b/structure | jq .

# 测试获取张量信息
echo -e "\n4. Get Tensor Info:"
curl -s http://localhost:8080/api/models/llama-2-7b/tensors/q_proj.weight | jq .

# 测试获取 Attention 数据
echo -e "\n5. Get Attention Data:"
curl -s http://localhost:8080/api/attention/llama-2-7b/layer/0/head/0 | jq .

echo -e "\nAll tests completed!"
