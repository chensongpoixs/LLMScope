package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"llmscope-backend/internal/config"
	"llmscope-backend/internal/types"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// extractParamsFromModelName 从模型名称中提取参数量
func extractParamsFromModelName(name string) string {
	// 匹配常见的参数量格式: 7B, 13B, 70B, 3B, 1.5B 等
	patterns := []string{
		`(\d+\.?\d*)[Bb]`,     // 7B, 13B, 1.5B
		`(\d+\.?\d*)[Mm]`,     // 7M (百万参数)
		`-(\d+\.?\d*)[Bb]-`,   // -7B-
		`_(\d+\.?\d*)[Bb]_`,   // _7B_
		`(\d+\.?\d*)billion`,  // 7billion
		`(\d+\.?\d*)million`,  // 7million
	}

	nameLower := strings.ToLower(name)

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(nameLower); len(matches) > 1 {
			value := matches[1]
			// 判断单位
			if strings.Contains(pattern, "[Mm]") || strings.Contains(pattern, "million") {
				return value + "M"
			}
			return value + "B"
		}
	}

	return "Unknown"
}

// extractQuantizationFromModelName 从模型名称中提取量化信息
func extractQuantizationFromModelName(name string) string {
	// 匹配常见的量化格式: Q4_0, Q4_K_M, Q5_K_S, Q8_0, F16, F32 等
	patterns := []string{
		`[Qq](\d+)_([0KkMmSs_]+)`, // Q4_K_M, Q5_K_S, Q4_0
		`[Ff](\d+)`,                // F16, F32
		`int(\d+)`,                 // int4, int8
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(name); len(matches) > 0 {
			return strings.ToUpper(matches[0])
		}
	}

	return "Unknown"
}

func toInt64(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int64:
		return n, true
	case int:
		return int64(n), true
	default:
		return 0, false
	}
}

// formatParamsFromMeta 从 llama.cpp meta.n_params 格式化参数量
func formatParamsFromMeta(nParams int64) string {
	if nParams <= 0 {
		return "Unknown"
	}
	billions := float64(nParams) / 1e9
	if billions >= 1 {
		if billions >= 10 {
			return fmt.Sprintf("%.0fB", math.Round(billions))
		}
		return fmt.Sprintf("%.1fB", billions)
	}
	millions := float64(nParams) / 1e6
	if millions >= 1 {
		return fmt.Sprintf("%.0fM", math.Round(millions))
	}
	return "Unknown"
}

// contextLengthK 将 n_ctx（token 数）转为前端展示的 K 单位
func contextLengthK(nCtx int64) int {
	if nCtx <= 0 {
		return 0
	}
	k := int(nCtx / 1024)
	if k < 1 {
		return 1
	}
	return k
}

func formatFileSize(bytes int64) string {
	if bytes <= 0 {
		return "Unknown"
	}
	const unit = 1024.0
	b := float64(bytes)
	if b >= unit*unit*unit {
		return fmt.Sprintf("%.1fGB", b/(unit*unit*unit))
	}
	if b >= unit*unit {
		return fmt.Sprintf("%.1fMB", b/(unit*unit))
	}
	return fmt.Sprintf("%.0fKB", b/unit)
}

func applyMetaToModel(model *types.ModelInfo, meta map[string]interface{}) {
	if nCtx, ok := toInt64(meta["n_ctx"]); ok && nCtx > 0 {
		model.ContextLength = contextLengthK(nCtx)
	}
	if nParams, ok := toInt64(meta["n_params"]); ok && nParams > 0 {
		model.Params = formatParamsFromMeta(nParams)
	}
	if size, ok := toInt64(meta["size"]); ok && size > 0 {
		model.FileSize = formatFileSize(size)
	}
}

// GetModels 获取模型列表
func GetModels(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从 llama.cpp 获取模型信息
		resp, err := http.Get(cfg.LlamaCpp.BaseURL + "/v1/models")

		var models []types.ModelInfo

		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			var llamaResponse map[string]interface{}
			if json.Unmarshal(body, &llamaResponse) == nil {
				// 解析 llama.cpp 返回的模型信息
				if data, ok := llamaResponse["data"].([]interface{}); ok && len(data) > 0 {
					for _, item := range data {
						if modelData, ok := item.(map[string]interface{}); ok {
							modelID := ""
							modelName := ""

							if id, ok := modelData["id"].(string); ok {
								modelID = id
								modelName = id
							}

							// 从文件名解析（meta 缺失时的回退）
							params := extractParamsFromModelName(modelName)
							quantization := extractQuantizationFromModelName(modelName)

							architecture := "LLaMA"
							nameLower := strings.ToLower(modelName)
							if strings.Contains(nameLower, "gemma") {
								architecture = "Gemma"
							} else if strings.Contains(nameLower, "mistral") {
								architecture = "Mistral"
							} else if strings.Contains(nameLower, "qwen") {
								architecture = "Qwen"
							} else if strings.Contains(nameLower, "llama") {
								architecture = "LLaMA"
							} else if strings.Contains(nameLower, "phi") {
								architecture = "Phi"
							}

							info := types.ModelInfo{
								ID:            modelID,
								Name:          modelName,
								Params:        params,
								Architecture:  architecture,
								ContextLength: 0,
								Quantization:  quantization,
								FileSize:      "Unknown",
								Loaded:        true,
							}

							// 优先使用 llama.cpp meta（n_ctx / n_params / size）
							if meta, ok := modelData["meta"].(map[string]interface{}); ok {
								applyMetaToModel(&info, meta)
							}

							// meta 未提供 context 时，尝试从 n_ctx_train 等字段
							if info.ContextLength == 0 {
								if meta, ok := modelData["meta"].(map[string]interface{}); ok {
									if nTrain, ok := toInt64(meta["n_ctx_train"]); ok && nTrain > 0 {
										info.ContextLength = contextLengthK(nTrain)
									}
								}
							}
							if info.ContextLength == 0 {
								info.ContextLength = 4 // 最后回退
							}

							models = append(models, info)
						}
					}
				}
			}
		}

		// 如果无法从 llama.cpp 获取，返回模拟数据
		if len(models) == 0 {
			models = []types.ModelInfo{
				{
					ID:            "llama-2-7b",
					Name:          "LLaMA-2-7B (示例)",
					Params:        "7B",
					Architecture:  "LLaMA",
					ContextLength: 4,
					Quantization:  "Q4_K_M",
					FileSize:      "3.8GB",
					Loaded:        false,
				},
			}
		}

		c.JSON(http.StatusOK, models)
	}
}

// LoadModel 加载模型
func LoadModel(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelID := c.Param("id")

		// 调用 llama.cpp 加载模型
		resp, err := http.Post(cfg.LlamaCpp.BaseURL+"/load", "application/json", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"modelId": modelID,
			"message": "Model loaded successfully",
		})
	}
}

// GetModelStructure 获取模型结构
func GetModelStructure(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelID := c.Param("id")

		// 模拟模型结构数据
		structure := types.ModelStructure{
			Embedding: types.LayerNode{
				Name:   "embedding",
				Type:   "embedding",
				Params: 32000 * 4096,
			},
			Transformers: make([]types.TransformerLayer, 32),
			LMHead: types.LayerNode{
				Name:   "lm_head",
				Type:   "linear",
				Params: 4096 * 32000,
			},
		}

		// 填充 Transformer 层
		for i := 0; i < 32; i++ {
			structure.Transformers[i] = types.TransformerLayer{
				ID: i,
				Attention: types.AttentionModule{
					QProj: types.TensorNode{Name: "q_proj", Shape: []int{4096, 4096}, Dtype: "float16"},
					KProj: types.TensorNode{Name: "k_proj", Shape: []int{4096, 4096}, Dtype: "float16"},
					VProj: types.TensorNode{Name: "v_proj", Shape: []int{4096, 4096}, Dtype: "float16"},
					OProj: types.TensorNode{Name: "o_proj", Shape: []int{4096, 4096}, Dtype: "float16"},
					Heads: 32,
				},
				FFN: types.FFNModule{
					UpProj:   types.TensorNode{Name: "up_proj", Shape: []int{4096, 11008}, Dtype: "float16"},
					DownProj: types.TensorNode{Name: "down_proj", Shape: []int{11008, 4096}, Dtype: "float16"},
					GateProj: &types.TensorNode{Name: "gate_proj", Shape: []int{4096, 11008}, Dtype: "float16"},
				},
				Norm: types.NormModule{
					Weight: types.TensorNode{Name: "norm_weight", Shape: []int{4096}, Dtype: "float16"},
				},
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"modelId":   modelID,
			"structure": structure,
		})
	}
}

// GetTensorInfo 获取张量信息
func GetTensorInfo(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelID := c.Param("id")
		tensorName := c.Param("name")

		// 模拟张量信息
		tensorInfo := types.TensorInfo{
			Name:         tensorName,
			Shape:        []int{4096, 4096},
			Dtype:        "float16",
			Params:       4096 * 4096,
			Device:       "GPU",
			Quantization: "Q4_K_M",
			Min:          -0.5,
			Max:          0.5,
			Mean:         0.001,
			Std:          0.12,
		}

		c.JSON(http.StatusOK, gin.H{
			"modelId": modelID,
			"tensor":  tensorInfo,
		})
	}
}

// GetModelInfo 获取模型详细信息
func GetModelInfo(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelID := c.Param("id")

		// 尝试从 llama.cpp 获取模型信息
		resp, err := http.Get(cfg.LlamaCpp.BaseURL + "/v1/models")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

		c.JSON(http.StatusOK, gin.H{
			"modelId": modelID,
			"info":    result,
		})
	}
}
