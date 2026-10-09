import { post } from '@/utils/request'

export function getList<T>(category?: string) {
  return post<T>({
    url: '/file/getList',
    data: category ? { category } : {},
  })
}

export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/file/deletes',
    data: { ids },
  })
}

// 设置文件分类（icon / wallpaper / unknown）
export function updateCategory<T>(id: number, category: string) {
  return post<T>({
    url: '/file/updateCategory',
    data: { id, category },
  })
}
