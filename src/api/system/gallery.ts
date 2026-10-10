import { get } from '@/utils/request'

export interface HdIconItem {
  name: string
  style: string
  file: string
}

export interface HdIconsIndex {
  list: HdIconItem[]
  name: string
  totalCount: number
  updateAt: string
  syncing?: boolean
}

export interface HdIconsListParams {
  page?: number
  pageSize?: number
  keyword?: string
  style?: string
}

export interface HdIconsListData {
  list: HdIconItem[]
  total: number
  page: number
  pageSize: number
  syncing: boolean
  syncedAt: string
  updateAt: string
  totalCount: number
}

// 获取 HD-Icons 图标索引（兼容旧接口，返回全量列表）
export function getHdIconsIndex<T>() {
  return get<T>({
    url: '/system/gallery/hdicons/index',
  })
}

// 分页 + 搜索 + 风格筛选（服务端本地数据，无需回源 GitHub）
export function getHdIconsList<T>(params: HdIconsListParams) {
  return get<T>({
    url: '/system/gallery/hdicons/list',
    data: params,
  })
}

// HD-Icons 图标地址（本地缓存位于 source_path/hdicon，
// 由 /uploads 静态路由直接提供，支持浏览器缓存/断点协商）
export function getHdIconsImageUrl(style: string, file: string) {
  return `/uploads/hdicon/${encodeURIComponent(style)}/${encodeURIComponent(file)}`
}

// 图标代理地址（后端回源兜底，静态文件缺失时可手动改用）
export function getHdIconsImageProxyUrl(style: string, file: string) {
  return `/api/system/gallery/hdicons/image?style=${encodeURIComponent(style)}&file=${encodeURIComponent(file)}`
}
