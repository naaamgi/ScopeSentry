import { AxiosResponse, AxiosRequestHeaders, InternalAxiosRequestConfig } from './types'
import { ElMessage } from 'element-plus'
import qs from 'qs'
import { SUCCESS_CODE } from '@/constants'
import { useUserStoreWithOut } from '@/store/modules/user'
import { useStorage } from '@/hooks/web/useStorage'

const { getStorage } = useStorage()

// 서버가 응답 메시지를 어떤 언어로 돌려줄지 결정하는 값.
// 서버 쪽에서 ko -> ko-KR, en -> en-US 로 환산하므로 UI 언어 코드를 그대로 보낸다.
const getRequestLocale = (): string => {
  try {
    return getStorage('lang') || 'zh-CN'
  } catch {
    return 'zh-CN'
  }
}

const defaultRequestInterceptors = (config: InternalAxiosRequestConfig) => {
  ;(config.headers as AxiosRequestHeaders)['Accept-Language'] = getRequestLocale()
  if (
    config.method === 'post' &&
    (config.headers as AxiosRequestHeaders)['Content-Type'] === 'application/x-www-form-urlencoded'
  ) {
    config.data = qs.stringify(config.data)
  }
  if (config.method === 'get' && config.params) {
    let url = config.url as string
    url += '?'
    const keys = Object.keys(config.params)
    for (const key of keys) {
      if (config.params[key] !== void 0 && config.params[key] !== null) {
        url += `${key}=${encodeURIComponent(config.params[key])}&`
      }
    }
    url = url.substring(0, url.length - 1)
    config.params = {}
    config.url = url
  }
  return config
}

const defaultResponseInterceptors = (response: AxiosResponse) => {
  if (response?.headers['content-type'] == 'application/octet-stream') {
    return response
  }
  if (response?.config?.responseType === 'blob') {
    // 如果是文件流，直接过
    return response
  } else if (response.data.code === SUCCESS_CODE) {
    if (response?.data?.message) {
      ElMessage.success(response?.data?.message)
    }
    return response.data
  } else {
    ElMessage.error(response?.data?.message)
    if (response?.data?.code === 401) {
      const userStore = useUserStoreWithOut()
      userStore.logout()
    }
    return response.data
  }
}

export { defaultResponseInterceptors, defaultRequestInterceptors }
