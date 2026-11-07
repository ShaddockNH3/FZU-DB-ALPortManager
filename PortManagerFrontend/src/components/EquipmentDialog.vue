<script setup>
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { equipShip, getEquipmentList, getShipById } from '@/api/ship'
import { acceptableTypesLabel, augmentsForShipType, equipmentTypeLabel, isMaxStars } from '@/constants/ship'

/**
 * 装备管理对话框：6 个普通槽 + 1 个兵装槽。
 *
 * 进来时拿到的是一份舰船快照，装/卸之后重新拉一份，避免本地状态和后端不一致；
 * 关闭时发 changed，让父级刷新列表。
 */
const props = defineProps({
  visible: { type: Boolean, default: false },
  ship: { type: Object, default: null }
})

const emit = defineEmits(['update:visible', 'changed'])

const AUGMENT_SLOT_INDEX = 5

const currentShip = ref(null)
const allEquipments = ref([])
const filteredEquipments = ref([])
const selectVisible = ref(false)
const currentSlot = ref(null)
const busy = ref(false)

watch(
  () => props.ship,
  (value) => {
    currentShip.value = value ? { ...value } : null
  },
  { immediate: true }
)

const augmentLocked = computed(() => !isMaxStars(currentShip.value))

function close() {
  emit('update:visible', false)
}

async function openSlot(slotIndex, acceptableTypes) {
  currentSlot.value = slotIndex

  if (allEquipments.value.length === 0) {
    try {
      const data = await getEquipmentList()
      allEquipments.value = data.equipments || []
    } catch (error) {
      console.error('获取装备列表失败:', error)
      ElMessage.error('获取装备列表失败')
      return
    }
  }

  let candidates = allEquipments.value.filter((item) => acceptableTypes.includes(item.type))
  if (slotIndex === AUGMENT_SLOT_INDEX) {
    candidates = augmentsForShipType(candidates, currentShip.value?.shipType)
  }

  filteredEquipments.value = candidates

  if (candidates.length === 0) {
    ElMessage.warning('该栏位暂无可用装备')
    return
  }

  selectVisible.value = true
}

async function applyEquipment(equipment) {
  busy.value = true
  try {
    await equipShip(currentShip.value.id, {
      slotIndex: currentSlot.value,
      equipmentId: equipment.id
    })
    ElMessage.success('装备成功')
    selectVisible.value = false
    await refreshShip()
  } catch (error) {
    console.error('装备失败:', error)
    ElMessage.error(`装备失败: ${error.message || error}`)
  } finally {
    busy.value = false
  }
}

async function unequip(slotIndex) {
  busy.value = true
  try {
    await equipShip(currentShip.value.id, { slotIndex, equipmentId: 0 })
    ElMessage.success('卸下装备成功')
    await refreshShip()
  } catch (error) {
    console.error('卸下装备失败:', error)
    ElMessage.error('卸下装备失败')
  } finally {
    busy.value = false
  }
}

async function refreshShip() {
  if (!currentShip.value) return
  try {
    currentShip.value = await getShipById(currentShip.value.id)
  } catch (error) {
    console.error('刷新舰船信息失败:', error)
  }
}
</script>

<template>
  <el-dialog
    :model-value="visible"
    :title="`装备管理 - ${currentShip?.shipName || ''}`"
    width="800px"
    @update:model-value="close"
    @closed="emit('changed')"
  >
    <div v-if="currentShip" v-loading="busy" class="equipment-container">
      <div v-for="(slot, index) in currentShip.equipmentSlots" :key="index" class="equipment-slot">
        <div class="slot-label">装备栏 {{ index + 1 }}</div>

        <div class="slot-content">
          <el-tag v-if="slot.equippedItem" type="success" closable @close="unequip(index)">
            {{ slot.equippedItem.name }}
          </el-tag>
          <el-button
            v-else
            type="primary"
            size="small"
            @click="openSlot(index, slot.acceptableTypes)"
          >
            选择装备
          </el-button>
        </div>

        <div class="acceptable-types">可装备类型: {{ acceptableTypesLabel(slot.acceptableTypes) }}</div>
      </div>

      <div v-if="currentShip.augmentSlot" class="equipment-slot augment-slot">
        <div class="slot-label">兵装栏 <span class="tip">(需满星)</span></div>

        <div class="slot-content">
          <el-tag
            v-if="currentShip.augmentSlot.equippedItem"
            type="warning"
            closable
            @close="unequip(5)"
          >
            {{ currentShip.augmentSlot.equippedItem.name }}
          </el-tag>
          <el-button
            v-else
            type="primary"
            size="small"
            :disabled="augmentLocked"
            @click="openSlot(5, currentShip.augmentSlot.acceptableTypes)"
          >
            {{ augmentLocked ? '未满星' : '选择兵装' }}
          </el-button>
        </div>
      </div>
    </div>

    <el-dialog v-model="selectVisible" title="选择装备" width="600px" append-to-body>
      <el-table
        :data="filteredEquipments"
        highlight-current-row
        style="width: 100%"
        @row-click="applyEquipment"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="装备名称" min-width="200" />
        <el-table-column prop="type" label="类型" width="150">
          <template #default="{ row }">{{ equipmentTypeLabel(row.type) }}</template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </el-dialog>
</template>

<style lang="scss" scoped>
.equipment-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.equipment-slot {
  padding: 12px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: #fff;

  &.augment-slot {
    border-color: var(--el-color-warning-light-5);
  }
}

.slot-label {
  font-weight: 600;
  margin-bottom: 8px;

  .tip {
    font-weight: 400;
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
}

.slot-content {
  min-height: 32px;
  display: flex;
  align-items: center;
}

.acceptable-types {
  margin-top: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
