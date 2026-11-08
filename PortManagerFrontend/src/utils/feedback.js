import { createDiscreteApi } from 'naive-ui'
import { themeOverrides } from '@/theme'

/**
 * 组件外面也要弹提示：api 拦截器、composable 里的异步流程。
 *
 * Naive UI 的 useMessage 只能在组件 setup 里用，所以这里用离散 API 挂一份全局实例。
 * 全站提示只从这一处出，别在业务代码里再 ElMessage / message 各弹一次。
 */
const { message, dialog } = createDiscreteApi(['message', 'dialog'], {
  configProviderProps: { themeOverrides }
})

export { message, dialog }
