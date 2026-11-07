<script setup>
import { computed } from 'vue'

/**
 * 分布卡片：标题 + 一组「名称 / 进度条 / 数量」。
 * 阵营、稀有度、舰种三块原来各写一遍同样的模板。
 */
const props = defineProps({
  title: { type: String, required: true },
  items: { type: Array, default: () => [] },
  total: { type: Number, default: 0 },
  /** 传颜色或 (name) => 颜色 都行 */
  color: { type: [String, Function], default: '#409eff' }
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
  <el-card class="chart-card" shadow="hover">
    <template #header>
      <div class="card-header">
        <span class="card-title">{{ title }}</span>
      </div>
    </template>

    <el-empty v-if="items.length === 0" description="暂无数据" :image-size="60" />

    <div v-else class="stat-list">
      <div v-for="item in items" :key="item.name" class="stat-item">
        <div class="item-label">{{ item.name }}</div>
        <div class="item-bar">
          <el-progress :percentage="percentageOf(item.count)" :color="colorOf(item.name)" />
        </div>
        <div class="item-value">{{ item.count }}</div>
      </div>
    </div>
  </el-card>
</template>

<style lang="scss" scoped>
.chart-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: none;
}

.card-title {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
}

.item-label {
  width: 100px;
  font-size: 14px;
  color: #606266;
  white-space: nowrap;
}

.item-bar {
  flex: 1;
}

.item-value {
  width: 50px;
  text-align: right;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}
</style>
