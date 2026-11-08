<script setup>
import { NButton, NCard, NForm, NFormItem, NIcon, NInput, NSelect } from 'naive-ui'
import { PhMagnifyingGlass, PhPlus, PhArrowCounterClockwise } from '@phosphor-icons/vue'
import { FACTIONS, RARITY_OPTIONS, SHIP_TYPES } from '@/constants/ship'

/**
 * 筛选条。model 是父级的 reactive 对象，这里直接改它的字段，
 * 改完发 search 让父级重新拉数据（下拉框 change 也即时触发）。
 */
defineProps({
  model: { type: Object, required: true },
  loading: { type: Boolean, default: false }
})

const emit = defineEmits(['search', 'reset', 'add'])

const rarityOptions = RARITY_OPTIONS.map((item) => ({ label: item.label, value: item.value }))
const shipTypeOptions = SHIP_TYPES.map((type) => ({ label: type, value: type }))
const factionOptions = FACTIONS.map((faction) => ({ label: faction, value: faction }))
</script>

<template>
  <n-card size="small" class="filter-card">
    <n-form inline label-placement="left" :show-feedback="false" class="filter-form">
      <n-form-item label="舰船名称">
        <n-input
          v-model:value="model.shipName"
          placeholder="搜索舰船名称"
          clearable
          class="input-name"
          @clear="emit('search')"
          @keyup.enter="emit('search')"
        />
      </n-form-item>

      <n-form-item label="稀有度">
        <n-select
          v-model:value="model.rarity"
          :options="rarityOptions"
          placeholder="全部"
          clearable
          class="select"
          @update:value="emit('search')"
        />
      </n-form-item>

      <n-form-item label="舰种">
        <n-select
          v-model:value="model.shipType"
          :options="shipTypeOptions"
          placeholder="全部"
          clearable
          class="select"
          @update:value="emit('search')"
        />
      </n-form-item>

      <n-form-item label="阵营">
        <n-select
          v-model:value="model.faction"
          :options="factionOptions"
          placeholder="全部"
          clearable
          class="select"
          @update:value="emit('search')"
        />
      </n-form-item>

      <n-form-item class="actions">
        <n-button type="primary" :loading="loading" @click="emit('search')">
          <template #icon>
            <n-icon><PhMagnifyingGlass /></n-icon>
          </template>
          搜索
        </n-button>

        <n-button quaternary @click="emit('reset')">
          <template #icon>
            <n-icon><PhArrowCounterClockwise /></n-icon>
          </template>
          重置
        </n-button>

        <n-button type="primary" secondary @click="emit('add')">
          <template #icon>
            <n-icon><PhPlus weight="bold" /></n-icon>
          </template>
          添加舰船
        </n-button>
      </n-form-item>
    </n-form>
  </n-card>
</template>

<style lang="scss" scoped>
.filter-card {
  margin-bottom: 16px;
}

// 一行放不下就换行，别把按钮挤到卡片外面
.filter-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;

  :deep(.n-form-item) {
    margin: 0;
  }

  :deep(.n-form-item-blank) {
    gap: 8px;
  }
}

.input-name {
  width: 200px;
}

.select {
  width: 148px;
}

.actions {
  margin-left: auto;
}

// 窄屏：每个字段独占一行，按钮回到左边
@media (max-width: 760px) {
  .filter-form :deep(.n-form-item) {
    width: 100%;
  }

  .input-name,
  .select {
    width: 100%;
  }

  .actions {
    margin-left: 0;
  }
}
</style>
