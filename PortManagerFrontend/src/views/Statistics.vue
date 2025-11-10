<template>
  <div class="statistics-container">
    <el-row :gutter="20">
      <!-- 总舰船数 -->
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card total-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="48"><Ship /></el-icon>
            <div class="stat-info">
              <div class="stat-label">总舰船数</div>
              <div class="stat-value">{{ statistics.totalShips }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 阵营数 -->
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card faction-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="48"><Flag /></el-icon>
            <div class="stat-info">
              <div class="stat-label">阵营数量</div>
              <div class="stat-value">{{ statistics.factionStats?.length || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 稀有度种类 -->
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card rarity-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="48"><Star /></el-icon>
            <div class="stat-info">
              <div class="stat-label">稀有度种类</div>
              <div class="stat-value">{{ statistics.rarityStats?.length || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 舰种数 -->
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card type-card" shadow="hover">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="48"><Grid /></el-icon>
            <div class="stat-info">
              <div class="stat-label">舰种数量</div>
              <div class="stat-value">{{ statistics.shipTypeStats?.length || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 阵营统计 -->
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">阵营分布</span>
            </div>
          </template>
          <div class="stat-list">
            <div
              v-for="item in statistics.factionStats"
              :key="item.name"
              class="stat-item"
            >
              <div class="item-label">{{ item.name }}</div>
              <div class="item-bar">
                <el-progress
                  :percentage="getPercentage(item.count, statistics.totalShips)"
                  :color="getFactionColor(item.name)"
                />
              </div>
              <div class="item-value">{{ item.count }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 稀有度统计 -->
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">稀有度分布</span>
            </div>
          </template>
          <div class="stat-list">
            <div
              v-for="item in statistics.rarityStats"
              :key="item.name"
              class="stat-item"
            >
              <div class="item-label">{{ item.name }}</div>
              <div class="item-bar">
                <el-progress
                  :percentage="getPercentage(item.count, statistics.totalShips)"
                  :color="getRarityColor(item.name)"
                />
              </div>
              <div class="item-value">{{ item.count }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 舰种统计 -->
      <el-col :span="24">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">舰种分布</span>
            </div>
          </template>
          <div class="stat-list">
            <div
              v-for="item in statistics.shipTypeStats"
              :key="item.name"
              class="stat-item"
            >
              <div class="item-label">{{ item.name }}</div>
              <div class="item-bar">
                <el-progress
                  :percentage="getPercentage(item.count, statistics.totalShips)"
                  color="#409eff"
                />
              </div>
              <div class="item-value">{{ item.count }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Ship, Flag, Star, Grid } from '@element-plus/icons-vue'
import { getStatistics } from '@/api/ship'

const statistics = ref({
  totalShips: 0,
  factionStats: [],
  rarityStats: [],
  shipTypeStats: []
})

// 计算百分比
const getPercentage = (count, total) => {
  if (!total) return 0
  return Math.round((count / total) * 100)
}

// 获取阵营颜色
const getFactionColor = (faction) => {
  const colorMap = {
    '白鹰': '#4169E1',
    '皇家': '#FFD700',
    '重樱': '#FF1493',
    '铁血': '#696969',
    '东煌': '#FF4500',
    '北方联合': '#DC143C',
    '自由鸢尾': '#9370DB',
    '维希教廷': '#8B4513',
    '鸢尾教国': '#BA55D3',
    '撒丁帝国': '#228B22',
    '郁金王国': '#FFB6C1',
    '飓风': '#00CED1',
    'META': '#8B008B'
  }
  return colorMap[faction] || '#409eff'
}

// 获取稀有度颜色
const getRarityColor = (rarity) => {
  const colorMap = {
    '普通': '#909399',
    '稀有': '#67C23A',
    '精锐': '#E6A23C',
    '超稀有': '#F56C6C',
    '海上传奇': '#C71585',
    '最高方案': '#8B008B',
    '决战方案': '#FF1493'
  }
  return colorMap[rarity] || '#409eff'
}

// 加载统计数据
const loadStatistics = async () => {
  try {
    const data = await getStatistics()
    statistics.value = data
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(() => {
  loadStatistics()
})
</script>

<style lang="scss" scoped>
.statistics-container {
  .stat-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: none;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-5px);
    }

    .stat-content {
      display: flex;
      align-items: center;
      gap: 20px;

      .stat-icon {
        color: #fff;
        padding: 15px;
        border-radius: 12px;
      }

      .stat-info {
        flex: 1;

        .stat-label {
          font-size: 14px;
          color: #909399;
          margin-bottom: 8px;
        }

        .stat-value {
          font-size: 28px;
          font-weight: bold;
          color: #303133;
        }
      }
    }
  }

  .total-card .stat-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .faction-card .stat-icon {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  }

  .rarity-card .stat-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }

  .type-card .stat-icon {
    background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  }

  .chart-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: none;

    .card-header {
      .card-title {
        font-size: 18px;
        font-weight: bold;
        color: #303133;
      }
    }

    .stat-list {
      .stat-item {
        display: flex;
        align-items: center;
        gap: 15px;
        margin-bottom: 20px;

        &:last-child {
          margin-bottom: 0;
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
      }
    }
  }
}
</style>
