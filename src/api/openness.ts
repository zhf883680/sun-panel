import { get } from '@/utils/request'

export function getLoginConfig<T>() {
  return get<T>({
    url: '/openness/loginConfig',
  })
}

// 获取免责声明
export function getDisclaimer<T>() {
  return get<T>({
    url: '/openness/getDisclaimer',
  })
}

// 获取关于的描述信息
export function getAboutDescription<T>() {
  return get<T>({
    url: '/openness/getAboutDescription',
  })
}

// 站点设置（登录页使用，无需鉴权）
export function getSiteSetting<T>() {
  return get<T>({
    url: '/openness/getSiteSetting',
  })
}
