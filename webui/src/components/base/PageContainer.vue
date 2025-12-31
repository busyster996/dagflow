<template>
  <div class="page-container" :class="{ compact, fullHeight }">
    <!-- 工具栏槽位 -->
    <div v-if="$slots.toolbar" class="page-toolbar">
      <slot name="toolbar" />
    </div>
    
    <!-- 统计卡片槽位 -->
    <div v-if="$slots.stats" class="page-stats">
      <slot name="stats" />
    </div>
    
    <!-- 主要内容区 -->
    <div class="page-content" :class="{ scrollable }">
      <slot />
    </div>
    
    <!-- 页脚槽位（通常用于分页） -->
    <div v-if="$slots.footer" class="page-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  // 紧凑模式（减少间距）
  compact: {
    type: Boolean,
    default: false,
  },
  // 占满高度
  fullHeight: {
    type: Boolean,
    default: true,
  },
  // 内容区是否可滚动
  scrollable: {
    type: Boolean,
    default: true,
  },
})
</script>

<style lang="scss" scoped>
.page-container {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-base);
  padding: var(--spacing-lg);
  overflow: hidden;
  
  &.fullHeight {
    height: 100%;
  }
  
  &.compact {
    gap: var(--spacing-sm);
    padding: var(--spacing-base);
  }
  
  .page-toolbar {
    flex-shrink: 0;
  }
  
  .page-stats {
    flex-shrink: 0;
  }
  
  .page-content {
    flex: 1;
    min-height: 0;
    
    &.scrollable {
      overflow-y: auto;
      overflow-x: hidden;
      padding-right: var(--spacing-xs);
    }
  }
  
  .page-footer {
    flex-shrink: 0;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .page-container {
    gap: var(--spacing-sm);
    padding: var(--spacing-base);
    
    &.compact {
      gap: var(--spacing-xs);
      padding: var(--spacing-sm);
    }
  }
}
</style>