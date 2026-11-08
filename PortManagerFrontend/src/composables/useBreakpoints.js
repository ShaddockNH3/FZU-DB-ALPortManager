import { computed, onMounted, onUnmounted, ref } from 'vue'

/**
 * 视口宽度断点。表格列和分页要按宽度伸缩，纯 CSS 做不到（列是 JS 定义的），
 * 所以监听一次 window resize，两个断点够用。
 */
export function useBreakpoints() {
  const width = ref(typeof window === 'undefined' ? 1280 : window.innerWidth)

  const update = () => {
    width.value = window.innerWidth
  }

  onMounted(() => window.addEventListener('resize', update))
  onUnmounted(() => window.removeEventListener('resize', update))

  return {
    width,
    /** < 1100：表格藏掉次要列 */
    isCompact: computed(() => width.value < 1100),
    /** < 720：筛选条改成竖排、分页退成简版 */
    isNarrow: computed(() => width.value < 720)
  }
}
