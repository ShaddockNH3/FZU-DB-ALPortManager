<script setup>
import { computed } from 'vue'
import { NCard, NEmpty, NProgress } from 'naive-ui'

/**
 * 分布卡片：标题 + 一组「名称 / 进度条 / 数量」。
 * 阵营、稀有度、舰种三块原来各写一遍同样的模板。
 */
const props = defineProps({
  title: { type: String, required: true },
  items: { type: Array, default: () => [] },
  total: { type: Number, default: 0 },
  /** 传颜色或 (name) => 颜色 都行 */
  color: { type: [String, Function], default: '#1e40af' }
})

const percentageOf = (count) => {
  if (!props.total) return 0
  return Math.round((count / props.total) * 100)
}

const colorOf = computed(() => {
  return typeof props.color === 'function' ? props.color : () => props.color
})
</script>

<template>
  <n-card size="small">
    <template #header>
      <span class="card-title">{{ title }}</span>
    </template>

    <n-empty v-if="items.length === 0" description="暂无数据" size="small" class="empty" />

    <div v-else class="rows">
      <div v-for="item in items" :key="item.name" class="row">
        <span class="name">{{ item.name }}</span>
        <n-progress
          type="line"
          :percentage="percentageOf(item.count)"
          :color="colorOf(item.name)"
          :height="8"
          :border-radius="4"
          :show-indicator="false"
          class="bar"
        />
        <span class="count">{{ item.count }}</span>
        <span class="percent">{{ percentageOf(item.count) }}%</span>
      </div>
    </div>
  </n-card>
</template>

<style lang="scss" scoped>
.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
}

.empty {
  padding: 24px 0;
}

.rows {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.row {
  display: flex;
  align-items: center;
  gap: 14px;
}

.name {
  width: 92px;
  flex-shrink: 0;
  font-size: 13px;
  color: #4b5563;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bar {
  flex: 1;
  // 进度条本身限宽，但外层仍撑满，数量才会贴到卡片右边成一条竖线
  :deep(.n-progress) {
    max-width: 520px;
  }
}

.percent {
  width: 38px;
  text-align: right;
  font-size: 12px;
  color: #8b95a7;
  font-variant-numeric: tabular-nums;
}

.count {
  width: 32px;
  text-align: right;
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  font-variant-numeric: tabular-nums;
}
</style>
