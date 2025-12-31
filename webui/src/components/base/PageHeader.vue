<template>
  <div class="page-header">
    <div class="header-left">
      <div v-if="icon" class="header-icon-wrapper">
        <el-icon class="header-icon" :size="iconSize">
          <component :is="icon" />
        </el-icon>
      </div>
      <div class="header-content">
        <h2 class="header-title">{{ title }}</h2>
        <p v-if="subtitle" class="header-subtitle">{{ subtitle }}</p>
      </div>
    </div>
    <div class="header-right">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  // 标题
  title: {
    type: String,
    required: true,
  },
  // 副标题
  subtitle: {
    type: String,
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
    default: 24,
  },
})
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-base) var(--spacing-lg);
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border-light);
  transition: all var(--transition-base);
  
  &:hover {
    box-shadow: var(--shadow-md);
  }
  
  .header-left {
    display: flex;
    align-items: center;
    gap: var(--spacing-base);
    flex: 1;
    min-width: 0;
    
    .header-icon-wrapper {
      width: 40px;
      height: 40px;
      border-radius: var(--radius-md);
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: var(--shadow-sm);
      flex-shrink: 0;
      
      .header-icon {
        color: white;
      }
    }
    
    .header-content {
      flex: 1;
      min-width: 0;
      
      .header-title {
        margin: 0;
        font-size: 18px;
        font-weight: 700;
        color: var(--color-text-primary);
        line-height: 1.4;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      
      .header-subtitle {
        margin: 4px 0 0 0;
        font-size: 13px;
        color: var(--color-text-secondary);
        font-weight: 400;
        line-height: 1.5;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    flex-shrink: 0;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-base);
    padding: var(--spacing-base);
    
    .header-left {
      gap: var(--spacing-sm);
      
      .header-icon-wrapper {
        width: 36px;
        height: 36px;
      }
      
      .header-content {
        .header-title {
          font-size: 16px;
        }
        
        .header-subtitle {
          font-size: 12px;
        }
      }
    }
    
    .header-right {
      width: 100%;
      justify-content: flex-end;
    }
  }
}
</style>