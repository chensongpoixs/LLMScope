package handler

import (
	"encoding/json"
	"llmscope-backend/internal/config"
	"llmscope-backend/internal/types"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler WebSocket 处理器
func WebSocketHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		log.Println("WebSocket client connected")

		// 处理客户端消息
		go handleClientMessages(conn)

		// 模拟推理过程，发送实时数据
		simulateInference(conn)
	}
}

func handleClientMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		log.Printf("Received message: %v", msg)

		// 根据消息类型处理
		msgType, ok := msg["type"].(string)
		if !ok {
			continue
		}

		switch msgType {
		case "start_inference":
			// 开始推理
			log.Println("Starting inference...")
		case "stop_inference":
			// 停止推理
			log.Println("Stopping inference...")
		}
	}
}

func simulateInference(conn *websocket.Conn) {
	tokens := []string{"Hello", "world", "this", "is", "a", "test", "of", "LLM", "inference"}

	for i, token := range tokens {
		// 发送 Token
		tokenMsg := types.WSMessage{
			Type:     "token",
			Token:    token,
			Position: i,
		}
		if err := conn.WriteJSON(tokenMsg); err != nil {
			log.Printf("Write error: %v", err)
			return
		}

		// 发送激活数据
		activationMsg := types.WSMessage{
			Type:  "activation",
			Layer: rand.Intn(32),
			Mean:  rand.Float64(),
			Max:   rand.Float64() * 2,
		}
		if err := conn.WriteJSON(activationMsg); err != nil {
			log.Printf("Write error: %v", err)
			return
		}

		// 发送注意力数据
		if i%3 == 0 {
			attentionMsg := types.WSMessage{
				Type:   "attention",
				Layer:  rand.Intn(32),
				Head:   rand.Intn(32),
				Matrix: generateRandomMatrix(10, 10),
			}
			if err := conn.WriteJSON(attentionMsg); err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Println("Inference simulation completed")
}

func generateRandomMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j] = rand.Float64()
		}
	}
	return matrix
}
