<template>
  <div class="home-container">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside :width="sidebarWidth" class="app-aside">
        <!-- Logo区域 -->
        <div class="logo-container">
          <transition name="fade" mode="out-in">
            <div v-if="!isCollapse" class="logo-expanded">
              <div class="logo-icon-wrapper">
                <span class="logo-emoji">🔄</span>
              </div>
              <div class="logo-text">
                <h1 class="logo-title">DagFlow</h1>
              </div>
            </div>
            <div v-else class="logo-collapsed">
              <span class="logo-emoji">🔄</span>
            </div>
          </transition>
        </div>

        <!-- 导航菜单 -->
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapse"
          :collapse-transition="false"
          class="side-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="/tasks">
            <el-icon><List /></el-icon>
            <template #title>任务管理</template>
          </el-menu-item>
          <el-menu-item index="/pipelines">
            <el-icon><Connection /></el-icon>
            <template #title>流水线管理</template>
          </el-menu-item>
          <el-menu-item index="/upload">
            <el-icon><Upload /></el-icon>
            <template #title>文件上传</template>
          </el-menu-item>
        </el-menu>

        <!-- 折叠按钮 -->
        <div class="collapse-trigger" @click="toggleCollapse">
          <el-icon class="trigger-icon">
            <DArrowLeft v-if="!isCollapse" />
            <DArrowRight v-else />
          </el-icon>
          <span v-if="!isCollapse" class="trigger-text">收起</span>
        </div>
      </el-aside>

      <!-- 主内容区 -->
      <el-container class="main-container">
        <!-- 顶部导航栏 -->
        <el-header class="app-header">
          <div class="header-left">
            <div class="breadcrumb-wrapper">
              <el-icon class="page-icon"><Menu /></el-icon>
              <h2 class="page-title">{{ pageTitle }}</h2>
            </div>
          </div>
          <div class="header-right">
            <div class="header-actions">
              <!-- 可以添加全局操作按钮 -->
            </div>
          </div>
        </el-header>

        <!-- 主要内容 -->
        <el-main class="app-main">
          <router-view v-slot="{ Component }">
            <transition name="slide-fade" mode="out-in">
              <component :is="Component" :key="$route.path" />
            </transition>
          </router-view>
        </el-main>
      </el-container>
    </el-container>

    <!-- 事件通知浮层 -->
    <div class="notification-container">
      <transition-group name="notification" tag="div">
        <div
          v-for="event in events"
          :key="event.id"
          class="notification-item"
          :class="`notification-${event.type || 'info'}`"
        >
          <div class="notification-icon">
            <el-icon><Bell /></el-icon>
          </div>
          <div class="notification-content">
            <p class="notification-message">{{ event.message }}</p>
          </div>
          <el-icon class="notification-close" @click="removeEvent(event.id)">
            <Close />
          </el-icon>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { API_BASE_URL, API_ENDPOINTS } from '@/config'
import { useSidebar } from '@/composables/useViewMode'

const router = useRouter()
const route = useRoute()

// 侧边栏管理
const { isCollapsed: isCollapse, toggle: toggleCollapse, width: sidebarWidth } = useSidebar()

const activeMenu = computed(() => route.path)
const events = ref([])
let eventSource = null
let eventIdCounter = 0

const pageTitle = computed(() => {
  const titles = {
    '/tasks': '任务管理',
    '/pipelines': '流水线管理',
    '/upload': '文件上传',
  }
  return titles[route.path] || 'DagFlow'
})

const handleMenuSelect = (index) => {
  router.push(index)
}

const removeEvent = (id) => {
  const index = events.value.findIndex(e => e.id === id)
  if (index > -1) {
    events.value.splice(index, 1)
  }
}

// SSE 事件监听
const initEventListener = () => {
  eventSource = new EventSource(`${API_BASE_URL}${API_ENDPOINTS.event}`)

  eventSource.onmessage = (event) => {
    const eventObj = {
      id: eventIdCounter++,
      message: event.data,
      type: 'info',
      timestamp: Date.now(),
    }
    events.value.unshift(eventObj)

    // 保留最近5条
    if (events.value.length > 5) {
      events.value.pop()
    }

    // 5秒后自动移除
    setTimeout(() => {
      removeEvent(eventObj.id)
    }, 5000)
  }

  eventSource.onerror = (error) => {
    console.error('SSE连接错误:', error)
  }

  eventSource.onopen = () => {
  }
}

onMounted(() => {
  initEventListener()
})

onUnmounted(() => {
  if (eventSource) {
    eventSource.close()
  }
})
</script>

<style lang="scss" scoped>
.home-container {
  height: 100vh;
  background: var(--color-bg-secondary);
  overflow: hidden;
}

.el-container {
  height: 100%;
}

// ==========================================
// 侧边栏样式
// ==========================================
.app-aside {
  background: linear-gradient(180deg, #1e293b 0%, #0f172a 100%);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-fast);
  box-shadow: var(--shadow-sm);
  position: relative;
  z-index: var(--z-fixed);
  border-right: 1px solid rgba(255, 255, 255, 0.05);
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-base);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.02);
}

.logo-expanded {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  width: 100%;

  .logo-icon-wrapper {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);

    .logo-emoji {
      font-size: 20px;
      filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
    }
  }

  .logo-text {
    flex: 1;
    min-width: 0;

    .logo-title {
      margin: 0;
      font-size: 18px;
      font-weight: 700;
      color: white;
      line-height: 1.3;
      letter-spacing: 0.5px;
      background: linear-gradient(135deg, #ffffff 0%, #a0aec0 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .logo-subtitle {
      margin: 2px 0 0 0;
      font-size: 12px;
      color: rgba(255, 255, 255, 0.6);
      font-weight: 500;
      letter-spacing: 0.5px;
    }
  }
}

.logo-collapsed {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);

  .logo-emoji {
    font-size: 20px;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
  }
}

.side-menu {
  flex: 1;
  border-right: none;
  background: transparent;
  overflow-y: auto;
  overflow-x: hidden;
  padding: var(--spacing-sm) 0;

  :deep(.el-menu-item) {
    margin: 0 var(--spacing-sm) var(--spacing-xs) var(--spacing-sm);
    padding: 0 var(--spacing-base);
    height: 44px;
    line-height: 44px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 14px;
    font-weight: 500;
    border-radius: var(--radius-md);
    transition: all var(--transition-base);

    &:hover {
      color: white;
      background: rgba(255, 255, 255, 0.1);
      transform: translateX(4px);
    }

    &.is-active {
      color: white;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      box-shadow: var(--shadow-md);
      transform: translateX(0);

      &::before {
        content: '';
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
        width: 3px;
        height: 50%;
        background: white;
        border-radius: 0 var(--radius-base) var(--radius-base) 0;
      }
    }

    .el-icon {
      margin-right: var(--spacing-sm);
      font-size: 18px;
      color: inherit;
    }
  }

  &::-webkit-scrollbar {
    width: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: var(--radius-base);
  }
}

.collapse-trigger {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  color: rgba(255, 255, 255, 0.6);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.02);
  transition: all var(--transition-base);
  font-size: 13px;
  font-weight: 500;

  &:hover {
    color: white;
    background: rgba(255, 255, 255, 0.08);

    .trigger-icon {
      transform: scale(1.1);
    }
  }

  .trigger-icon {
    font-size: 16px;
    transition: transform var(--transition-base);
  }
}

// ==========================================
// 主容器样式
// ==========================================
.main-container {
  background: var(--color-bg-secondary);
  overflow: hidden;
}

.app-header {
  background: white;
  border-bottom: 1px solid var(--color-border-light);
  padding: 0 var(--spacing-lg);
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-sm);
  height: 56px;
  z-index: var(--z-sticky);
  position: relative;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);

  .breadcrumb-wrapper {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);

    .page-icon {
      font-size: 20px;
      color: var(--color-primary);
    }

    .page-title {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: var(--color-text-primary);
      letter-spacing: 0.3px;
      line-height: 1.4;
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.app-main {
  padding: 0;
  overflow: hidden;
  background: var(--color-bg-secondary);
  position: relative;
}

// ==========================================
// 通知浮层样式
// ==========================================
.notification-container {
  position: fixed;
  top: 64px;
  right: var(--spacing-lg);
  z-index: var(--z-tooltip);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  pointer-events: none;
  max-width: 400px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-sm);
  padding: var(--spacing-base) var(--spacing-lg);
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  border-left: 4px solid var(--color-primary);
  pointer-events: auto;
  transition: all var(--transition-base);
  cursor: pointer;
  min-width: 320px;

  &:hover {
    transform: translateX(-4px);
    box-shadow: var(--shadow-xl);
  }

  &.notification-success {
    border-left-color: var(--color-success);

    .notification-icon {
      color: var(--color-success);
    }
  }

  &.notification-warning {
    border-left-color: var(--color-warning);

    .notification-icon {
      color: var(--color-warning);
    }
  }

  &.notification-danger,
  &.notification-error {
    border-left-color: var(--color-danger);

    .notification-icon {
      color: var(--color-danger);
    }
  }

  .notification-icon {
    flex-shrink: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-primary);
    font-size: 18px;
  }

  .notification-content {
    flex: 1;
    min-width: 0;

    .notification-message {
      margin: 0;
      font-size: 14px;
      color: var(--color-text-primary);
      line-height: 1.6;
      word-break: break-word;
    }
  }

  .notification-close {
    flex-shrink: 0;
    width: 20px;
    height: 20px;
    cursor: pointer;
    color: var(--color-text-tertiary);
    transition: color var(--transition-base);
    font-size: 16px;

    &:hover {
      color: var(--color-text-primary);
    }
  }
}

// 通知动画
.notification-enter-active {
  animation: slideInRight var(--transition-fast);
}

.notification-leave-active {
  animation: slideOutRight var(--transition-fast);
}

@keyframes slideInRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slideOutRight {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 1200px) {
  .notification-container {
    right: var(--spacing-sm);
    max-width: 300px;
  }

  .notification-item {
    min-width: 260px;
  }
}

@media (max-width: 768px) {
  .logo-container {
    height: 56px;
    padding: var(--spacing-sm);
  }
  
  .logo-expanded {
    .logo-icon-wrapper {
      width: 32px;
      height: 32px;
      
      .logo-emoji {
        font-size: 18px;
      }
    }
    
    .logo-text {
      .logo-title {
        font-size: 16px;
      }
    }
  }
  
  .side-menu {
    padding: var(--spacing-sm) 0;
    
    :deep(.el-menu-item) {
      margin: 0 var(--spacing-sm) var(--spacing-xs) var(--spacing-sm);
      height: 40px;
      line-height: 40px;
      font-size: 13px;
      
      .el-icon {
        font-size: 16px;
      }
    }
  }
  
  .collapse-trigger {
    height: 44px;
    font-size: 12px;
    
    .trigger-icon {
      font-size: 14px;
    }
  }

  .app-aside {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: var(--z-modal);
    box-shadow: var(--shadow-xl);
  }

  .app-header {
    padding: 0 var(--spacing-base);
    height: 52px;
    
    .page-icon {
      font-size: 18px !important;
    }
    
    .page-title {
      font-size: 16px !important;
    }
  }

  .notification-container {
    top: 60px;
    right: var(--spacing-base);
    left: var(--spacing-base);
    max-width: none;
  }

  .notification-item {
    min-width: 0;
    padding: var(--spacing-sm) var(--spacing-base);
    
    .notification-icon {
      width: 20px;
      height: 20px;
      font-size: 16px;
    }
    
    .notification-content {
      .notification-message {
        font-size: 13px;
      }
    }
    
    .notification-close {
      width: 18px;
      height: 18px;
      font-size: 14px;
    }
  }
}
</style>