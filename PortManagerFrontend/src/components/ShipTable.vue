<script setup>
import { Edit, Delete, Plus, Minus } from '@element-plus/icons-vue'
import { initialStars, maxStars, rarityTagType } from '@/constants/ship'

/**
 * 舰船表格 + 分页。纯展示，所有动作都往上抛事件，状态在父级。
 */
defineProps({
  ships: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  total: { type: Number, default: 0 }
})

const emit = defineEmits(['edit', 'delete', 'equip', 'star', 'page-change', 'size-change'])
</script>

<template>
  <el-card class="table-card" shadow="hover">
    <el-table
      :data="ships"
      v-loading="loading"
      stripe
      style="width: 100%"
      :default-sort="{ prop: 'id', order: 'descending' }"
    >
      <el-table-column prop="id" label="ID" width="80" sortable />
      <el-table-column prop="shipName" label="舰船名称" min-width="150" />

      <el-table-column prop="rarity" label="稀有度" width="150">
        <template #default="{ row }">
          <el-tag :type="rarityTagType(row.rarity)">{{ row.rarity }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="shipType" label="舰种" width="120" />
      <el-table-column prop="faction" label="阵营" width="150" />

      <el-table-column label="星级" width="180">
        <template #default="{ row }">
          <div class="star-cell">
            <el-button
              size="small"
              :icon="Minus"
              circle
              :disabled="row.stars <= initialStars(row.rarity)"
              @click="emit('star', row, -1)"
            />
            <span class="star-count">
              {{ row.stars }}★ {{ row.stars >= maxStars(row.rarity) ? '(满)' : '' }}
            </span>
            <el-button
              size="small"
              :icon="Plus"
              circle
              :disabled="row.stars >= maxStars(row.rarity)"
              @click="emit('star', row, 1)"
            />
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="level" label="等级" width="100" sortable />
      <el-table-column prop="createdAt" label="创建时间" width="180" />

      <el-table-column label="操作" width="300" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" :icon="Edit" @click="emit('edit', row)">
            编辑
          </el-button>
          <el-button type="success" size="small" @click="emit('equip', row)">装备</el-button>
          <el-button type="danger" size="small" :icon="Delete" @click="emit('delete', row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pagination"
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @current-change="(value) => emit('page-change', value)"
      @size-change="(value) => emit('size-change', value)"
    />
  </el-card>
</template>

<style lang="scss" scoped>
.table-card {
  background: rgba(255, 255, 255, 0.95);
}

.star-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.star-count {
  min-width: 80px;
  text-align: center;
}

.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>
