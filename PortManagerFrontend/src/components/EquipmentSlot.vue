<script setup>
import { acceptableTypesLabel } from '@/constants/ship'

/**
 * 单个装备槽位：已装就显示一个可移除的标签，没装就一个按钮。
 * 普通槽和兵装槽共用，差别只在标题、配色和锁定原因。
 */
defineProps({
  slot: { type: Object, required: true },
  index: { type: Number, default: 0 },
  title: { type: String, required: true },
  kind: { type: String, default: 'normal' },
  locked: { type: Boolean, default: false },
  lockedText: { type: String, default: '不可选' },
  disabled: { type: Boolean, default: false }
})

const emit = defineEmits(['pick', 'unequip'])
</script>

<template>
  <div class="equipment-slot" :class="{ 'is-augment': kind === 'augment' }">
    <div class="slot-label">
      {{ title }}
      <span v-if="locked" class="tip">({{ lockedText }})</span>
    </div>

    <div class="slot-content">
      <el-tag
        v-if="slot.equippedItem"
        :type="kind === 'augment' ? 'warning' : 'success'"
        closable
        @close="emit('unequip', index)"
      >
        {{ slot.equippedItem.name }}
      </el-tag>

      <el-button
        v-else
        type="primary"
        size="small"
        :disabled="locked || disabled"
        @click="emit('pick', index, slot.acceptableTypes)"
      >
        {{ locked ? lockedText : '选择装备' }}
      </el-button>
    </div>

    <div class="acceptable-types">可装备类型: {{ acceptableTypesLabel(slot.acceptableTypes) }}</div>
  </div>
</template>

<style lang="scss" scoped>
.equipment-slot {
  padding: 12px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: #fff;

  &.is-augment {
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
