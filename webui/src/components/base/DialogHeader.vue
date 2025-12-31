<template>
  <div class="modern-dialog-header">
    <div class="header-left">
      <div v-if="icon" class="header-icon-wrapper">
        <el-icon class="header-icon">
          <component :is="icon" />
        </el-icon>
      </div>
      <h2 class="header-title">{{ title }}</h2>
    </div>
    <div class="header-right">
      <slot name="actions">
        <el-button 
          v-if="showClose"
          type="danger" 
          :icon="Close" 
          circle 
          size="small"
          @click="handleClose"
          class="close-button"
        />
      </slot>
    </div>
  </div>
</template>

<script setup>
import { Close } from '@element-plus/icons-vue'

const props = defineProps({
  // 标题
  title: {
    type: String,
    required: true,
  },
  // 图标组件
  icon: {
    type: Object,
    default: null,
  },
  // 显示关闭按钮
  showClose: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['close'])

const handleClose = () => {
  emit('close')
}
</script>

<style lang="scss" scoped>
.modern-dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-base) var(--spacing-lg);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: var(--shadow-md);
  min-height: 68px;
  
  .header-left {
    display: flex;
    align-items: center;
    gap: var(--spacing-base);
    flex: 1;
    min-width: 0;
    
    .header-icon-wrapper {
      width: 36px;
      height: 36px;
      border-radius: var(--radius-md);
      background: rgba(255, 255, 255, 0.2);
      display: flex;
      align-items: center;
      justify-content: center;
      backdrop-filter: blur(10px);
      flex-shrink: 0;
      
      .header-icon {
        font-size: 20px;
        color: white;
      }
    }
    
    .header-title {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: white;
      line-height: 1.4;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    flex-shrink: 0;
    
    :deep(.el-button) {
      &:not(.close-button) {
        background: rgba(255, 255, 255, 0.2);
        border: 1px solid rgba(255, 255, 255, 0.3);
        color: white;
        backdrop-filter: blur(10px);
        
        &:hover {
          background: rgba(255, 255, 255, 0.3);
          transform: translateY(-2px);
          box-shadow: var(--shadow-sm);
        }
        
        &.el-button--primary {
          background: white;
          color: #667eea;
          border: none;
          
          &:hover {
            background: rgba(255, 255, 255, 0.9);
            transform: translateY(-2px);
            box-shadow: var(--shadow-md);
          }
        }
      }
    }
    
    .close-button {
      background: rgba(255, 255, 255, 0.2);
      border: none;
      color: white;
      backdrop-filter: blur(10px);
      
      &:hover {
        background: rgba(255, 255, 255, 0.3);
        transform: rotate(90deg);
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .modern-dialog-header {
    padding: var(--spacing-sm) var(--spacing-base);
    flex-wrap: wrap;
    gap: var(--spacing-sm);
    min-height: 64px;
    
    .header-left {
      gap: var(--spacing-sm);
      
      .header-icon-wrapper {
        width: 32px;
        height: 32px;
        
        .header-icon {
          font-size: 18px;
        }
      }
      
      .header-title {
        font-size: 16px;
      }
    }
    
    .header-right {
      :deep(.el-button) {
        font-size: 13px;
        padding: var(--spacing-xs) var(--spacing-sm);
      }
    }
  }
}
</style>