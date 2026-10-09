import { ss } from '@/utils/storage'

const LOCAL_NAME = 'appSetting'

export type Theme = 'light' | 'dark' | 'auto'

export type Language = 'zh-CN' | 'en-US'

export interface SiteSetting {
  title: string
  favicoUrl: string
  loginBackgroundUrl: string
  isEnableLoginCaptcha: boolean
  isDisableItemCardCache: boolean
}

export function defaultSiteSetting(): SiteSetting {
  return {
    title: 'Sun-Panel',
    favicoUrl: '',
    loginBackgroundUrl: '',
    isEnableLoginCaptcha: false,
    isDisableItemCardCache: false,
  }
}

export interface AppState {
  siderCollapsed: boolean
  theme: Theme
  language: Language
  siteSetting: SiteSetting
}

export function defaultSetting(): AppState {
  const lan = (navigator.language).toLowerCase()
  let language: Language = 'en-US'
  if (lan.includes('zh'))
    language = 'zh-CN'

  return { siderCollapsed: false, theme: 'auto', language, siteSetting: defaultSiteSetting() }
}

export function getLocalSetting(): AppState {
  const localSetting: AppState | undefined = ss.get(LOCAL_NAME)
  return { ...defaultSetting(), ...localSetting }
}

export function setLocalSetting(setting: AppState): void {
  ss.set(LOCAL_NAME, setting)
}

export function removeLocalState() {
  ss.remove(LOCAL_NAME)
}
