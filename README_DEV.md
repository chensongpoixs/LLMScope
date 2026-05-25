# LLMScope 开发指南

## 项目结构

```
src/
├── api/              # API 接口
├── components/       # 可复用组件
├── views/           # 页面视图
├── stores/          # Pinia 状态管理
├── router/          # 路由配置
├── composables/     # 组合式函数
├── config/          # 配置管理
├── types/           # TypeScript 类型定义
└── utils/           # 工具函数

public/
└── config.js        # 后端配置文件（可在构建后修改）
```

## 后端配置

项目使用 `public/config.js` 配置后端地址，该文件在构建后可以直接修改，无需重新编译前端代码。

### 配置文件位置
- **开发环境**: `public/config.js`
- **生产环境**: `dist/config.js`

### 配置示例

```javascript
window.LLMSCOPE_CONFIG = {
  apiBaseUrl: 'http://localhost:8080',
  wsUrl: 'ws://localhost:8080/ws',
  timeout: 30000,
  maxRetries: 3
}
```

### 修改后端地址

**开发环境**：修改 `public/config.js`

**生产环境**：
1. 构建项目：`npm run build`
2. 修改 `dist/config.js` 中的地址
3. 部署 `dist` 目录

## 安装依赖

```bash
npm install
```

## 开发

```bash
npm run dev
```

访问 http://localhost:3000

## 构建

```bash
npm run build
```

构建产物在 `dist` 目录，包含 `config.js` 配置文件。

## 核心功能模块

### 1. 首页 (HomeView)
- 模型卡片列表
- 模型导入功能

### 2. 仪表盘 (DashboardView)
- 参数统计
- 模型结构预览
- 层级参数分布图

### 3. 结构显微镜 (MicroscopeView)
- 无限缩放画布
- 焦平面系统 (Level 1-4)
- 节点详情面板

### 4. Attention 显微镜 (AttentionView)
- Attention Head 网格
- 热力图可视化
- 层间对比

### 5. 激活显微镜 (ActivationView)
- Prompt 输入
- 实时激活可视化
- Token 时间轴

## 后端接口要求

需要 llama.cpp 提供以下接口：

- `GET /api/models` - 获取模型列表
- `POST /api/models/:id/load` - 加载模型
- `GET /api/models/:id/structure` - 获取模型结构
- `GET /api/models/:id/tensors/:name` - 获取张量信息
- `WebSocket ws://localhost:8080/ws` - 实时数据流

## WebSocket 消息格式

### Token Stream
```json
{
  "type": "token",
  "token": "Hello",
  "position": 12
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
  "matrix": []
}
```

## 待实现功能

- [ ] PixiJS 可视化引擎集成
- [ ] WebGL Shader 热力图
- [ ] D3-force 图布局
- [ ] WebWorker 数据处理
- [ ] 双模型对比模式
- [ ] PDF 导出功能
- [ ] 暗色模式
