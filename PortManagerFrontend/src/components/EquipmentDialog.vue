<script setup>
import { AUGMENT_SLOT_INDEX, equipmentTypeLabel } from '@/constants/ship'
import { useEquipmentDialog } from '@/composables/useEquipmentDialog'
import EquipmentSlot from '@/components/EquipmentSlot.vue'

/**
 * 装备管理对话框：6 个普通槽 + 1 个兵装槽。
 *
 * 逻辑在 useEquipmentDialog，单个槽位的样式在 EquipmentSlot，
 * 这里只剩版面和一个装备选择弹窗。
 */
const props = defineProps({
  visible: { type: Boolean, default: false },
  ship: { type: Object, default: null }
})

const emit = defineEmits(['update:visible', 'changed'])

const {
  currentShip,
  candidates,
  pickerVisible,
  busy,
  augmentLocked,
  openPicker,
  apply,
  unequip
} = useEquipmentDialog(() => props.ship)

function close() {
  emit('update:visible', false)
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
      <equipment-slot
        v-for="(slot, index) in currentShip.equipmentSlots"
        :key="index"
        :slot="slot"
        :index="index"
        :title="`装备栏 ${index + 1}`"
        :disabled="busy"
        @pick="openPicker"
        @unequip="unequip"
      />

      <equipment-slot
        v-if="currentShip.augmentSlot"
        :slot="currentShip.augmentSlot"
        :index="AUGMENT_SLOT_INDEX"
        title="兵装栏"
        kind="augment"
        :locked="augmentLocked"
        locked-text="需满星"
        :disabled="busy"
        @pick="openPicker"
        @unequip="unequip"
      />
    </div>

    <el-dialog v-model="pickerVisible" title="选择装备" width="600px" append-to-body>
      <el-table :data="candidates" highlight-current-row style="width: 100%" @row-click="apply">
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
</style>
