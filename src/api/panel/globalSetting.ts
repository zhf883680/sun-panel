import { post } from '@/utils/request'

export interface SiteSetting {
  title: string
  favicoUrl: string
  loginBackgroundUrl: string
  isEnableLoginCaptcha: boolean
  isDisableItemCardCache: boolean
}

export interface CustomJsAndCss {
  jsContent: string
  cssContent: string
}

// 获取站点设置
export function getSiteStting<T>() {
  return post<T>({
    url: '/panel/globalSetting/getSiteStting',
  })
}

// 保存站点设置
export function saveSiteStting<T>(data: SiteSetting) {
  return post<T>({
    url: '/panel/globalSetting/saveSiteStting',
    data,
  })
}

// 自定义JS/CSS所在目录
export function getCustomStylePath<T>() {
  return post<T>({
    url: '/panel/globalSetting/getCustomStylePath',
  })
}

// 读取自定义JS/CSS内容
export function getCustomJsAndCssCode<T>() {
  return post<T>({
    url: '/panel/globalSetting/getCustomJsAndCssCode',
  })
}

// 保存自定义JS/CSS内容
export function saveCustomJsAndCssCode<T>(data: CustomJsAndCss) {
  return post<T>({
    url: '/panel/globalSetting/saveCustomJsAndCssCode',
    data,
  })
}
