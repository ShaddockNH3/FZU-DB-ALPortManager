<script setup>
import { NButton, NCard, NIcon, NTag } from 'naive-ui'
import { PhPlusCircle } from '@phosphor-icons/vue'
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
  <n-card size="small" :class="['slot', { 'is-augment': kind === 'augment' }]">
    <template #header>
      <div class="slot-head">
        <span class="slot-title">{{ title }}</span>
        <n-tag v-if="locked" size="tiny" :bordered="false" type="warning">{{ lockedText }}</n-tag>
      </div>
    </template>

    <div class="slot-body">
      <n-tag
        v-if="slot.equippedItem"
        :type="kind === 'augment' ? 'warning' : 'success'"
        closable
        :bordered="false"
        @close="emit('unequip', index)"
      >
        {{ slot.equippedItem.name }}
      </n-tag>

      <n-button
        v-else
        size="small"
        type="primary"
        secondary
        :disabled="locked || disabled"
        @click="emit('pick', index, slot.acceptableTypes)"
      >
        <template #icon>
          <n-icon><PhPlusCircle /></n-icon>
        </template>
        {{ locked ? lockedText : '选择装备' }}
      </n-button>
    </div>

    <div class="slot-types">可装备：{{ acceptableTypesLabel(slot.acceptableTypes) }}</div>
  </n-card>
</template>

<style lang="scss" scoped>
.slot {
  background: #fff;
}

.slot.is-augment {
  background: #fffdf7;
}

.slot-head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.slot-title {
  font-weight: 600;
  color: #1f2937;
}

.slot-body {
  min-height: 32px;
  display: flex;
  align-items: center;

  // 装备名可能很长（双联装57mmL60博福斯…），让它换行，别撑破卡片
  :deep(.n-tag) {
    max-width: 100%;
    height: auto;
    min-height: 26px;
    padding: 2px 8px;
    white-space: normal;
    line-height: 1.4;
    text-align: left;
  }
}

.slot-types {
  margin-top: 10px;
  font-size: 12px;
  color: #8b95a7;
}
</style>
