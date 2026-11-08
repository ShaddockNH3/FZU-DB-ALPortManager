import { computed, ref, watch } from 'vue'
import { equipShip, getEquipmentList, getShipById } from '@/api/ship'
import { message } from '@/utils/feedback'
import { AUGMENT_SLOT_INDEX, augmentsForShipType, isMaxStars } from '@/constants/ship'

/**
 * 装备对话框的用例：选装备、装上、卸下、刷新。
 *
 * 拿到的是舰船快照，装/卸之后重新拉一份，避免本地状态和后端不一致。
 * 组件只负责画，判断都在这里。
 *
 * @param {() => object|null} getShip 父级传下来的舰船（可能为 null）
 */
export function useEquipmentDialog(getShip) {
  const currentShip = ref(null)
  const equipments = ref([])
  const candidates = ref([])
  const pickerVisible = ref(false)
  const currentSlot = ref(null)
  const busy = ref(false)

  watch(
    getShip,
    (value) => {
      currentShip.value = value ? { ...value } : null
    },
    { immediate: true }
  )

  /** 兵装要满星才解锁，规则和后端一致。 */
  const augmentLocked = computed(() => !isMaxStars(currentShip.value))

  /** 装备列表只拉一次，之后按槽位可装类型过滤。 */
  async function ensureEquipments() {
    if (equipments.value.length > 0) return true

    try {
      const data = await getEquipmentList()
      equipments.value = data.equipments || []
      return true
    } catch (error) {
      console.error('获取装备列表失败:', error)
      return false
    }
  }

  async function openPicker(slotIndex, acceptableTypes) {
    if (!(await ensureEquipments())) return

    let pool = equipments.value.filter((item) => (acceptableTypes || []).includes(item.type))
    if (slotIndex === AUGMENT_SLOT_INDEX) {
      pool = augmentsForShipType(pool, currentShip.value?.shipType)
    }

    candidates.value = pool
    currentSlot.value = slotIndex

    if (pool.length === 0) {
      message.warning('该栏位暂无可用装备')
      return
    }

    pickerVisible.value = true
  }

  async function apply(equipment) {
    busy.value = true
    try {
      await equipShip(currentShip.value.id, {
        slotIndex: currentSlot.value,
        equipmentId: equipment.id
      })
      message.success('装备成功')
      pickerVisible.value = false
      await refresh()
    } catch (error) {
      console.error('装备失败:', error)
    } finally {
      busy.value = false
    }
  }

  async function unequip(slotIndex) {
    busy.value = true
    try {
      await equipShip(currentShip.value.id, { slotIndex, equipmentId: 0 })
      message.success('卸下装备成功')
      await refresh()
    } catch (error) {
      console.error('卸下装备失败:', error)
    } finally {
      busy.value = false
    }
  }

  async function refresh() {
    if (!currentShip.value) return

    try {
      currentShip.value = await getShipById(currentShip.value.id)
    } catch (error) {
      console.error('刷新舰船信息失败:', error)
    }
  }

  return {
    currentShip,
    candidates,
    pickerVisible,
    busy,
    augmentLocked,
    openPicker,
    apply,
    unequip
  }
}
