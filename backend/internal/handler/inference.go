package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"llmscope-backend/internal/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Completion 文本补全，转发到 llama.cpp
func Completion(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := forwardToLlama(cfg, "/completion", body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		writeUpstreamResponse(c, resp)
	}
}

// Chat 聊天补全，转发到 llama.cpp /v1/chat/completions
func Chat(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req map[string]interface{}
		if err := json.Unmarshal(body, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stream, _ := req["stream"].(bool)

		resp, err := forwardToLlama(cfg, "/v1/chat/completions", body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   "无法连接到 llama.cpp",
				"message": err.Error(),
				"url":     cfg.LlamaCpp.BaseURL,
			})
			return
		}
		defer resp.Body.Close()

		if stream {
			proxyStreamResponse(c, resp)
			return
		}

		writeUpstreamResponse(c, resp)
	}
}

func forwardToLlama(cfg *config.Config, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, cfg.LlamaCpp.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return http.DefaultClient.Do(req)
}

func writeUpstreamResponse(c *gin.Context, resp *http.Response) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	c.Data(resp.StatusCode, contentType, respBody)
}

var hopByHopHeaders = map[string]bool{
	"connection": true, "keep-alive": true, "proxy-authenticate": true,
	"proxy-authorization": true, "te": true, "trailers": true,
	"transfer-encoding": true, "upgrade": true,
}

func proxyStreamResponse(c *gin.Context, resp *http.Response) {
	for key, values := range resp.Header {
		lower := strings.ToLower(key)
		// 保留本中间件 CORS 头，避免上游空值/重复导致浏览器误判
		if strings.HasPrefix(lower, "access-control-") || hopByHopHeaders[lower] {
			continue
		}
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	if c.Writer.Header().Get("Content-Type") == "" {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
	}
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Status(resp.StatusCode)

	flusher, canFlush := c.Writer.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, _ = c.Writer.Write(buf[:n])
			if canFlush {
				flusher.Flush()
			}
		}
		if err != nil {
			break
		}
	}
}
