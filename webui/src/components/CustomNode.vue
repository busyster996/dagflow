<template>
  <div class="custom-node" :class="`node-${data.step.state}`">
    <!-- 状态指示条 -->
    <div class="status-bar" :style="{ background: statusGradient }">
      <div class="status-pulse" v-if="data.step.state === 'running'"></div>
    </div>

    <!-- 节点头部 -->
    <div class="node-header">
      <div class="header-icon-wrapper">
        <svg class="step-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z"/>
        </svg>
      </div>
      <div class="header-content">
        <h3 class="step-name" :title="data.step.name">{{ data.step.name }}</h3>
        <div class="step-type">
          <el-icon><Setting /></el-icon>
          <span>{{ data.step.type || 'exec' }}</span>
        </div>
      </div>
    </div>

    <!-- 节点主体 -->
    <div class="node-body">
      <!-- 状态与代码 -->
      <div class="info-section">
        <div class="info-row">
          <div class="info-item" v-if="data.step.code !== undefined">
            <span class="item-label">返回码</span>
            <span class="item-value code" :class="getCodeClass()">
              {{ data.step.code }}
            </span>
          </div>
          <div class="info-item">
            <span class="item-label">状态</span>
            <div class="status-badge" :style="{ background: nodeColor }">
              <span class="status-dot"></span>
              <span class="status-text">{{ getStateLabel(data.step.state) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 分隔线 -->
      <div class="node-divider"></div>

      <!-- 时间信息 -->
      <div class="time-section">
        <div class="time-row">
          <div class="time-item">
            <el-icon class="time-icon"><Timer /></el-icon>
            <div class="time-content">
              <span class="time-label">开始</span>
              <span class="time-value">{{ formatTime(data.step.time?.start) }}</span>
            </div>
          </div>
        </div>
        <div class="time-row">
          <div class="time-item">
            <el-icon class="time-icon"><Calendar /></el-icon>
            <div class="time-content">
              <span class="time-label">结束</span>
              <span class="time-value">{{ formatTime(data.step.time?.end) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 悬浮操作按钮 -->
    <div class="node-actions">
      <el-icon class="action-icon"><More /></el-icon>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { STATUS_COLORS } from '@/config'

const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
})

const nodeColor = computed(() => {
  return STATUS_COLORS[props.data.step.state] || STATUS_COLORS.unknown
})

const statusGradient = computed(() => {
  const color = nodeColor.value
  return `linear-gradient(90deg, ${color} 0%, ${adjustColor(color, 20)} 100%)`
})

// 调整颜色亮度
const adjustColor = (color, amount) => {
  const clamp = (val) => Math.min(Math.max(val, 0), 255)
  const num = parseInt(color.replace('#', ''), 16)
  const r = clamp((num >> 16) + amount)
  const g = clamp(((num >> 8) & 0x00FF) + amount)
  const b = clamp((num & 0x0000FF) + amount)
  return '#' + ((r << 16) | (g << 8) | b).toString(16).padStart(6, '0')
}

// 获取状态标签
const getStateLabel = (state) => {
  const labels = {
    'pending': '待执行',
    'running': '执行中',
    'stopped': '已停止',
    'failed': '已失败',
    'success': '成功',
    'paused': '已暂停',
    'killed': '已终止',
    'timeout': '超时',
    'canceled': '已取消',
    'skipped': '已跳过',
  }
  return labels[state] || state
}

// 获取状态码样式
const getCodeClass = () => {
  const code = props.data.step.code
  if (code === 0) return 'success'
  if (code > 0) return 'error'
  return 'default'
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '---'
  if (typeof time === 'string' && time.includes(':')) {
    // 提取时分秒
    const match = time.match(/(\d{2}):(\d{2}):(\d{2})/)
    if (match) {
      return `${match[1]}:${match[2]}:${match[3]}`
    }
    return time
  }
  return time
}
</script>

<style lang="scss" scoped>
.custom-node {
  width: 300px;
  min-height: 160px;
  background: white;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  transition: all var(--transition-base);
  border: 1px solid var(--color-border-light);
  position: relative;

  &:hover {
    box-shadow: var(--shadow-xl);
    transform: translateY(-3px) scale(1.02);
    border-color: var(--color-primary);

    .node-actions {
      opacity: 1;
    }
  }

  // 不同状态的节点样式
  &.node-running {
    animation: nodeGlow 2s ease-in-out infinite;
  }

  &.node-failed {
    .node-header {
      background: linear-gradient(135deg, #fee2e2 0%, #fef2f2 100%);
    }
  }

  &.node-stopped {
    .node-header {
      background: linear-gradient(135deg, #d1fae5 0%, #ecfdf5 100%);
    }
  }
}

// 状态指示条
.status-bar {
  height: 4px;
  width: 100%;
  position: relative;
  overflow: hidden;

  .status-pulse {
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.6), transparent);
    animation: statusPulse 2s ease-in-out infinite;
  }
}

@keyframes statusPulse {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

@keyframes nodeGlow {
  0%, 100% {
    box-shadow: var(--shadow-sm), 0 0 12px rgba(102, 126, 234, 0.2);
  }
  50% {
    box-shadow: var(--shadow-md), 0 0 16px rgba(102, 126, 234, 0.3);
  }
}

// 节点头部
.node-header {
  padding: var(--spacing-sm) var(--spacing-base);
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border-bottom: 1px solid var(--color-border-light);
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  min-height: 56px;

  .header-icon-wrapper {
    width: 32px;
    height: 32px;
    border-radius: var(--radius-md);
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-sm);
    flex-shrink: 0;

    .step-icon {
      width: 18px;
      height: 18px;
      color: white;
    }
  }

  .header-content {
    flex: 1;
    min-width: 0;

    .step-name {
      margin: 0 0 4px 0;
      font-size: 14px;
      font-weight: 600;
      color: var(--color-text-primary);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      line-height: 1.4;
    }

    .step-type {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;
      color: var(--color-text-secondary);
      font-weight: 500;
      line-height: 1.4;

      .el-icon {
        font-size: 13px;
      }
    }
  }
}

// 节点主体
.node-body {
  padding: var(--spacing-sm) var(--spacing-base);
  background: white;
}

// 信息区域
.info-section {
  margin-bottom: var(--spacing-base);

  .info-row {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-sm);
  }

  .info-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: var(--spacing-sm);
    background: var(--color-bg-tertiary);
    border-radius: var(--radius-md);

    .item-label {
      font-size: 11px;
      color: var(--color-text-tertiary);
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      line-height: 1.4;
    }

    .item-value {
      font-size: 14px;
      font-weight: 600;
      color: var(--color-text-primary);
      line-height: 1.5;

      &.code {
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        padding: 2px 8px;
        border-radius: var(--radius-base);
        display: inline-block;
        font-size: 13px;

        &.success {
          background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
          color: #059669;
        }

        &.error {
          background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%);
          color: #dc2626;
        }

        &.default {
          background: var(--color-bg-secondary);
          color: var(--color-text-secondary);
        }
      }
    }

    .status-badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 4px 10px;
      border-radius: var(--radius-base);
      font-size: 12px;
      font-weight: 600;
      color: white;
      box-shadow: var(--shadow-sm);

      .status-dot {
        width: 6px;
        height: 6px;
        border-radius: var(--radius-full);
        background: white;
        animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
      }

      .status-text {
        line-height: 1;
      }
    }
  }
}

// 分隔线
.node-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--color-border-base) 20%, var(--color-border-base) 80%, transparent);
  margin: var(--spacing-base) 0;
}

// 时间区域
.time-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);

  .time-row {
    display: flex;
  }

  .time-item {
    flex: 1;
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    padding: var(--spacing-sm);
    background: var(--color-bg-secondary);
    border-radius: var(--radius-base);
    transition: all var(--transition-fast);

    &:hover {
      background: var(--color-bg-tertiary);
      transform: translateX(2px);
    }

    .time-icon {
      font-size: 16px;
      color: var(--color-primary);
      flex-shrink: 0;
    }

    .time-content {
      flex: 1;
      min-width: 0;
      display: flex;
      flex-direction: column;
      gap: 2px;

      .time-label {
        font-size: 10px;
        color: var(--color-text-tertiary);
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
      }

      .time-value {
        font-size: 11px;
        color: var(--color-text-primary);
        font-weight: 500;
        font-family: 'Consolas', monospace;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

// 操作按钮
.node-actions {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
  width: 28px;
  height: 28px;
  border-radius: var(--radius-base);
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0;
  transition: all var(--transition-base);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(10px);

  &:hover {
    background: var(--color-primary);
    box-shadow: var(--shadow-md);

    .action-icon {
      color: white;
      transform: scale(1.2);
    }
  }

  .action-icon {
    font-size: 16px;
    color: var(--color-text-secondary);
    transition: all var(--transition-fast);
  }
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 768px) {
  .custom-node {
    width: 260px;
  }
}
</style>