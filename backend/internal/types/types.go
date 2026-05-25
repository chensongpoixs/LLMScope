package types

// ModelInfo 模型信息
type ModelInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Params        string `json:"params"`
	Architecture  string `json:"architecture"`
	ContextLength int    `json:"contextLength"`
	Quantization  string `json:"quantization"`
	FileSize      string `json:"fileSize"`
	Loaded        bool   `json:"loaded"`
}

// ModelStructure 模型结构
type ModelStructure struct {
	Embedding    LayerNode          `json:"embedding"`
	Transformers []TransformerLayer `json:"transformers"`
	LMHead       LayerNode          `json:"lmHead"`
}

// TransformerLayer Transformer 层
type TransformerLayer struct {
	ID        int              `json:"id"`
	Attention AttentionModule  `json:"attention"`
	FFN       FFNModule        `json:"ffn"`
	Norm      NormModule       `json:"norm"`
}

// LayerNode 层节点
type LayerNode struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Params int    `json:"params"`
}

// AttentionModule 注意力模块
type AttentionModule struct {
	QProj TensorNode `json:"qProj"`
	KProj TensorNode `json:"kProj"`
	VProj TensorNode `json:"vProj"`
	OProj TensorNode `json:"oProj"`
	Heads int        `json:"heads"`
}

// FFNModule 前馈网络模块
type FFNModule struct {
	UpProj   TensorNode  `json:"upProj"`
	DownProj TensorNode  `json:"downProj"`
	GateProj *TensorNode `json:"gateProj,omitempty"`
}

// NormModule 归一化模块
type NormModule struct {
	Weight TensorNode  `json:"weight"`
	Bias   *TensorNode `json:"bias,omitempty"`
}

// TensorNode 张量节点
type TensorNode struct {
	Name  string `json:"name"`
	Shape []int  `json:"shape"`
	Dtype string `json:"dtype"`
}

// TensorInfo 张量详细信息
type TensorInfo struct {
	Name         string  `json:"name"`
	Shape        []int   `json:"shape"`
	Dtype        string  `json:"dtype"`
	Params       int     `json:"params"`
	Device       string  `json:"device"`
	Quantization string  `json:"quantization,omitempty"`
	Min          float64 `json:"min,omitempty"`
	Max          float64 `json:"max,omitempty"`
	Mean         float64 `json:"mean,omitempty"`
	Std          float64 `json:"std,omitempty"`
}

// AttentionData 注意力数据
type AttentionData struct {
	ModelID string      `json:"modelId"`
	Layer   int         `json:"layer"`
	Head    int         `json:"head"`
	Tokens  []string    `json:"tokens"`
	Matrix  [][]float64 `json:"matrix"`
}

// WSMessage WebSocket 消息
type WSMessage struct {
	Type     string      `json:"type"`
	Token    string      `json:"token,omitempty"`
	Position int         `json:"position,omitempty"`
	Layer    int         `json:"layer,omitempty"`
	Head     int         `json:"head,omitempty"`
	Mean     float64     `json:"mean,omitempty"`
	Max      float64     `json:"max,omitempty"`
	Matrix   [][]float64 `json:"matrix,omitempty"`
}
