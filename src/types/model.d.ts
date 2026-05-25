export interface ModelStructure {
  embedding: LayerNode
  transformers: TransformerLayer[]
  lmHead: LayerNode
}

export interface TransformerLayer {
  id: number
  attention: AttentionModule
  ffn: FFNModule
  norm: NormModule
}

export interface LayerNode {
  name: string
  type: string
  params: number
}

export interface AttentionModule {
  qProj: TensorNode
  kProj: TensorNode
  vProj: TensorNode
  oProj: TensorNode
  heads: number
}

export interface FFNModule {
  upProj: TensorNode
  downProj: TensorNode
  gateProj?: TensorNode
}

export interface NormModule {
  weight: TensorNode
  bias?: TensorNode
}

export interface TensorNode {
  name: string
  shape: number[]
  dtype: string
}
