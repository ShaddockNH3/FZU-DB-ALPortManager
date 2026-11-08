<script setup>
import { computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NIcon, NLayout, NLayoutContent, NLayoutHeader, NMenu } from 'naive-ui'
import { PhAnchor } from '@phosphor-icons/vue'

const router = useRouter()
const route = useRoute()

const activeKey = computed(() => route.path)

/**
 * 菜单项从路由表里取：标题和图标都在 meta 上，加页面只改 router/index.js。
 */
const menuOptions = computed(() => {
  const layout = router.options.routes.find((item) => item.path === '/')
  return (layout?.children || [])
    .filter((child) => child.meta?.title)
    .map((child) => ({
      key: child.path,
      label: child.meta.title,
      icon: child.meta.icon
        ? () => h(NIcon, null, { default: () => h(child.meta.icon) })
        : undefined
    }))
})

function go(key) {
  router.push(key)
}
</script>

<template>
  <n-layout class="shell">
    <n-layout-header bordered class="header">
      <div class="header-inner">
        <div class="brand">
          <n-icon :size="26" color="#1e40af"><PhAnchor weight="duotone" /></n-icon>
          <div class="brand-text">
            <h1>港区管理系统</h1>
            <span>AL Port Manager</span>
          </div>
        </div>

        <n-menu
          mode="horizontal"
          :value="activeKey"
          :options="menuOptions"
          :responsive="false"
          @update:value="go"
        />
      </div>
    </n-layout-header>

    <n-layout-content class="content">
      <div class="content-inner">
        <router-view />
      </div>
    </n-layout-content>
  </n-layout>
</template>

<style lang="scss" scoped>
.shell {
  min-height: 100vh;
  background: #f4f6fa;
}

.header {
  background: #fff;
}

.header-inner {
  max-width: 1320px;
  margin: 0 auto;
  height: 64px;
  padding: 0 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 32px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;

  h1 {
    font-size: 17px;
    font-weight: 600;
    color: #1f2937;
    letter-spacing: 0.2px;
  }

  span {
    font-size: 11px;
    color: #8b95a7;
    letter-spacing: 1.4px;
    text-transform: uppercase;
  }
}

.content {
  background: #f4f6fa;
}

.content-inner {
  max-width: 1320px;
  margin: 0 auto;
  padding: 24px 28px 40px;
}

// 手机上把留白收一收，不然表格挤得只剩两条
@media (max-width: 760px) {
  .header-inner {
    padding: 0 12px;
    gap: 12px;
  }

  .brand-text span {
    display: none;
  }

  .brand-text h1 {
    font-size: 15px;
  }

  .content-inner {
    padding: 16px 12px 32px;
  }
}

// 更窄的时候菜单只留图标，不然文字被压成「航…」
@media (max-width: 520px) {
  :deep(.n-menu-item-content-header) {
    display: none;
  }

  :deep(.n-menu-item-content) {
    padding-left: 10px;
    padding-right: 10px;
  }
}
</style>
