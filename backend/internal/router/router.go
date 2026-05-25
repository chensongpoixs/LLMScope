package router

import (
	"llmscope-backend/internal/config"
	"llmscope-backend/internal/handler"
	"llmscope-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS(cfg.CORS))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由
	api := r.Group("/api")
	{
		// llama.cpp 相关
		llama := api.Group("/llama")
		{
			llama.GET("/status", handler.GetLlamaCppStatus(cfg))
			llama.GET("/models", handler.GetLlamaCppModels(cfg))
			llama.Any("/proxy/*path", handler.ProxyToLlamaCpp(cfg))
		}

		// 模型相关
		models := api.Group("/models")
		{
			models.GET("", handler.GetModels(cfg))
			models.POST("/:id/load", handler.LoadModel(cfg))
			models.GET("/:id/structure", handler.GetModelStructure(cfg))
			models.GET("/:id/tensors/:name", handler.GetTensorInfo(cfg))
			models.GET("/:id/info", handler.GetModelInfo(cfg))
		}

		// 推理相关
		inference := api.Group("/inference")
		{
			inference.POST("/completion", handler.Completion(cfg))
			inference.POST("/chat", handler.Chat(cfg))
		}

		// Attention 相关
		attention := api.Group("/attention")
		{
			attention.GET("/:modelId/layer/:layer/head/:head", handler.GetAttention(cfg))
		}
	}

	// WebSocket
	r.GET("/ws", handler.WebSocketHandler(cfg))

	return r
}
