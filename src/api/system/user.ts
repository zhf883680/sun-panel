import { post } from '@/utils/request'

// 可指定 token，用于多账号场景查询其它账号信息
export function getInfo<T>(token?: string) {
  return post<T>({
    url: '/user/getInfo',
    headers: token ? { token } : undefined,
  })
}

export function getAuthInfo<T>() {
  return post<T>({
    url: '/user/getAuthInfo',
  })
}

export function getReferralCode<T>() {
  return post<T>({
    url: '/user/getReferralCode',
  })
}

export function updateInfo<T>(name: string) {
  return post<T>({
    url: '/user/updateInfo',
    data: { name },
  })
}

export function updatePassword<T>(oldPassword: string, newPassword: string) {
  return post<T>({
    url: '/user/updatePassword',
    data: { newPassword, oldPassword },
  })
}
