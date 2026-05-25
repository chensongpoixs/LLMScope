package handler

import (
	"llmscope-backend/internal/config"
	"llmscope-backend/internal/types"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAttention 获取注意力矩阵
func GetAttention(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelID := c.Param("modelId")
		layer, _ := strconv.Atoi(c.Param("layer"))
		head, _ := strconv.Atoi(c.Param("head"))

		// 模拟注意力矩阵数据
		seqLen := 20
		matrix := make([][]float64, seqLen)
		for i := 0; i < seqLen; i++ {
			matrix[i] = make([]float64, seqLen)
			for j := 0; j < seqLen; j++ {
				// 生成模拟的注意力分数
				if j <= i {
					matrix[i][j] = rand.Float64()
				} else {
					matrix[i][j] = 0
				}
			}
		}

		tokens := make([]string, seqLen)
		for i := 0; i < seqLen; i++ {
			tokens[i] = "token_" + strconv.Itoa(i)
		}

		attention := types.AttentionData{
			ModelID: modelID,
			Layer:   layer,
			Head:    head,
			Tokens:  tokens,
			Matrix:  matrix,
		}

		c.JSON(http.StatusOK, attention)
	}
}
