<script setup>
import { onMounted, ref } from 'vue'
import { Flag, Grid, Ship, Star } from '@element-plus/icons-vue'
import { getStatistics } from '@/api/ship'
import { factionColor, rarityColor } from '@/constants/ship'
import StatCard from '@/components/StatCard.vue'
import DistributionCard from '@/components/DistributionCard.vue'

const statistics = ref({
  totalShips: 0,
  factionStats: [],
  rarityStats: [],
  shipTypeStats: []
})

async function loadStatistics() {
  try {
    statistics.value = await getStatistics()
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(loadStatistics)
</script>

<template>
  <div class="statistics-container">
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :md="6">
        <stat-card label="总舰船数" :value="statistics.totalShips" :icon="Ship" tone="total" />
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <stat-card
          label="阵营数量"
          :value="statistics.factionStats?.length || 0"
          :icon="Flag"
          tone="faction"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <stat-card
          label="稀有度种类"
          :value="statistics.rarityStats?.length || 0"
          :icon="Star"
          tone="rarity"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <stat-card
          label="舰种数量"
          :value="statistics.shipTypeStats?.length || 0"
          :icon="Grid"
          tone="type"
        />
      </el-col>
    </el-row>

    <el-row :gutter="20" class="second-row">
      <el-col :xs="24" :md="12">
        <distribution-card
          title="阵营分布"
          :items="statistics.factionStats"
          :total="statistics.totalShips"
          :color="factionColor"
        />
      </el-col>
      <el-col :xs="24" :md="12">
        <distribution-card
          title="稀有度分布"
          :items="statistics.rarityStats"
          :total="statistics.totalShips"
          :color="rarityColor"
        />
      </el-col>
    </el-row>

    <el-row :gutter="20" class="second-row">
      <el-col :span="24">
        <distribution-card
          title="舰种分布"
          :items="statistics.shipTypeStats"
          :total="statistics.totalShips"
          color="#409eff"
        />
      </el-col>
    </el-row>
  </div>
</template>

<style lang="scss" scoped>
.second-row {
  margin-top: 20px;
}
</style>
