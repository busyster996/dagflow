<template>
  <el-tag 
    :type="statusType" 
    :size="size" 
    :effect="effect"
    :class="['status-tag', `status-${status}`, { 'with-icon': showIcon }]"
  >
    <el-icon v-if="showIcon" class="tag-icon">
      <component :is="statusIcon" />
    </el-icon>
    <span class="tag-text">{{ statusText }}</span>
  </el-tag>
</template>

<script setup>
import { computed } from 'vue'
import { getStatusType, getStatusText } from '@/utils/status'
import {
  VideoPlay,
  CircleCheckFilled,
  CircleCloseFilled,
  Warning,
  VideoPause,
  Close,
} from '@element-plus/icons-vue'

const props = defineProps({
  // 状态值
  status: {
    type: String,
    required: true,
  },
  // 显示图标
  showIcon: {
    type: Boolean,
    default: false,
  },
  // 标签大小
  size: {
    type: String,
    default: 'default',
    validator: (value) => ['large', 'default', 'small'].includes(value),
  },
  // 标签效果
  effect: {
    type: String,
    default: 'light',
    validator: (value) => ['dark', 'light', 'plain'].includes(value),
  },
  // 自定义文本
  customText: {
    type: String,
    default: null,
  },
})

/**
 * 状态类型（Element Plus）
 */
const statusType = computed(() => getStatusType(props.status))

/**
 * 状态文本
 */
const statusText = computed(() => props.customText || getStatusText(props.status))

/**
 * 状态图标
 */
const statusIcon = computed(() => {
  const iconMap = {
    running: VideoPlay,
    stopped: CircleCheckFilled,
    success: CircleCheckFilled,
    failed: CircleCloseFilled,
    pending: Warning,
    paused: VideoPause,
    timeout: Close,
    killed: Close,
    canceled: Close,
  }
  return iconMap[props.status] || Warning
})
</script>

<style lang="scss" scoped>
.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
  border-radius: var(--radius-base);
  transition: all var(--transition-fast);
  
  &.with-icon {
    padding-left: 6px;
  }
  
  .tag-icon {
    font-size: 14px;
  }
  
  .tag-text {
    line-height: 1;
  }
  
  // 运行中状态添加脉冲动画
  &.status-running {
    animation: tagPulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }
}

@keyframes tagPulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}
</style>