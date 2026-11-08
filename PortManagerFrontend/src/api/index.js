import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})

/**
 * 只做两件事：把 axios 的响应拆成业务数据，把错因摊平成一句能读的话。
 *
 * 提示在这里弹一次就够了——之前拦截器弹一次、调用方再弹一次，
 * 后端挂掉的时候页面会连着蹦两条一样的红条。
 */
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const message = error.response?.data?.msg || error.message || '请求失败'
    ElMessage.error(message)
    return Promise.reject(new Error(message))
  }
)

export default request
