# LLMScope

> Large Language Model Microscope
> 大模型显微镜

一个基于 `llama.cpp` 的大模型内部结构可视化与推理观测平台。

---

# 一、项目定位

LLMScope 是一个：

```text
“像操作显微镜一样观察大模型内部结构”的前端可视化系统
```

核心目标：

* 可视化模型结构
* 探查参数与张量
* 实时观察推理激活
* 分析 Attention 行为
* 对比不同模型内部差异
* 为 AI Interpretability（可解释性）研究提供交互工具

后端推理与模型加载基于：

# llama.cpp

---

# 二、技术架构

# 前端技术栈

| 模块     | 技术                      |
| ------ | ----------------------- |
| 框架     | Vue 3                   |
| 语言     | TypeScript              |
| 构建工具   | Vite                    |
| 状态管理   | Pinia                   |
| 路由     | Vue Router              |
| UI     | Naive UI / Element Plus |
| 可视化引擎  | PixiJS（核心推荐）            |
| 图布局    | D3-force                |
| 图表     | ECharts                 |
| 热力图    | WebGL Shader            |
| 动画     | GSAP                    |
| Worker | WebWorker               |
| 数据通信   | WebSocket + REST API    |

---

# llama.cpp 后端结构

推荐：

```text
llama.cpp
 ├── HTTP Server
 ├── Token Streaming
 ├── KV Cache
 ├── Attention Export
 ├── Activation Hook
 ├── Tensor Metadata
 └── GGUF Parser
```

建议基于：

```bash
./server -m model.gguf --port 8080
```

扩展自定义接口。

---

# 三、系统整体架构

```text
Frontend (Vue3)
 ├── Visualization Engine
 ├── Model Explorer
 ├── Tensor Inspector
 ├── Activation Microscope
 ├── Attention Viewer
 ├── Compare Mode
 └── Export System

          ↓ REST/WebSocket

Backend (llama.cpp)
 ├── Model Runtime
 ├── Token Generator
 ├── Tensor Reader
 ├── Activation Collector
 ├── Attention Hook
 └── Metadata Service
```

---

# 四、页面结构设计

# 1. 首页（模型实验室）

## 功能

### 模型卡片列表

显示：

* 模型名称
* 参数量
* 架构类型
* Context Length
* Quantization 类型
* 文件大小

示例：

```text
LLaMA-2-7B
Q4_K_M
7B Params
4K Context
```

---

## 模型导入

支持：

* 本地 GGUF
* HuggingFace config.json
* llama.cpp server endpoint

---

## 模型加载状态

显示：

* 加载进度
* 显存占用
* KV Cache 状态
* GPU Offload 层数

---

# 五、模型总览仪表盘（Low Magnification）

# 页面目标

快速了解模型整体结构。

---

## 模块 1：参数统计卡片

显示：

| 指标               | 内容             |
| ---------------- | -------------- |
| Total Params     | 总参数量           |
| Trainable Params | 可训练参数          |
| Layers           | Transformer 层数 |
| Hidden Size      | 隐藏维度           |
| Heads            | Attention 头数   |
| Vocab Size       | 词表大小           |

---

## 模块 2：模型结构缩略图

```text
Embedding
   ↓
Transformer × N
   ↓
LM Head
```

点击可跳转结构显微镜。

---

## 模块 3：层级参数分布图

图表：

* 每层参数量
* Attention / FFN 占比
* 显存占用

使用：

```text
ECharts Bar Chart
```

---

# 六、模型结构显微镜（核心模块）

# 核心目标

像显微镜一样：

```text
宏观 → 模块 → 张量 → 元素
```

逐层深入观察模型。

---

# 1. 可交互结构树

## 支持：

* 展开 / 折叠
* 节点高亮
* 搜索定位
* 动态加载

结构：

```text
Model
 ├── Embedding
 ├── Layer.0
 │    ├── Attention
 │    ├── FFN
 │    └── Norm
 ├── Layer.1
 └── LM Head
```

---

# 2. 无限缩放画布

## 功能

| 操作      | 行为   |
| ------- | ---- |
| 滚轮      | 缩放   |
| 左键拖动    | 平移   |
| 双击      | 聚焦节点 |
| MiniMap | 鸟瞰导航 |

---

# 3. 焦平面系统（核心创新）

## Focus Level

### Level 1：Layer

仅显示：

```text
Transformer Blocks
```

---

### Level 2：Module

展开：

```text
Attention
FFN
Norm
```

---

### Level 3：Tensor

显示：

```text
q_proj.weight
k_proj.weight
v_proj.weight
```

---

### Level 4：Element

查看：

```text
Tensor 数值级细节
```

---

# 4. 节点详情面板

右侧动态面板：

| 字段           | 内容        |
| ------------ | --------- |
| Name         | 参数名       |
| Shape        | 张量维度      |
| Dtype        | FP16/BF16 |
| Params       | 参数量       |
| Device       | CPU/GPU   |
| Quantization | Q4/Q5/Q8  |
| Min/Max      | 数值范围      |
| Mean/Std     | 统计值       |

---

# 七、Attention 显微镜

# 目标

观察 Transformer 注意力机制。

---

# 1. Attention Head 网格

显示：

```text
Head 0
Head 1
Head 2
...
```

每个头：

* 小型热力图预览
* Hover 高亮
* 点击放大

---

# 2. 全屏 Heatmap

支持：

* Token Hover
* Attention Score 查看
* 缩放
* 平移
* 截图

---

# 3. 层间对比

比较：

```text
Layer 1 vs Layer 20
```

观察：

* 注意力模式变化
* 长距离依赖
* Token 聚焦差异

---

# 4. llama.cpp Attention Hook

后端需导出：

```cpp
attention_scores[layer][head]
```

前端实时接收：

```json
{
  "layer": 12,
  "head": 3,
  "tokens": [...],
  "matrix": [...]
}
```

---

# 八、参数探查系统

# 1. 参数搜索

支持：

```text
q_proj.weight
mlp.down_proj
```

模糊搜索。

---

# 2. 参数表格

显示：

* Shape
* Dtype
* Quantization
* Mean
* Std
* Sparsity

---

# 3. 权重分布图

图表：

* Histogram
* Density Curve
* Outlier Detection

---

# 4. Tensor Viewer

支持：

| 类型  | 展示           |
| --- | ------------ |
| 1D  | 曲线           |
| 2D  | 热力图          |
| 3D+ | Slice Viewer |

---

# 九、激活显微镜（Inference Microscope）

# 核心目标

观察：

```text
Prompt → 激活 → 输出
```

全过程。

---

# 1. Prompt 输入区

用户输入：

```text
Explain transformers
```

点击：

```text
Run Forward Pass
```

---

# 2. 激活热力覆盖

结构图节点颜色：

```text
颜色越亮 → 激活越强
```

---

# 3. 神经元激活柱状图

显示：

* Top-K neurons
* 激活强度
* 稀疏度

---

# 4. Token 时间轴

支持：

```text
Token-by-token replay
```

观察：

* Attention 演化
* Hidden State 变化
* 输出 logits

---

# 5. llama.cpp Hook 扩展

后端导出：

```cpp
hidden_states[layer]
mlp_activations[layer]
attention_scores[layer]
```

推荐：

* WebSocket 实时推送
* Binary tensor stream

---

# 十、双模型对比模式

# 页面布局

```text
Left  : Model A
Right : Model B
```

同步：

* 缩放
* 平移
* Token Position

---

# 对比内容

| 类型         | 内容                |
| ---------- | ----------------- |
| 参数量        | 差异                |
| Attention  | Pattern Diff      |
| Activation | 激活变化              |
| Tensor     | Weight Difference |

---

# 十一、导出系统

# 1. PNG 导出

导出：

* 当前画布
* 热力图
* Attention View

---

# 2. JSON 导出

导出：

```json
{
  "layer": 12,
  "mean": 0.002,
  "std": 0.12
}
```

---

# 3. PDF 报告

生成：

* 参数统计
* 模型结构
* Attention 分析

---

# 十二、状态管理设计（Pinia）

# store 结构

```text
stores/
 ├── modelStore.ts
 ├── viewportStore.ts
 ├── tensorStore.ts
 ├── activationStore.ts
 ├── compareStore.ts
 ├── websocketStore.ts
 └── uiStore.ts
```

---

# 十三、WebSocket 实时协议

# Token Stream

```json
{
  "type": "token",
  "token": "Hello",
  "position": 12
}
```

---

# Activation Stream

```json
{
  "type": "activation",
  "layer": 15,
  "mean": 0.23,
  "max": 1.2
}
```

---

# Attention Stream

```json
{
  "type": "attention",
  "layer": 10,
  "head": 5,
  "matrix": []
}
```

---

# 十四、目录结构（推荐）

```text
src/
 ├── api/
 ├── components/
 ├── views/
 ├── stores/
 ├── workers/
 ├── engine/
 ├── composables/
 ├── shaders/
 ├── router/
 └── utils/
```

---

# 十五、性能优化方案

# 1. WebWorker

用于：

* Tensor 统计
* Heatmap 计算
* 大矩阵处理

---

# 2. LOD（Level of Detail）

缩放不同层级：

| 缩放级别 | 显示      |
| ---- | ------- |
| 远景   | Layer   |
| 中景   | Module  |
| 近景   | Tensor  |
| 超近景  | Element |

---

# 3. 虚拟渲染

避免：

```text
10万节点同时渲染
```

---

# 4. GPU Heatmap

使用：

```text
WebGL Shader
```

绘制 attention matrix。

---

# 十六、快捷键系统

| 快捷键    | 功能           |
| ------ | ------------ |
| F      | 搜索           |
| 1      | Layer View   |
| 2      | Module View  |
| 3      | Tensor View  |
| 4      | Element View |
| Space  | Reset Camera |
| Ctrl+Z | 返回上一步        |

---

# 十七、暗色模式

默认：

```text
白色实验室风格
```

支持：

```text
Dark Research Mode
```

---

# 十八、MVP 开发优先级（非常重要）

# Phase 1（核心 MVP）

## 必做

* 模型加载
* 结构树
* 缩放画布
* 参数查看
* Attention Heatmap
* Prompt 推理
* Activation Overlay

---

# Phase 2（增强）

* 双模型对比
* Token Replay
* PDF 导出
* 高维 Tensor Viewer

---

# Phase 3（高级研究功能）

* 神经元聚类
* Feature Visualization
* SAE 可视化
* Mechanistic Interpretability

---

# 十九、推荐 llama.cpp 修改点

# 需要新增：

| 功能                | llama.cpp 修改          |
| ----------------- | --------------------- |
| Attention 导出      | Hook attention scores |
| Hidden State 导出   | Intermediate tensors  |
| Activation Stream | WebSocket             |
| Tensor Metadata   | GGUF parser           |
| Token Replay      | Step execution        |

---

# 二十、最终产品目标

LLMScope 不只是：

```text
“模型参数查看器”
```

而是：

# 一个 AI 模型内部世界的交互式显微镜。

用户可以：

```text
从模型整体结构
一路深入到：
Attention Head
Tensor
Neuron
甚至单个权重值
```

并实时观察：

```text
Prompt 如何在模型内部流动与演化。
```
