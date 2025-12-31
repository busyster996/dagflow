/**
 * 视图模式管理 Composable
 */

import { ref, watch } from 'vue'

/**
 * 视图模式管理（卡片/列表切换）
 * @param {string} defaultMode - 默认模式，'grid' 或 'table'
 * @param {string} storageKey - 本地存储键名
 * @returns {Object} 视图模式管理对象
 */
export function useViewMode(defaultMode = 'grid', storageKey = null) {
  // 从本地存储恢复视图模式
  const storedMode = storageKey ? localStorage.getItem(storageKey) : null
  const viewMode = ref(storedMode || defaultMode)
  
  /**
   * 切换视图模式
   * @param {string} mode - 新的视图模式
   */
  const setViewMode = (mode) => {
    viewMode.value = mode
  }
  
  /**
   * 切换到卡片视图
   */
  const toGridView = () => {
    setViewMode('grid')
  }
  
  /**
   * 切换到列表视图
   */
  const toTableView = () => {
    setViewMode('table')
  }
  
  /**
   * 切换视图模式
   */
  const toggleViewMode = () => {
    viewMode.value = viewMode.value === 'grid' ? 'table' : 'grid'
  }
  
  /**
   * 是否为卡片视图
   */
  const isGridView = computed(() => viewMode.value === 'grid')
  
  /**
   * 是否为列表视图
   */
  const isTableView = computed(() => viewMode.value === 'table')
  
  // 监听变化并保存到本地存储
  if (storageKey) {
    watch(viewMode, (newMode) => {
      localStorage.setItem(storageKey, newMode)
    })
  }
  
  return {
    viewMode,
    setViewMode,
    toGridView,
    toTableView,
    toggleViewMode,
    isGridView,
    isTableView,
  }
}

/**
 * 侧边栏折叠状态管理
 * @param {boolean} defaultCollapsed - 默认是否折叠
 * @param {string} storageKey - 本地存储键名
 * @returns {Object} 侧边栏管理对象
 */
export function useSidebar(defaultCollapsed = false, storageKey = 'sidebar-collapsed') {
  // 从本地存储恢复折叠状态
  const storedState = localStorage.getItem(storageKey)
  const isCollapsed = ref(storedState ? storedState === 'true' : defaultCollapsed)
  
  /**
   * 切换折叠状态
   */
  const toggle = () => {
    isCollapsed.value = !isCollapsed.value
  }
  
  /**
   * 展开侧边栏
   */
  const expand = () => {
    isCollapsed.value = false
  }
  
  /**
   * 折叠侧边栏
   */
  const collapse = () => {
    isCollapsed.value = true
  }
  
  /**
   * 侧边栏宽度
   */
  const width = computed(() => {
    return isCollapsed.value ? '48px' : '160px'
  })
  
  // 监听变化并保存到本地存储
  watch(isCollapsed, (newValue) => {
    localStorage.setItem(storageKey, newValue.toString())
  })
  
  return {
    isCollapsed,
    toggle,
    expand,
    collapse,
    width,
  }
}

/**
 * 标签页管理
 * @param {string} defaultTab - 默认激活的标签页
 * @returns {Object} 标签页管理对象
 */
export function useTabs(defaultTab = 'default') {
  const activeTab = ref(defaultTab)
  const tabHistory = ref([defaultTab])
  
  /**
   * 切换标签页
   * @param {string} tab - 标签页名称
   */
  const switchTab = (tab) => {
    if (activeTab.value !== tab) {
      activeTab.value = tab
      tabHistory.value.push(tab)
      
      // 保留最近10个历史记录
      if (tabHistory.value.length > 10) {
        tabHistory.value.shift()
      }
    }
  }
  
  /**
   * 返回上一个标签页
   */
  const goBack = () => {
    if (tabHistory.value.length > 1) {
      tabHistory.value.pop() // 移除当前
      const previousTab = tabHistory.value[tabHistory.value.length - 1]
      activeTab.value = previousTab
    }
  }
  
  /**
   * 检查是否为激活标签页
   * @param {string} tab - 标签页名称
   */
  const isActive = (tab) => {
    return activeTab.value === tab
  }
  
  /**
   * 重置到默认标签页
   */
  const reset = () => {
    activeTab.value = defaultTab
    tabHistory.value = [defaultTab]
  }
  
  return {
    activeTab,
    tabHistory,
    switchTab,
    goBack,
    isActive,
    reset,
  }
}