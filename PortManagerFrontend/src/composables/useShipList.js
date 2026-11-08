import { onMounted, reactive } from 'vue'
import { storeToRefs } from 'pinia'
import { useShipStore } from '@/stores/ship'
import { dialog, message } from '@/utils/feedback'
import { DEFAULT_PAGE_SIZE, initialStars, maxStars } from '@/constants/ship'

// 下拉框空值用 null：空字符串在 Naive 里算“有值”，placeholder 就不显示了
const emptyFilter = () => ({ shipName: '', rarity: null, shipType: null, faction: null })

/**
 * 舰船列表页的状态与用例：加载、筛选、分页、升星、删除。
 *
 * 列表数据只有一份：store 里的 ships / total / loading。
 * 之前这里又存了一份 ref，两边靠手写赋值同步，改一处忘了另一处就开始飘。
 */
export function useShipList() {
  const store = useShipStore()
  const { ships, total, loading } = storeToRefs(store)

  const filterForm = reactive(emptyFilter())
  const pagination = reactive({ page: 1, pageSize: DEFAULT_PAGE_SIZE })

  /** 拉列表。空筛选项不发，免得后端按空字符串过滤。 */
  async function load() {
    const params = { page: pagination.page, pageSize: pagination.pageSize }
    Object.entries(filterForm).forEach(([key, value]) => {
      if (value) params[key] = value
    })

    try {
      await store.fetchShipList(params)
    } catch (error) {
      // 提示由拦截器统一弹，这里只留日志
      console.error('加载舰船列表失败:', error)
    }
  }

  function search() {
    pagination.page = 1
    return load()
  }

  function reset() {
    Object.assign(filterForm, emptyFilter())
    return search()
  }

  function changePage(page) {
    pagination.page = page
    return load()
  }

  function changePageSize(size) {
    pagination.pageSize = size
    pagination.page = 1
    return load()
  }

  /** 升星/降星：夹在初始星与满星之间，成功后整页重拉。 */
  async function changeStars(ship, delta) {
    const nextStars = ship.stars + delta
    if (nextStars < initialStars(ship.rarity) || nextStars > maxStars(ship.rarity)) {
      return
    }

    try {
      await store.updateShip(ship.id, {
        shipName: ship.shipName,
        rarity: ship.rarity,
        shipType: ship.shipType,
        faction: ship.faction,
        level: ship.level,
        stars: nextStars
      })
      message.success(delta > 0 ? '升星成功' : '降星成功')
      await load()
    } catch (error) {
      console.error('星级调整失败:', error)
    }
  }

  /** 删除前确认。Naive 的 dialog 没有 Promise 版本，这里包一层。 */
  function confirmDelete(ship) {
    return new Promise((resolve) => {
      dialog.warning({
        title: '删除确认',
        content: `确定要删除舰船「${ship.shipName}」吗？`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: () => resolve(true),
        onNegativeClick: () => resolve(false),
        onClose: () => resolve(false),
        onMaskClick: () => resolve(false)
      })
    })
  }

  async function remove(ship) {
    if (!(await confirmDelete(ship))) {
      return
    }

    try {
      await store.deleteShip(ship.id)
      message.success('删除成功')
      await load()
    } catch (error) {
      console.error('删除失败:', error)
    }
  }

  onMounted(load)

  return {
    filterForm,
    pagination,
    ships,
    total,
    loading,
    load,
    search,
    reset,
    changePage,
    changePageSize,
    changeStars,
    remove
  }
}
