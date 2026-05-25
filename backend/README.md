# LLMScope Backend

基于 Go 语言的 LLMScope 后端服务，作为 llama.cpp 的代理层。

## 功能特性

- 模型管理 API
- 模型结构查询
- 张量信息获取
- Attention 可视化数据
- WebSocket 实时推理流
- 激活数据流式传输

## 项目结构

```
backend/
├── main.go                    # 入口文件
├── go.mod                     # Go 模块定义
├── internal/
│   ├── config/               # 配置管理
│   ├── handler/              # HTTP 处理器
│   ├── middleware/           # 中间件
│   ├── router/               # 路由配置
│   └── types/                # 类型定义
├── scripts/                  # 辅助脚本
└── .env.example              # 环境变量示例
```

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod download
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件配置 llama.cpp 地址
```

### 3. 运行服务

**Windows:**
```powershell
go run main.go
```

**Linux/Mac:**
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 4. 测试 API

**Windows:**
```powershell
.\test-api.ps1
```

**Linux/Mac:**
```bash
chmod +x test-api.sh
./test-api.sh
```

### 5. 构建

```bash
go build -o llmscope-backend main.go
```

**Windows:**
```powershell
.\llmscope-backend.exe
```

**Linux/Mac:**
```bash
./llmscope-backend
```

## API 接口

### 健康检查
- `GET /health` - 服务健康检查

### 模型管理

- `GET /api/models` - 获取模型列表
- `POST /api/models/:id/load` - 加载模型
- `GET /api/models/:id/structure` - 获取模型结构
- `GET /api/models/:id/tensors/:name` - 获取张量信息
- `GET /api/models/:id/info` - 获取模型详细信息

### 推理

- `POST /api/inference/completion` - 文本补全
- `POST /api/inference/chat` - 聊天补全

### Attention

- `GET /api/attention/:modelId/layer/:layer/head/:head` - 获取注意力矩阵

### WebSocket

- `ws://localhost:8080/ws` - 实时数据流

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| LLAMA_CPP_HOST | llama.cpp 主机 | localhost |
| LLAMA_CPP_PORT | llama.cpp 端口 | 8081 |

## 开发说明

### 目录说明

- `internal/config/` - 配置加载和管理
- `internal/handler/` - API 处理器实现
  - `model.go` - 模型相关接口
  - `attention.go` - Attention 数据接口
  - `inference.go` - 推理接口
  - `websocket.go` - WebSocket 实时流
- `internal/middleware/` - 中间件（CORS等）
- `internal/router/` - 路由配置
- `internal/types/` - 数据类型定义

### 添加新接口

1. 在 `internal/types/types.go` 定义数据结构
2. 在 `internal/handler/` 创建处理器
3. 在 `internal/router/router.go` 注册路由

## WebSocket 消息格式

### Token Stream
```json
{
  "type": "token",
  "token": "Hello",
  "position": 0
}
```

### Activation Stream
```json
{
  "type": "activation",
  "layer": 15,
  "mean": 0.23,
  "max": 1.2
}
```

### Attention Stream
```json
{
  "type": "attention",
  "layer": 10,
  "head": 5,
  "matrix": [[...]]
}
```

## 故障排除

### 依赖问题

如果遇到依赖错误，运行：
```bash
rm go.sum
go mod tidy
```

### 端口占用

如果 8080 端口被占用，修改 `.env` 文件中的 `SERVER_PORT`

### CORS 问题

CORS 已默认配置为允许所有来源，如需限制请修改 `internal/config/config.go`
