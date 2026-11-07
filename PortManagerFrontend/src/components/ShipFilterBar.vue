<script setup>
import { Search, RefreshRight, Plus } from '@element-plus/icons-vue'
import { RARITY_OPTIONS, SHIP_TYPES, FACTIONS } from '@/constants/ship'

/**
 * 筛选条。model 是父级的 reactive 对象，这里直接改它的字段，
 * 改完发 search 让父级重新拉数据（下拉框 change 也即时触发）。
 */
defineProps({
  model: { type: Object, required: true },
  loading: { type: Boolean, default: false }
})

const emit = defineEmits(['search', 'reset', 'add'])
</script>

<template>
  <el-card class="filter-card" shadow="hover">
    <el-form :inline="true" :model="model" class="filter-form">
      <el-form-item label="舰船名称">
        <el-input
          v-model="model.shipName"
          placeholder="请输入舰船名称"
          clearable
          @clear="emit('search')"
        />
      </el-form-item>

      <el-form-item label="稀有度">
        <el-select
          v-model="model.rarity"
          placeholder="请选择稀有度"
          clearable
          class="select-rarity"
          @change="emit('search')"
        >
          <el-option
            v-for="item in RARITY_OPTIONS"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="舰种">
        <el-select
          v-model="model.shipType"
          placeholder="请选择舰种"
          clearable
          class="select-narrow"
          @change="emit('search')"
        >
          <el-option v-for="type in SHIP_TYPES" :key="type" :label="type" :value="type" />
        </el-select>
      </el-form-item>

      <el-form-item label="阵营">
        <el-select
          v-model="model.faction"
          placeholder="请选择阵营"
          clearable
          class="select-narrow"
          @change="emit('search')"
        >
          <el-option v-for="faction in FACTIONS" :key="faction" :label="faction" :value="faction" />
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :icon="Search" :loading="loading" @click="emit('search')">
          搜索
        </el-button>
        <el-button :icon="RefreshRight" @click="emit('reset')">重置</el-button>
        <el-button type="success" :icon="Plus" @click="emit('add')">添加舰船</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<style lang="scss" scoped>
.filter-card {
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.select-rarity {
  width: 180px;
}

.select-narrow {
  width: 150px;
}

.filter-form {
  row-gap: 8px;
}
</style>
