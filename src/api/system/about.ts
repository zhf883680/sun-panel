import { post } from '@/utils/request'

export function get<T>() {
  return post<T>({
    url: '/about',
  })
}

// 运行信息
export function runningInfo<T>() {
  return post<T>({
    url: '/about/version/runningInfo',
  })
}

// 检查新版本
export function checkNewVersion<T>(isCheckBeta: boolean) {
  return post<T>({
    url: '/about/version/checkNewVersion',
    data: { isCheckBeta },
  })
}
