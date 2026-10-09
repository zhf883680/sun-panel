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
}

// 获取 HD-Icons 图标索引（后端代理 + 缓存）
export function getHdIconsIndex<T>() {
  return get<T>({
    url: '/system/gallery/hdicons/index',
  })
}

// HD-Icons 图标地址（走后端代理，解决跨域）
export function getHdIconsImageUrl(style: string, file: string) {
  return `/api/system/gallery/hdicons/image?style=${encodeURIComponent(style)}&file=${encodeURIComponent(file)}`
}
