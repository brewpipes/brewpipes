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

  export function sankey<N = any, L = any> (): SankeyGenerator<N, L>

  export function sankeyLinkHorizontal<L = any> (): (link: L) => string | null

  export function sankeyLeft<N = any> (node: N, n: number): number

  export function sankeyRight<N = any> (node: N, n: number): number

  export function sankeyCenter<N = any> (node: N, n: number): number

  export function sankeyJustify<N = any> (node: N, n: number): number
}
