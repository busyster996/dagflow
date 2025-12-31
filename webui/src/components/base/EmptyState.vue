<template>
  <div class="empty-state" :class="size">
    <div class="empty-icon-wrapper">
      <el-icon class="empty-icon" :size="iconSize">
        <component :is="icon" />
      </el-icon>
    </div>
    <h3 v-if="title" class="empty-title">{{ title }}</h3>
    <p class="empty-description">{{ description }}</p>
    <div v-if="$slots.actions" class="empty-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { FolderOpened } from '@element-plus/icons-vue'

const props = defineProps({
  // 图标组件
  icon: {
    type: Object,
    default: () => FolderOpened,
  },
  // 标题
  title: {
    type: String,
    default: '',
  },
  // 描述文本
  description: {
    type: String,
    default: '暂无数据',
  },
  // 大小
  size: {
    type: String,
    default: 'default',
    validator: (value) => ['small', 'default', 'large'].includes(value),
  },
})

/**
 * 图标大小
 */
const iconSize = computed(() => {
  const sizeMap = {
    small: 48,
    default: 64,
    large: 80,
  }
  return sizeMap[props.size] || 64
})
</script>

<style lang="scss" scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-3xl);
  gap: var(--spacing-base);
  color: var(--color-text-tertiary);
  min-height: 320px;
  
  &.small {
    min-height: 200px;
    padding: var(--spacing-xl);
    gap: var(--spacing-sm);
  }
  
  &.large {
    min-height: 400px;
    padding: var(--spacing-4xl);
    gap: var(--spacing-lg);
  }
  
  .empty-icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: var(--spacing-sm);
    
    .empty-icon {
      opacity: 0.4;
      color: var(--color-text-tertiary);
      transition: all var(--transition-base);
    }
  }
  
  .empty-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text-secondary);
    line-height: 1.5;
  }
  
  .empty-description {
    margin: 0;
    font-size: 14px;
    color: var(--color-text-tertiary);
    text-align: center;
    max-width: 480px;
    line-height: 1.6;
  }
  
  .empty-actions {
    margin-top: var(--spacing-base);
    display: flex;
    gap: var(--spacing-sm);
  }
  
  &:hover {
    .empty-icon {
      opacity: 0.6;
      transform: scale(1.05);
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .empty-state {
    padding: var(--spacing-xl);
    min-height: 240px;
    
    &.small {
      min-height: 180px;
      padding: var(--spacing-lg);
    }
    
    &.large {
      min-height: 300px;
      padding: var(--spacing-2xl);
    }
    
    .empty-title {
      font-size: 15px;
    }
    
    .empty-description {
      font-size: 13px;
      max-width: 100%;
    }
  }
}
</style>