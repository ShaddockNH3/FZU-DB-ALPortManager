<script setup>
import { NCard, NIcon, NStatistic } from 'naive-ui'

/**
 * 顶部指标卡。四张卡只有图标、文案、数值和颜色不同，模板只留一份。
 */
defineProps({
  label: { type: String, required: true },
  value: { type: [Number, String], default: 0 },
  icon: { type: Object, required: true },
  color: { type: String, default: '#1e40af' },
  unit: { type: String, default: '' }
})
</script>

<template>
  <n-card size="small" class="stat-card">
    <div class="row">
      <div class="badge" :style="{ background: `${color}14`, color }">
        <n-icon :size="22"><component :is="icon" weight="duotone" /></n-icon>
      </div>

      <n-statistic :label="label" :value="value">
        <template v-if="unit" #suffix>
          <span class="unit">{{ unit }}</span>
        </template>
      </n-statistic>
    </div>
  </n-card>
</template>

<style lang="scss" scoped>
.row {
  display: flex;
  align-items: center;
  gap: 14px;
}

.badge {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.unit {
  font-size: 13px;
  color: #8b95a7;
  margin-left: 2px;
}
</style>
