import axios, { type AxiosResponse } from 'axios'
import { useAuthStore } from '@/store'

const service = axios.create({
  baseURL: import.meta.env.VITE_GLOB_API_URL,
})

// 顶部加载进度条（naive-ui LoadingBarProvider）
let requestCount = 0

function startLoading() {
  requestCount++
  if (requestCount === 1)
    window.$loadingBar?.start()
}

function stopLoading() {
  requestCount = Math.max(0, requestCount - 1)
  if (requestCount === 0)
    window.$loadingBar?.finish()
}

service.interceptors.request.use(
  (config) => {
    const token = useAuthStore().token
    // 请求级 token 优先：多账号场景会使用其它账号的 token 调用接口
    if (token && !config.headers.token)
      config.headers.token = token
    startLoading()
    return config
  },
  (error) => {
    stopLoading()
    return Promise.reject(error.response)
  },
)

service.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    stopLoading()
    if (response.status === 200)
      return response

    throw new Error(response.status.toString())
  },
  (error) => {
    stopLoading()
    return Promise.reject(error)
  },
)

export default service
