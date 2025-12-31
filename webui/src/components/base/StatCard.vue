<template>
  <div 
    class="stat-card" 
    :class="{ clickable, loading: isLoading }"
    @click="handleClick"
  >
    <div class="stat-icon" :class="variant">
      <el-icon :size="iconSize">
        <component :is="icon" />
      </el-icon>
    </div>
    <div class="stat-content">
      <p class="stat-label">{{ label }}</p>
      <h3 class="stat-value">
        {{ isLoading ? '...' : formattedValue }}
      </h3>
      <p v-if="subLabel" class="stat-sublabel">{{ subLabel }}</p>
    </div>
    <div v-if="trend" class="stat-trend" :class="trendDirection">
      <el-icon>
        <component :is="trendIcon" />
      </el-icon>
      <span>{{ trend }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { formatNumber } from '@/utils/format'
import { ArrowUp, ArrowDown, Minus } from '@element-plus/icons-vue'

const props = defineProps({
  // 图标组件
  icon: {
    type: Object,
    required: true,
  },
  // 标签文本
  label: {
    type: String,
    required: true,
  },
  // 数值
  value: {
    type: [Number, String],
    required: true,
  },
  // 副标签
  subLabel: {
    type: String,
    default: '',
  },
  // 变体样式
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => [
      'primary', 'success', 'warning', 'danger', 'info', 
      'running', 'failed', 'active', 'disabled', 'jinja', 'total'
    ].includes(value),
  },
  // 是否可点击
  clickable: {
    type: Boolean,
    default: false,
  },
  // 加载状态
  isLoading: {
    type: Boolean,
    default: false,
  },
  // 图标大小
  iconSize: {
    type: Number,
    default: 24,
  },
  // 趋势（可选）
  trend: {
    type: String,
    default: '',
  },
  // 趋势方向：up, down, neutral
  trendDirection: {
    type: String,
    default: 'neutral',
    validator: (value) => ['up', 'down', 'neutral'].includes(value),
  },
  // 数值格式化
  formatValue: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['click'])

/**
 * 格式化后的数值
 */
const formattedValue = computed(() => {
  if (typeof props.value === 'string') return props.value
  return props.formatValue ? formatNumber(props.value) : props.value
})

/**
 * 趋势图标
 */
const trendIcon = computed(() => {
  const iconMap = {
    up: ArrowUp,
    down: ArrowDown,
    neutral: Minus,
  }
  return iconMap[props.trendDirection] || Minus
})

/**
 * 点击处理
 */
const handleClick = () => {
  if (props.clickable && !props.isLoading) {
    emit('click')
  }
}
</script>

<style lang="scss" scoped>
.stat-card {
  background: white;
  border-radius: var(--radius-md);
  padding: var(--spacing-base);
  display: flex;
  align-items: center;
  gap: var(--spacing-base);
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
  border: 1px solid var(--color-border-light);
  position: relative;
  overflow: hidden;
  min-height: 88px;
  
  &.clickable {
    cursor: pointer;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: var(--shadow-lg);
      border-color: var(--color-primary);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
  
  &.loading {
    pointer-events: none;
    opacity: 0.6;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: var(--radius-md);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    box-shadow: var(--shadow-sm);
    
    &.primary {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
    }
    
    &.success {
      background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
      color: white;
    }
    
    &.warning {
      background: linear-gradient(135deg, #fbc2eb 0%, #a6c1ee 100%);
      color: white;
    }
    
    &.danger {
      background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
      color: white;
    }
    
    &.info {
      background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
      color: var(--color-text-primary);
    }
    
    &.running {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
    }
    
    &.failed {
      background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
      color: white;
    }
    
    &.active {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
    }
    
    &.disabled {
      background: linear-gradient(135deg, #868f96 0%, #596164 100%);
      color: white;
    }
    
    &.jinja {
      background: linear-gradient(135deg, #fbc2eb 0%, #a6c1ee 100%);
      color: white;
    }
    
    &.total {
      background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
      color: var(--color-text-primary);
    }
  }
  
  .stat-content {
    flex: 1;
    min-width: 0;
    
    .stat-label {
      margin: 0 0 4px 0;
      font-size: 13px;
      color: var(--color-text-secondary);
      font-weight: 500;
      line-height: 1.4;
    }
    
    .stat-value {
      margin: 0;
      font-size: 24px;
      font-weight: 700;
      color: var(--color-text-primary);
      line-height: 1.2;
    }
    
    .stat-sublabel {
      margin: 4px 0 0 0;
      font-size: 12px;
      color: var(--color-text-tertiary);
      line-height: 1.4;
    }
  }
  
  .stat-trend {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    border-radius: var(--radius-base);
    font-size: 12px;
    font-weight: 600;
    
    &.up {
      background: rgba(82, 196, 26, 0.1);
      color: var(--color-success);
    }
    
    &.down {
      background: rgba(255, 77, 79, 0.1);
      color: var(--color-danger);
    }
    
    &.neutral {
      background: var(--color-bg-tertiary);
      color: var(--color-text-secondary);
    }
    
    .el-icon {
      font-size: 11px;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .stat-card {
    padding: var(--spacing-sm);
    min-height: 72px;
    
    .stat-icon {
      width: 40px;
      height: 40px;
    }
    
    .stat-content {
      .stat-label {
        font-size: 12px;
      }
      
      .stat-value {
        font-size: 20px;
      }
      
      .stat-sublabel {
        font-size: 11px;
      }
    }
    
    .stat-trend {
      font-size: 11px;
      padding: 3px 6px;
    }
  }
}
</style>