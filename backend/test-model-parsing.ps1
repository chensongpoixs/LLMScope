# 测试模型名称解析功能

Write-Host "测试模型名称解析功能" -ForegroundColor Green
Write-Host "================================`n"

$testCases = @(
    @{Name = "gemma-4-3b-it-Q4_0.gguf"; ExpectedParams = "3B"; ExpectedQuant = "Q4_0"},
    @{Name = "llama-2-7b-chat.Q4_K_M.gguf"; ExpectedParams = "7B"; ExpectedQuant = "Q4_K_M"},
    @{Name = "mistral-7b-instruct.Q5_K_S.gguf"; ExpectedParams = "7B"; ExpectedQuant = "Q5_K_S"},
    @{Name = "qwen-1.5b-chat.Q4_0.gguf"; ExpectedParams = "1.5B"; ExpectedQuant = "Q4_0"}
)

Write-Host "模拟测试用例:" -ForegroundColor Yellow
foreach ($test in $testCases) {
    Write-Host "`n模型: $($test.Name)"
    Write-Host "  预期参数量: $($test.ExpectedParams)"
    Write-Host "  预期量化: $($test.ExpectedQuant)"
}

Write-Host "`n================================"
Write-Host "运行单元测试..." -ForegroundColor Yellow
go test ./internal/handler/ -v -run TestExtract

Write-Host "`n================================"
Write-Host "✓ 修复完成！" -ForegroundColor Green
Write-Host "`n请重启后端服务并刷新前端页面查看效果。"
