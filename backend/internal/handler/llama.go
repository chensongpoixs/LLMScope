package handler

import (
	"encoding/json"
	"io"
	"llmscope-backend/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetLlamaCppStatus 获取 llama.cpp 服务状态
func GetLlamaCppStatus(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试连接 llama.cpp
		resp, err := http.Get(cfg.LlamaCpp.BaseURL + "/health")
		
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"connected": false,
				"url":       cfg.LlamaCpp.BaseURL,
				"error":     err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		c.JSON(http.StatusOK, gin.H{
			"connected": true,
			"url":       cfg.LlamaCpp.BaseURL,
			"status":    resp.StatusCode,
		})
	}
}

// GetLlamaCppModels 从 llama.cpp 获取模型信息
func GetLlamaCppModels(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := http.Get(cfg.LlamaCpp.BaseURL + "/v1/models")
		
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "无法连接到 llama.cpp",
				"message": err.Error(),
				"url":     cfg.LlamaCpp.BaseURL,
			})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// ProxyToLlamaCpp 代理请求到 llama.cpp
func ProxyToLlamaCpp(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("path")
		targetURL := cfg.LlamaCpp.BaseURL + path

		// 创建新请求
		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 复制请求头
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "无法连接到 llama.cpp",
				"url":   targetURL,
			})
			return
		}
		defer resp.Body.Close()

		// 复制响应头
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// 复制响应体
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}
