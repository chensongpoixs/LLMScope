package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"llmscope-backend/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Completion 文本补全
func Completion(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 转发到 llama.cpp
		jsonData, _ := json.Marshal(req)
		resp, err := http.Post(
			cfg.LlamaCpp.BaseURL+"/completion",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		c.JSON(http.StatusOK, result)
	}
}

// Chat 聊天补全
func Chat(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 转发到 llama.cpp
		jsonData, _ := json.Marshal(req)
		resp, err := http.Post(
			cfg.LlamaCpp.BaseURL+"/v1/chat/completions",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		c.JSON(http.StatusOK, result)
	}
}
