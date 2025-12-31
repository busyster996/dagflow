<template>
  <div class="info-item" :class="[layout, { hoverable }]">
    <div v-if="icon" class="item-icon">
      <el-icon :size="iconSize">
        <component :is="icon" />
      </el-icon>
    </div>
    <div class="item-content">
      <span class="item-label">{{ label }}</span>
      <span v-if="!$slots.value" class="item-value" :class="valueClass">
        {{ value }}
      </span>
      <div v-else class="item-value">
        <slot name="value" />
      </div>
    </div>
    <div v-if="$slots.action" class="item-action">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  // 标签文本
  label: {
    type: String,
    required: true,
  },
  // 值
  value: {
    type: [String, Number],
    default: '',
  },
  // 图标组件
  icon: {
    type: Object,
    default: null,
  },
  // 图标大小
  iconSize: {
    type: Number,
    default: 16,
  },
  // 布局方式：horizontal, vertical
  layout: {
    type: String,
    default: 'horizontal',
    validator: (value) => ['horizontal', 'vertical'].includes(value),
  },
  // 值的样式类
  valueClass: {
    type: String,
    default: '',
  },
  // 是否可悬浮
  hoverable: {
    type: Boolean,
    default: false,
  },
})
</script>

<style lang="scss" scoped>
.info-item {
  display: flex;
  padding: var(--spacing-sm) var(--spacing-base);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-md);
  transition: all var(--transition-base);
  border: 1px solid transparent;
  min-height: 44px;
  
  &.hoverable:hover {
    background: var(--color-bg-secondary);
    transform: translateX(2px);
    border-color: var(--color-border-base);
    box-shadow: var(--shadow-sm);
  }
  
  &.horizontal {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    gap: var(--spacing-base);
    
    .item-content {
      flex-direction: row;
      align-items: center;
      justify-content: space-between;
      flex: 1;
    }
  }
  
  &.vertical {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
    
    .item-content {
      flex-direction: column;
      align-items: flex-start;
      width: 100%;
    }
  }
  
  .item-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-primary);
    flex-shrink: 0;
  }
  
  .item-content {
    display: flex;
    gap: var(--spacing-sm);
    min-width: 0;
    
    .item-label {
      font-size: 13px;
      color: var(--color-text-secondary);
      font-weight: 500;
      white-space: nowrap;
      line-height: 1.6;
    }
    
    .item-value {
      font-size: 14px;
      color: var(--color-text-primary);
      font-weight: 600;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 1.6;
      
      &.success {
        color: var(--color-success);
      }
      
      &.error {
        color: var(--color-danger);
      }
      
      &.warning {
        color: var(--color-warning);
      }
      
      &.highlight {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        font-size: 14px;
      }
      
      &.code {
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        font-size: 12px;
      }
      
      &.time {
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        font-size: 12px;
      }
    }
  }
  
  .item-action {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .info-item {
    padding: var(--spacing-xs) var(--spacing-sm);
    min-height: 40px;
    
    &.horizontal {
      flex-direction: column;
      align-items: flex-start;
      gap: var(--spacing-xs);
      
      .item-content {
        flex-direction: column;
        align-items: flex-start;
        width: 100%;
        gap: var(--spacing-xs);
      }
    }
    
    .item-content {
      .item-label {
        font-size: 12px;
      }
      
      .item-value {
        font-size: 13px;
      }
    }
  }
}
</style>