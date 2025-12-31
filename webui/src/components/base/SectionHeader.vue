<template>
  <div class="section-header" :class="{ clickable, bordered }">
    <div class="header-main">
      <el-icon v-if="icon" class="section-icon" :size="iconSize">
        <component :is="icon" />
      </el-icon>
      <span class="section-title">{{ title }}</span>
      <el-badge v-if="count !== null && count !== undefined" :value="count" :max="max" class="count-badge" />
      <el-tag v-if="tag" :type="tagType" size="small" class="section-tag">{{ tag }}</el-tag>
    </div>
    <div v-if="$slots.actions || expandable" class="header-actions">
      <slot name="actions" />
      <el-icon v-if="expandable" class="expand-icon" :class="{ expanded: isExpanded }">
        <ArrowDown />
      </el-icon>
    </div>
  </div>
</template>

<script setup>
import { ArrowDown } from '@element-plus/icons-vue'

const props = defineProps({
  // 标题文本
  title: {
    type: String,
    required: true,
  },
  // 图标组件
  icon: {
    type: Object,
    default: null,
  },
  // 图标大小
  iconSize: {
    type: Number,
    default: 18,
  },
  // 计数徽章
  count: {
    type: Number,
    default: null,
  },
  // 徽章最大值
  max: {
    type: Number,
    default: 99,
  },
  // 标签文本
  tag: {
    type: String,
    default: '',
  },
  // 标签类型
  tagType: {
    type: String,
    default: 'info',
  },
  // 是否可点击
  clickable: {
    type: Boolean,
    default: false,
  },
  // 是否可展开/收起
  expandable: {
    type: Boolean,
    default: false,
  },
  // 展开状态
  isExpanded: {
    type: Boolean,
    default: false,
  },
  // 是否显示底部边框
  bordered: {
    type: Boolean,
    default: false,
  },
})
</script>

<style lang="scss" scoped>
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-base);
  padding: 0;
  transition: all var(--transition-base);
  
  &.clickable {
    cursor: pointer;
    user-select: none;
    
    &:hover {
      .section-title {
        color: var(--color-primary);
      }
    }
  }
  
  &.bordered {
    border-bottom: 1px solid var(--color-border-light);
    padding-bottom: var(--spacing-sm);
  }
  
  .header-main {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    flex: 1;
    min-width: 0;
    
    .section-icon {
      color: var(--color-primary);
      flex-shrink: 0;
    }
    
    .section-title {
      font-size: 15px;
      font-weight: 600;
      color: var(--color-text-primary);
      transition: color var(--transition-base);
      white-space: nowrap;
      line-height: 1.5;
    }
    
    .count-badge {
      flex-shrink: 0;
    }
    
    .section-tag {
      margin-left: var(--spacing-sm);
      flex-shrink: 0;
    }
  }
  
  .header-actions {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    flex-shrink: 0;
    
    .expand-icon {
      font-size: 18px;
      transition: transform var(--transition-base);
      color: var(--color-text-tertiary);
      
      &.expanded {
        transform: rotate(180deg);
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .section-header {
    gap: var(--spacing-sm);
    
    .header-main {
      gap: var(--spacing-xs);
      
      .section-title {
        font-size: 14px;
      }
    }
    
    .header-actions {
      .expand-icon {
        font-size: 16px;
      }
    }
  }
}
</style>