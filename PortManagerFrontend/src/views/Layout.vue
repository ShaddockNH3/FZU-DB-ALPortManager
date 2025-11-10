<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="header-left">
        <el-icon class="logo-icon" :size="32"><Ship /></el-icon>
        <h1 class="title">港区管理系统</h1>
      </div>
      <el-menu
        :default-active="activeMenu"
        mode="horizontal"
        :ellipsis="false"
        background-color="transparent"
        text-color="#fff"
        active-text-color="#ffd04b"
        @select="handleSelect"
      >
        <el-menu-item index="/ships">
          <el-icon><List /></el-icon>
          <span>舰船管理</span>
        </el-menu-item>
        <el-menu-item index="/statistics">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据统计</span>
        </el-menu-item>
      </el-menu>
    </el-header>
    <el-main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Ship, List, DataAnalysis } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const activeMenu = computed(() => route.path)

const handleSelect = (key) => {
  router.push(key)
}
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100%;
  background: transparent;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  padding: 0 30px;

  .header-left {
    display: flex;
    align-items: center;
    gap: 15px;

    .logo-icon {
      color: #fff;
    }

    .title {
      color: #fff;
      font-size: 24px;
      font-weight: bold;
      margin: 0;
      text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
    }
  }

  :deep(.el-menu) {
    border: none;
  }

  :deep(.el-menu-item) {
    border: none !important;
    font-size: 16px;
    
    &:hover {
      background-color: rgba(255, 255, 255, 0.1) !important;
    }
  }
}

.main-content {
  padding: 30px;
  overflow-y: auto;

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.3);
    border-radius: 4px;

    &:hover {
      background: rgba(255, 255, 255, 0.5);
    }
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
