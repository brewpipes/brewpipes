declare module 'd3-sankey' {
  export interface SankeyGraph<N, L> {
    nodes: N[]
    links: L[]
  }

  export interface SankeyGenerator<N, L, OutputLink = L> {
    (data: SankeyGraph<N, L>): SankeyGraph<N, OutputLink>
    nodeId(): (node: N) => string | number
    nodeId(id: (node: N) => string | number): this
    nodeWidth(): number
    nodeWidth(width: number): this
    nodePadding(): number
    nodePadding(padding: number): this
    extent(): [[number, number], [number, number]]
    extent(extent: [[number, number], [number, number]]): this
    nodeAlign(): (node: N, n: number) => number
    nodeAlign(align: (node: N, n: number) => number): this
    nodeSort(): ((a: N, b: N) => number) | null
    nodeSort(sort: ((a: N, b: N) => number) | null): this
    linkSort(): ((a: L, b: L) => number) | null
    linkSort(sort: ((a: L, b: L) => number) | null): this
  }

  export interface SankeyLink<N> {
    source: N
    target: N
    width?: number
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function sankey<N = any, L = any>(): SankeyGenerator<N, L>
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function sankeyLinkHorizontal<L = any>(): (link: L) => string | null
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function sankeyLeft<N = any>(node: N, n: number): number
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function sankeyRight<N = any>(node: N, n: number): number
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function sankeyCenter<N = any>(node: N, n: number): number
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function sankeyJustify<N = any>(node: N, n: number): number
}
