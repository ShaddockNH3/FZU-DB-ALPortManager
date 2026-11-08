<script setup>
import { computed } from 'vue'
import { NDataTable, NModal } from 'naive-ui'
import { AUGMENT_SLOT_INDEX, equipmentTypeLabel } from '@/constants/ship'
import { useEquipmentDialog } from '@/composables/useEquipmentDialog'
import EquipmentSlot from '@/components/EquipmentSlot.vue'

/**
 * 装备管理对话框：5～6 个普通槽 + 1 个兵装槽。
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

const show = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const pickerColumns = [
  { title: 'ID', key: 'id', width: 70 },
  { title: '装备名称', key: 'name', minWidth: 200 },
  {
    title: '类型',
    key: 'type',
    width: 120,
    render: (row) => equipmentTypeLabel(row.type)
  }
]

function close() {
  show.value = false
}

function onClosed() {
  emit('changed')
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="`装备管理 · ${currentShip?.shipName || ''}`"
    style="width: 860px; max-width: 94vw"
    @after-leave="onClosed"
  >
    <div v-if="currentShip" class="slots" :class="{ 'is-busy': busy }">
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

    <n-modal v-model:show="pickerVisible" preset="card" title="选择装备" style="width: 620px">
      <n-data-table
        :columns="pickerColumns"
        :data="candidates"
        :bordered="false"
        :max-height="360"
        size="small"
        :row-props="(row) => ({ style: 'cursor:pointer', onClick: () => apply(row) })"
      />
    </n-modal>
  </n-modal>
</template>

<style lang="scss" scoped>
.slots {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

.slots.is-busy {
  opacity: 0.6;
  pointer-events: none;
}
</style>
