<template>
  <div 
    class="card-grid" 
    :class="[`cols-${columns}`, { compact, bordered }]"
    :style="gridStyle"
  >
    <slot />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // 列数（响应式）
  columns: {
    type: [Number, String, Object],
    default: 'auto',
    validator: (value) => {
      if (typeof value === 'number') return value > 0 && value <= 12
      if (typeof value === 'string') return value === 'auto'
      if (typeof value === 'object') {
        // 支持响应式列数配置
        return true
      }
      return false
    },
  },
  // 最小列宽（当 columns 为 'auto' 时生效）
  minColumnWidth: {
    type: [String, Number],
    default: 340,
  },
  // 间距
  gap: {
    type: String,
    default: 'var(--spacing-base)',
  },
  // 紧凑模式
  compact: {
    type: Boolean,
    default: false,
  },
  // 显示边框
  bordered: {
    type: Boolean,
    default: false,
  },
})

/**
 * 网格样式
 */
const gridStyle = computed(() => {
  let gridTemplateColumns
  
  if (props.columns === 'auto') {
    const minWidth = typeof props.minColumnWidth === 'number' 
      ? `${props.minColumnWidth}px` 
      : props.minColumnWidth
    gridTemplateColumns = `repeat(auto-fill, minmax(${minWidth}, 1fr))`
  } else if (typeof props.columns === 'number') {
    gridTemplateColumns = `repeat(${props.columns}, 1fr)`
  } else if (typeof props.columns === 'object') {
    // 响应式列数配置
    gridTemplateColumns = `repeat(auto-fill, minmax(${props.minColumnWidth}px, 1fr))`
  }
  
  return {
    gridTemplateColumns,
    gap: props.gap,
  }
})
</script>

<style lang="scss" scoped>
.card-grid {
  display: grid;
  padding: 0;
  
  &.compact {
    gap: var(--spacing-sm);
    padding: 0;
  }
  
  &.bordered {
    :deep(> *) {
      border: 1px solid var(--color-border-light);
    }
  }
}

// 响应式设计
@media (max-width: 1400px) {
  .card-grid {
    &:not(.compact) {
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)) !important;
    }
  }
}

@media (max-width: 768px) {
  .card-grid {
    grid-template-columns: 1fr !important;
    gap: var(--spacing-sm);
    
    &.compact {
      gap: var(--spacing-xs);
    }
  }
}
</style>