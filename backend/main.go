package main

import (
	"log"
	"llmscope-backend/internal/config"
	"llmscope-backend/internal/router"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化路由
	r := router.Setup(cfg)

	// 启动服务器
	log.Printf("LLMScope Backend starting on %s", cfg.Server.Address)
	if err := r.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
