import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useShipStore } from '@/stores/ship'
import { initialStars, maxStars } from '@/constants/ship'

const emptyFilter = () => ({ shipName: '', rarity: '', shipType: '', faction: '' })

/**
 * 舰船列表页的状态与用例：加载、筛选、分页、升星、删除。
 *
 * 原来这些全写在 ShipList.vue 的 setup 里，和模板、对话框搅在一起；
 * 抽出来之后页面组件只负责拼装，逻辑可以单独读。
 */
export function useShipList() {
  const shipStore = useShipStore()

  const filterForm = reactive(emptyFilter())
  const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
  const ships = ref([])
  const loading = ref(false)

  /** 拉列表。空筛选项不发，免得后端按空字符串过滤。 */
  async function load() {
    loading.value = true
    try {
      const params = { page: pagination.page, pageSize: pagination.pageSize }
      Object.entries(filterForm).forEach(([key, value]) => {
        if (value) params[key] = value
      })

      const data = await shipStore.fetchShipList(params)
      ships.value = data.ships || []
      pagination.total = data.total || 0
    } catch (error) {
      console.error('加载数据失败:', error)
      ElMessage.error('加载舰船列表失败')
    } finally {
      loading.value = false
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
      await shipStore.updateShip(ship.id, {
        shipName: ship.shipName,
        rarity: ship.rarity,
        shipType: ship.shipType,
        faction: ship.faction,
        level: ship.level,
        stars: nextStars
      })
      ElMessage.success(delta > 0 ? '升星成功' : '降星成功')
      await load()
    } catch (error) {
      console.error('星级调整失败:', error)
      ElMessage.error('星级调整失败')
    }
  }

  async function remove(ship) {
    try {
      await ElMessageBox.confirm(`确定要删除舰船"${ship.shipName}"吗？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return // 用户点了取消
    }

    try {
      await shipStore.deleteShip(ship.id)
      ElMessage.success('删除成功')
      await load()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }

  onMounted(load)

  return {
    filterForm,
    pagination,
    ships,
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
