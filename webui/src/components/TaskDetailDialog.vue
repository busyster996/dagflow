<template>
  <el-dialog
    v-model="dialogVisible"
    :close-on-click-modal="false"
    :show-close="false"
    @close="handleClose"
    fullscreen
    class="task-detail-dialog"
  >
    <template #header>
      <DialogHeader
        :icon="Menu"
        :title="taskName"
        @close="dialogVisible = false"
      />
    </template>

    <div v-loading="loading" class="dialog-content">
      <!-- 侧边信息面板 -->
      <div class="sidebar-panel">
        <!-- 任务信息卡片 -->
        <el-card shadow="hover" class="info-card">
          <template #header>
            <SectionHeader
              :icon="InfoFilled"
              title="任务信息"
            />
          </template>
          <div class="info-content" v-if="taskData">
            <InfoItem
              label="任务名称"
              :value="taskData.name"
              hoverable
            />
            <InfoItem
              label="任务类型"
              :value="taskData.kind || 'dag'"
              hoverable
            />
            <InfoItem label="任务状态" hoverable>
              <template #value>
                <StatusTag :status="taskData.state" effect="light" size="large" show-icon />
              </template>
            </InfoItem>
            <InfoItem
              label="步骤数量"
              :value="taskData.count"
              value-class="highlight"
              hoverable
            />
          </div>
        </el-card>

        <!-- 执行信息卡片 -->
        <el-card shadow="hover" class="message-card">
          <template #header>
            <SectionHeader
              :icon="ChatLineSquare"
              title="执行信息"
            />
          </template>
          <div class="message-content">
            <p v-if="taskData?.message" class="message-text">{{ taskData.message }}</p>
            <EmptyState v-else description="暂无执行信息" size="small" />
          </div>
        </el-card>

        <!-- 环境变量卡片 -->
        <el-card shadow="hover" v-if="taskData?.env && taskData.env.length > 0" class="env-card">
          <template #header>
            <SectionHeader
              :icon="Key"
              title="环境变量"
              :count="taskData.env.length"
            />
          </template>
          <el-scrollbar max-height="280px">
            <div class="env-list">
              <div v-for="(env, index) in taskData.env" :key="index" class="env-item">
                <div class="env-name">
                  <el-icon class="env-icon"><Setting /></el-icon>
                  <span>{{ env.name }}</span>
                </div>
                <div class="env-value">{{ env.value }}</div>
              </div>
            </div>
          </el-scrollbar>
        </el-card>
      </div>

      <!-- 主要内容区 - DAG图形 -->
      <div class="main-content">
        <el-card shadow="hover" class="graph-card">
          <template #header>
            <div class="graph-header">
              <div class="graph-header-left">
                <SectionHeader
                  :icon="Connection"
                  title="DAG 流程图"
                />
              </div>
              <div class="graph-controls">
                <el-button-group size="small">
                  <el-button @click="handleFitView">
                    <el-icon><FullScreen /></el-icon>
                    适应画布
                  </el-button>
                  <el-button @click="handleZoomIn">
                    <el-icon><ZoomIn /></el-icon>
                    放大
                  </el-button>
                  <el-button @click="handleZoomOut">
                    <el-icon><ZoomOut /></el-icon>
                    缩小
                  </el-button>
                  <el-button @click="handleActualSize">
                    <el-icon><Refresh /></el-icon>
                    重置
                  </el-button>
                </el-button-group>
              </div>
            </div>
          </template>
          
          <div class="graph-wrapper">
            <VueFlowGraph
              v-if="steps.length > 0"
              ref="graphRef"
              :steps="steps"
              :task-name="taskName"
              :show-toolbar="false"
              @node-click="handleNodeClick"
            />
            <EmptyState v-else :icon="Connection" description="暂无步骤数据" />
          </div>
        </el-card>
      </div>
    </div>

    <!-- 步骤详情抽屉 -->
    <el-drawer
      v-for="step in openedSteps"
      :key="step.name"
      v-model="step.visible"
      direction="rtl"
      size="600px"
      @close="closeStep(step.name)"
      class="step-drawer"
    >
      <template #header>
        <div class="drawer-header">
          <div class="drawer-header-left">
            <div class="drawer-icon-wrapper">
              <el-icon><Operation /></el-icon>
            </div>
            <h3 class="drawer-title">{{ step.name }}</h3>
          </div>
          <el-button
            type="danger"
            :icon="Close"
            circle
            size="small"
            @click="closeStep(step.name)"
          />
        </div>
      </template>
      <StepContent :task-name="taskName" :step-name="step.name" />
    </el-drawer>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onUnmounted, computed, shallowRef } from 'vue'
import {
  Menu,
  Close,
  InfoFilled,
  ChatLineSquare,
  Key,
  Setting,
  Connection,
  FullScreen,
  ZoomIn,
  ZoomOut,
  Refresh,
  Operation,
} from '@element-plus/icons-vue'
import { getTaskDetail } from '@/api/task'
import { useWebSocketList } from '@/composables/useWebSocket'
import { API_ENDPOINTS } from '@/config'
import { DialogHeader, SectionHeader, InfoItem, StatusTag, EmptyState } from '@/components/base'
import VueFlowGraph from './VueFlowGraph.vue'
import StepContent from './StepContent.vue'
import { throttle, debounce } from '@/utils/throttle'

const props = defineProps({
  modelValue: Boolean,
  taskName: String,
})

const emit = defineEmits(['update:modelValue'])

const dialogVisible = ref(false)
const loading = ref(false)
const taskData = ref(null)
const steps = shallowRef([]) // 使用shallowRef避免深度响应
const openedSteps = ref([])
const graphRef = ref(null)
let wsManager = null

// 性能监控
const performanceStats = ref({
  wsMessageCount: 0,
  httpRequestCount: 0,
  lastUpdateTime: 0,
  renderCount: 0,
})

// 记录上次任务状态，避免无意义的HTTP请求
let lastTaskState = null
let taskDetailCache = null

watch(
  () => props.modelValue,
  async (val) => {
    dialogVisible.value = val
    if (val && props.taskName) {
      await loadTaskDetail()
      initWebSocket()
    }
  }
)

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
  if (!val) {
    handleClose()
  }
})

/**
 * 节流更新步骤数据 - 避免高频渲染
 */
const throttledUpdateSteps = throttle((newSteps) => {
  performanceStats.value.renderCount++
  
  // 增量更新：只更新变化的步骤
  if (steps.value.length === newSteps.length) {
    let hasChanges = false
    const updatedSteps = steps.value.map((oldStep, index) => {
      const newStep = newSteps[index]
      // 比较关键属性是否变化
      if (
        oldStep.name === newStep.name &&
        oldStep.state === newStep.state &&
        oldStep.code === newStep.code &&
        JSON.stringify(oldStep.time) === JSON.stringify(newStep.time)
      ) {
        return oldStep // 无变化，保持原对象引用
      }
      hasChanges = true
      return newStep // 有变化，使用新对象
    })
    
    if (hasChanges) {
      steps.value = updatedSteps
    }
  } else {
    // 步骤数量变化，全量更新
    steps.value = newSteps
  }
}, 100) // 100ms节流

/**
 * 防抖加载任务详情 - 避免频繁HTTP请求
 */
const debouncedLoadTaskDetail = debounce(async (forceUpdate = false) => {
  performanceStats.value.httpRequestCount++
  
  try {
    const response = await getTaskDetail(props.taskName)
    const newTaskData = response.data
    
    // 检查任务状态是否真的变化
    const currentState = newTaskData?.state
    if (!forceUpdate && lastTaskState === currentState && taskDetailCache) {
      // 状态未变化，使用缓存
      return
    }
    
    lastTaskState = currentState
    taskDetailCache = newTaskData
    taskData.value = newTaskData
  } catch (error) {
    console.error('获取任务详情失败:', error)
  }
}, 500, { maxWait: 2000 }) // 500ms防抖，最多2秒必须执行一次

/**
 * 智能更新任务信息 - 仅在关键状态变化时更新
 */
const smartUpdateTaskInfo = (stepsData) => {
  // 检测任务状态变化
  const hasRunning = stepsData.some(s => s.state === 'running')
  const hasFailed = stepsData.some(s => s.state === 'failed')
  const allStopped = stepsData.length > 0 && stepsData.every(s => s.state === 'stopped' || s.state === 'failed')
  
  const currentTaskState = taskData.value?.state
  
  // 只在以下情况触发HTTP请求：
  // 1. 没有缓存数据
  // 2. 有步骤正在运行且任务状态不是running
  // 3. 所有步骤完成但任务状态不是stopped
  // 4. 有失败步骤但任务状态不是failed
  if (!taskData.value ||
      (hasRunning && currentTaskState !== 'running') ||
      (allStopped && currentTaskState !== 'stopped') ||
      (hasFailed && currentTaskState !== 'failed')) {
    debouncedLoadTaskDetail()
  }
}

/**
 * 初始化WebSocket连接
 */
const initWebSocket = () => {
  if (wsManager) {
    wsManager.close()
    wsManager = null
  }

  const ws = useWebSocketList(
    `${API_ENDPOINTS.task}/${props.taskName}/step`,
    {
      onMessage: (response) => {
        performanceStats.value.wsMessageCount++
        performanceStats.value.lastUpdateTime = Date.now()
        
        if (response.data && Array.isArray(response.data)) {
          // 节流更新步骤数据
          throttledUpdateSteps(response.data)
          
          // 智能更新任务信息（带防抖）
          smartUpdateTaskInfo(response.data)
        }
      }
    }
  )
  
  wsManager = ws
  ws.connect()
}

/**
 * 加载任务详情（初始加载）
 */
const loadTaskDetail = async () => {
  loading.value = true
  try {
    const response = await getTaskDetail(props.taskName)
    taskData.value = response.data
    lastTaskState = response.data?.state
    taskDetailCache = response.data
    performanceStats.value.httpRequestCount++
  } catch (error) {
    console.error('获取任务详情失败:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 处理节点点击
 */
const handleNodeClick = (stepName) => {
  const existing = openedSteps.value.find(s => s.name === stepName)
  if (!existing) {
    openedSteps.value.push({ name: stepName, visible: true })
  } else {
    existing.visible = true
  }
}

/**
 * 关闭步骤抽屉
 */
const closeStep = (stepName) => {
  const index = openedSteps.value.findIndex(s => s.name === stepName)
  if (index > -1) {
    openedSteps.value.splice(index, 1)
  }
}

/**
 * 图形控制方法
 */
const handleFitView = () => {
  if (graphRef.value) {
    graphRef.value.fitView()
  }
}

const handleZoomIn = () => {
  if (graphRef.value) {
    graphRef.value.zoomIn()
  }
}

const handleZoomOut = () => {
  if (graphRef.value) {
    graphRef.value.zoomOut()
  }
}

const handleActualSize = () => {
  if (graphRef.value) {
    graphRef.value.actualSize()
  }
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  // 取消所有待执行的防抖/节流任务
  if (debouncedLoadTaskDetail.cancel) {
    debouncedLoadTaskDetail.cancel()
  }
  if (throttledUpdateSteps.cancel) {
    throttledUpdateSteps.cancel()
  }
  
  if (wsManager) {
    wsManager.close()
    wsManager = null
  }
  
  openedSteps.value = []
  steps.value = []
  taskData.value = null
  lastTaskState = null
  taskDetailCache = null
  
  // 打印性能统计（开发模式）
  if (import.meta.env.DEV) {
    console.log('📊 TaskDetail性能统计:', performanceStats.value)
  }
  
  // 重置性能统计
  performanceStats.value = {
    wsMessageCount: 0,
    httpRequestCount: 0,
    lastUpdateTime: 0,
    renderCount: 0,
  }
}

onUnmounted(() => {
  handleClose()
})

// 开发模式下监控性能
if (import.meta.env.DEV) {
  watch(performanceStats, (stats) => {
    if (stats.wsMessageCount > 100 || stats.httpRequestCount > 50) {
      console.warn('⚠️ 性能警告: WebSocket消息或HTTP请求过多', stats)
    }
  }, { deep: true })
}
</script>

<style lang="scss" scoped>
.task-detail-dialog {
  :deep(.el-dialog) {
    background: var(--color-bg-secondary);
  }

  :deep(.el-dialog__header) {
    padding: 0 !important;
    margin: 0 !important;
    border: none !important;
  }

  :deep(.el-dialog__body) {
    padding: 0 !important;
    height: calc(100vh - 68px) !important;
    overflow: hidden !important;
  }
}

// ==========================================
// 对话框内容
// ==========================================
.dialog-content {
  display: grid;
  grid-template-columns: 340px 1fr;
  gap: var(--spacing-base);
  height: 100%;
  padding: var(--spacing-lg);
  overflow: hidden;
}

// 侧边信息面板
.sidebar-panel {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-base);
  overflow-y: auto;
  overflow-x: hidden;

  .info-card {
    .info-content {
      display: flex;
      flex-direction: column;
      gap: var(--spacing-xs);
    }
  }

  .message-card {
    .message-content {
      min-height: 80px;

      .message-text {
        margin: 0;
        padding: var(--spacing-base);
        background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
        border-radius: var(--radius-md);
        border-left: 3px solid var(--color-primary);
        color: var(--color-text-primary);
        line-height: 1.6;
        word-break: break-word;
        font-size: 14px;
      }
    }
  }

  .env-card {
    .env-list {
      .env-item {
        padding: var(--spacing-sm) var(--spacing-base);
        margin-bottom: var(--spacing-sm);
        background: var(--color-bg-secondary);
        border-radius: var(--radius-md);
        border-left: 3px solid var(--color-primary);
        transition: all var(--transition-base);

        &:hover {
          background: var(--color-bg-tertiary);
          transform: translateX(2px);
          box-shadow: var(--shadow-sm);
        }

        &:last-child {
          margin-bottom: 0;
        }

        .env-name {
          display: flex;
          align-items: center;
          gap: 6px;
          margin-bottom: 4px;
          font-weight: 600;
          color: var(--color-text-primary);
          font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
          font-size: 13px;
          line-height: 1.5;

          .env-icon {
            font-size: 14px;
            color: var(--color-primary);
          }
        }

        .env-value {
          color: var(--color-text-secondary);
          font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
          font-size: 12px;
          word-break: break-all;
          padding-left: 20px;
          line-height: 1.6;
        }
      }
    }
  }
}

// 主内容区
.main-content {
  overflow: hidden;
  display: flex;
  flex-direction: column;

  .graph-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    :deep(.el-card__header) {
      padding: var(--spacing-base) var(--spacing-lg);
      background: var(--color-bg-secondary);
      border-bottom: 1px solid var(--color-border-light);
    }

    :deep(.el-card__body) {
      flex: 1;
      padding: 0;
      overflow: hidden;
      min-height: 0;
      display: flex;
      flex-direction: column;
    }

    .graph-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
    }

    .graph-wrapper {
      flex: 1;
      width: 100%;
      min-height: 0;
      overflow: hidden;
      background: linear-gradient(135deg, #f6f8fb 0%, #ffffff 100%);
    }
  }
}

// ==========================================
// 步骤抽屉样式
// ==========================================
.step-drawer {
  :deep(.el-drawer__header) {
    padding: 0 !important;
    margin: 0 !important;
    border: none !important;
  }

  :deep(.el-drawer__body) {
    padding: 0;
  }

  .drawer-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--spacing-sm) var(--spacing-md);
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    box-shadow: none;

    .drawer-header-left {
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);

      .drawer-icon-wrapper {
        width: 24px;
        height: 24px;
        border-radius: var(--radius-sm);
        background: rgba(255, 255, 255, 0.2);
        display: flex;
        align-items: center;
        justify-content: center;
        backdrop-filter: blur(10px);

        .el-icon {
          font-size: 14px;
          color: white;
        }
      }

      .drawer-title {
        margin: 0;
        font-size: 13px;
        font-weight: 600;
        color: white;
        line-height: 1;
      }
    }

    .el-button {
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

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 1200px) {
  .dialog-content {
    grid-template-columns: 320px 1fr;
  }

  .step-drawer {
    :deep(.el-drawer) {
      width: 500px !important;
    }
  }
}

@media (max-width: 768px) {
  .dialog-content {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
    padding: var(--spacing-md);
  }

  .sidebar-panel {
    max-height: 40vh;
  }

  .step-drawer {
    :deep(.el-drawer) {
      width: 90% !important;
    }
  }
}
</style>