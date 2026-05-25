# PowerShell 测试脚本

Write-Host "Testing LLMScope Backend API..." -ForegroundColor Green

# 测试健康检查
Write-Host "`n1. Health Check:" -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get | ConvertTo-Json

# 测试获取模型列表
Write-Host "`n2. Get Models:" -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8080/api/models" -Method Get | ConvertTo-Json

# 测试获取模型结构
Write-Host "`n3. Get Model Structure:" -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8080/api/models/llama-2-7b/structure" -Method Get | ConvertTo-Json -Depth 5

# 测试获取张量信息
Write-Host "`n4. Get Tensor Info:" -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8080/api/models/llama-2-7b/tensors/q_proj.weight" -Method Get | ConvertTo-Json

# 测试获取 Attention 数据
Write-Host "`n5. Get Attention Data:" -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8080/api/attention/llama-2-7b/layer/0/head/0" -Method Get | ConvertTo-Json -Depth 3

Write-Host "`nAll tests completed!" -ForegroundColor Green
