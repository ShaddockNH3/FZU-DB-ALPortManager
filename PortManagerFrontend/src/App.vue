<template>
  <div id="app" class="app-container">
    <router-view />
  </div>
</template>

<!-- 全局重置与底色在 assets/styles/main.scss，这里只留骨架 -->
<style lang="scss" scoped>
.app-container {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
</style>
