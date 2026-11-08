<script setup>
import { onMounted, ref } from 'vue'
import { NGi, NGrid } from 'naive-ui'
import { PhBoat, PhFlag, PhStack, PhStar } from '@phosphor-icons/vue'
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
  <div class="page">
    <header class="page-head">
      <h2>数据统计</h2>
      <p>按阵营、稀有度、舰种三个维度看舰队的构成</p>
    </header>

    <n-grid cols="1 s:2 m:4" responsive="screen" :x-gap="16" :y-gap="16" item-responsive>
      <n-gi>
        <stat-card
          label="总舰船数"
          :value="statistics.totalShips"
          :icon="PhBoat"
          color="#1e40af"
        />
      </n-gi>
      <n-gi>
        <stat-card
          label="阵营数量"
          :value="statistics.factionStats?.length || 0"
          :icon="PhFlag"
          color="#2f9e6f"
        />
      </n-gi>
      <n-gi>
        <stat-card
          label="稀有度种类"
          :value="statistics.rarityStats?.length || 0"
          :icon="PhStar"
          color="#c8963e"
        />
      </n-gi>
      <n-gi>
        <stat-card
          label="舰种数量"
          :value="statistics.shipTypeStats?.length || 0"
          :icon="PhStack"
          color="#6b5bd2"
        />
      </n-gi>
    </n-grid>

    <n-grid cols="1 m:2" responsive="screen" :x-gap="16" :y-gap="16" item-responsive class="charts">
      <n-gi>
        <distribution-card
          title="阵营分布"
          :items="statistics.factionStats"
          :total="statistics.totalShips"
          :color="factionColor"
        />
      </n-gi>
      <n-gi>
        <distribution-card
          title="稀有度分布"
          :items="statistics.rarityStats"
          :total="statistics.totalShips"
          :color="rarityColor"
        />
      </n-gi>
      <n-gi span="1 m:2">
        <distribution-card
          title="舰种分布"
          :items="statistics.shipTypeStats"
          :total="statistics.totalShips"
          color="#1e40af"
        />
      </n-gi>
    </n-grid>
  </div>
</template>

<style lang="scss" scoped>
.page-head {
  margin-bottom: 18px;

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: #1f2937;
  }

  p {
    margin-top: 4px;
    font-size: 13px;
    color: #8b95a7;
  }
}

.charts {
  margin-top: 16px;
}
</style>
