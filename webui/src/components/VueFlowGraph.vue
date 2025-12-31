<template>
  <div class="vue-flow-wrapper">
    <VueFlow
      v-model:nodes="nodes"
      v-model:edges="edges"
      :default-zoom="1"
      :min-zoom="0.2"
      :max-zoom="4"
      fit-view-on-init
      @node-click="handleNodeClick"
      @node-context-menu="handleNodeContextMenu"
      class="modern-vue-flow"
    >
      <Background 
        :pattern-color="patternColor" 
        :gap="20"
        variant="dots"
      />
      <Controls v-if="showToolbar" position="bottom-right" />
      <MiniMap v-if="showToolbar" position="bottom-left" />
      
      <template #node-custom="{ data }">
        <CustomNode :data="data" />
      </template>
    </VueFlow>

    <!-- 现代化右键菜单 -->
    <transition name="context-menu-fade">
      <div
        v-if="contextMenu.visible"
        class="modern-context-menu"
        :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
      >
        <div class="context-menu-header">
          <el-icon class="menu-header-icon"><Operation /></el-icon>
          <span>步骤操作</span>
        </div>
        <div class="context-menu-body">
          <div
            v-if="contextMenu.step?.state === 'running'"
            class="menu-item danger"
            @click="handleAction('kill')"
          >
            <el-icon><VideoPause /></el-icon>
            <span>强制终止</span>
          </div>
          <template v-else-if="contextMenu.step?.state === 'paused'">
            <div class="menu-item success" @click="handleAction('resume')">
              <el-icon><VideoPlay /></el-icon>
              <span>恢复执行</span>
            </div>
            <div class="menu-item danger" @click="handleAction('kill')">
              <el-icon><VideoPause /></el-icon>
              <span>强制终止</span>
            </div>
          </template>
          <template v-else-if="contextMenu.step?.state === 'pending'">
            <div class="menu-item warning" @click="handleAction('pause')">
              <el-icon><CirclePause /></el-icon>
              <span>暂停挂起</span>
            </div>
            <div class="menu-item danger" @click="handleAction('kill')">
              <el-icon><VideoPause /></el-icon>
              <span>强制终止</span>
            </div>
          </template>
          <div v-else class="menu-item disabled">
            <el-icon><InfoFilled /></el-icon>
            <span>暂无可用操作</span>
          </div>
        </div>
      </div>
    </transition>

    <!-- 图形统计信息浮层 -->
    <div class="graph-stats" v-if="stats">
      <div class="stat-item">
        <span class="stat-label">节点</span>
        <span class="stat-value">{{ stats.nodes }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">连接</span>
        <span class="stat-value">{{ stats.edges }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted, nextTick, shallowRef } from 'vue'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { MarkerType } from '@vue-flow/core'
import CustomNode from './CustomNode.vue'
import { stepAction } from '@/api/task'
import { rafThrottle } from '@/utils/throttle'

const props = defineProps({
  steps: {
    type: Array,
    required: true,
  },
  taskName: {
    type: String,
    required: true,
  },
  showToolbar: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['node-click'])

const nodes = shallowRef([]) // 使用shallowRef避免深度响应
const edges = shallowRef([]) // 使用shallowRef避免深度响应
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  step: null,
  nodeId: null,
})

const patternColor = computed(() => '#c7d2fe')

const stats = computed(() => {
  if (!nodes.value.length) return null
  return {
    nodes: nodes.value.length,
    edges: edges.value.length,
  }
})

let isFirstRender = true
let savedPositions = {}
let lastStepsHash = ''

// 性能监控
const renderStats = {
  updateCount: 0,
  fullRenderCount: 0,
  incrementalUpdateCount: 0,
}

const { fitView, zoomIn: flowZoomIn, zoomOut: flowZoomOut, setViewport } = useVueFlow()

/**
 * 生成步骤数据哈希用于对比
 */
const generateStepsHash = (steps) => {
  return steps.map(s => `${s.name}:${s.state}:${s.code}`).join('|')
}

/**
 * 增量更新节点 - 只更新变化的节点
 */
const incrementalUpdateNodes = (newSteps) => {
  renderStats.incrementalUpdateCount++
  
  const updatedNodes = nodes.value.map(node => {
    const newStep = newSteps.find(s => s.name === node.id)
    if (!newStep) return node
    
    // 检查步骤状态是否变化
    if (node.data.step.state === newStep.state &&
        node.data.step.code === newStep.code) {
      return node // 无变化，保持原对象
    }
    
    // 状态变化，创建新节点对象
    return {
      ...node,
      data: {
        ...node.data,
        step: newStep,
      },
    }
  })
  
  // 更新边的动画状态
  const updatedEdges = edges.value.map(edge => {
    const targetStep = newSteps.find(s => s.name === edge.target)
    if (!targetStep) return edge
    
    const shouldAnimate = targetStep.state === 'running'
    if (edge.animated === shouldAnimate &&
        edge.style.stroke === getEdgeColor(targetStep.state)) {
      return edge // 无变化
    }
    
    return {
      ...edge,
      animated: shouldAnimate,
      style: {
        ...edge.style,
        stroke: getEdgeColor(targetStep.state),
      },
    }
  })
  
  nodes.value = updatedNodes
  edges.value = updatedEdges
}

/**
 * 转换步骤数据为 Vue Flow 格式
 */
const convertStepsToFlow = () => {
  renderStats.fullRenderCount++
  
  const newNodes = []
  const newEdges = []

  props.steps.forEach((step, index) => {
    const nodeData = {
      id: step.name,
      type: 'custom',
      data: {
        step: step,
        taskName: props.taskName,
      },
      position: savedPositions[step.name] || {
        x: 0,
        y: index * 120,
      },
    }

    newNodes.push(nodeData)

    if (step.depends && step.depends.length > 0) {
      step.depends.forEach((depend) => {
        newEdges.push({
          id: `${depend}-${step.name}`,
          source: depend,
          target: step.name,
          type: 'smoothstep',
          animated: step.state === 'running',
          markerEnd: MarkerType.ArrowClosed,
          style: {
            strokeWidth: 2,
            stroke: getEdgeColor(step.state),
          },
        })
      })
    }
  })

  return { newNodes, newEdges }
}

const getEdgeColor = (state) => {
  const colorMap = {
    running: '#667eea',
    stopped: '#52c41a',
    failed: '#ff4d4f',
    pending: '#faad14',
  }
  return colorMap[state] || '#94a3b8'
}

// 更新图形
const updateGraph = async () => {
  const { newNodes, newEdges } = convertStepsToFlow()
  
  if (isFirstRender) {
    nodes.value = newNodes
    edges.value = newEdges
    
    await nextTick()
    
    applyDagreLayout()
    
    setTimeout(() => {
      nodes.value.forEach(node => {
        savedPositions[node.id] = { ...node.position }
      })
      isFirstRender = false
    }, 500)
  } else {
    nodes.value = newNodes
    edges.value = newEdges
  }
}

// 应用 Dagre 布局算法
const applyDagreLayout = () => {
  const levels = new Map()
  const visited = new Set()
  
  const roots = props.steps.filter(step => !step.depends || step.depends.length === 0)
  const queue = roots.map(step => ({ name: step.name, level: 0 }))
  
  while (queue.length > 0) {
    const { name, level } = queue.shift()
    
    if (visited.has(name)) continue
    visited.add(name)
    
    if (!levels.has(level)) {
      levels.set(level, [])
    }
    levels.get(level).push(name)
    
    const dependents = props.steps.filter(step => 
      step.depends && step.depends.includes(name)
    )
    
    dependents.forEach(dep => {
      queue.push({ name: dep.name, level: level + 1 })
    })
  }
  
  const nodeWidth = 300
  const nodeHeight = 160
  const horizontalGap = 150
  const verticalGap = 100
  
  levels.forEach((nodeNames, level) => {
    const x = level * (nodeWidth + horizontalGap)
    
    nodeNames.forEach((nodeName, index) => {
      const y = index * (nodeHeight + verticalGap)
      const node = nodes.value.find(n => n.id === nodeName)
      if (node) {
        node.position = { x, y }
      }
    })
  })
}

// 处理节点点击
const handleNodeClick = (event) => {
  const stepName = event.node.id
  emit('node-click', stepName)
}

// 处理右键菜单
const handleNodeContextMenu = (event) => {
  event.event.preventDefault()
  
  const step = event.node.data.step
  
  contextMenu.value = {
    visible: true,
    x: event.event.clientX,
    y: event.event.clientY,
    step: step,
    nodeId: event.node.id,
  }
}

// 处理菜单操作
const handleAction = async (action) => {
  if (contextMenu.value.step) {
    try {
      await stepAction(props.taskName, contextMenu.value.step.name, action)
    } catch (error) {
      console.error('操作失败:', error)
    }
  }
  closeContextMenu()
}

// 关闭右键菜单
const closeContextMenu = () => {
  contextMenu.value.visible = false
}

// 监听点击事件关闭菜单
const handleClickOutside = (event) => {
  if (contextMenu.value.visible && !event.target.closest('.modern-context-menu')) {
    closeContextMenu()
  }
}

// 监听 steps 变化
watch(
  () => props.steps,
  () => {
    updateGraph()
  },
  { deep: true }
)

onMounted(() => {
  updateGraph()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// 暴露控制方法
defineExpose({
  fitView: () => {
    fitView({ padding: 0.2, duration: 300 })
  },
  zoomIn: () => {
    flowZoomIn({ duration: 300 })
  },
  zoomOut: () => {
    flowZoomOut({ duration: 300 })
  },
  actualSize: () => {
    setViewport({ x: 0, y: 0, zoom: 1 }, { duration: 300 })
  },
})
</script>

<style lang="scss">
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
@import '@vue-flow/controls/dist/style.css';
@import '@vue-flow/minimap/dist/style.css';

.vue-flow-wrapper {
  width: 100%;
  height: 100%;
  position: relative;
  background: linear-gradient(135deg, #f6f8fb 0%, #ffffff 100%);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.modern-vue-flow {
  :deep(.vue-flow__background) {
    background: transparent;
  }

  :deep(.vue-flow__edge) {
    .vue-flow__edge-path {
      stroke-width: 2;
      filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
    }

    &.animated {
      .vue-flow__edge-path {
        stroke-dasharray: 5;
        animation: dashdraw 0.5s linear infinite;
      }
    }
  }

  :deep(.vue-flow__controls) {
    background: white;
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-md);
    padding: var(--spacing-sm);
    border: 1px solid var(--color-border-light);

    button {
      border-radius: var(--radius-base);
      border: 1px solid var(--color-border-light);
      transition: all var(--transition-base);
      width: 32px;
      height: 32px;
      font-size: 14px;

      &:hover {
        background: var(--color-primary);
        border-color: var(--color-primary);
        color: white;
        transform: scale(1.1);
        box-shadow: var(--shadow-sm);
      }
    }
  }

  :deep(.vue-flow__minimap) {
    background: white;
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-md);
    border: 1px solid var(--color-border-light);
    overflow: hidden;

    .vue-flow__minimap-mask {
      fill: rgba(102, 126, 234, 0.08);
      stroke: var(--color-primary);
      stroke-width: 2;
    }

    .vue-flow__minimap-node {
      fill: var(--color-bg-tertiary);
      stroke: var(--color-border-base);
    }
  }
}

@keyframes dashdraw {
  to {
    stroke-dashoffset: -10;
  }
}

// ==========================================
// 现代化右键菜单
// ==========================================
.modern-context-menu {
  position: fixed;
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-xl);
  z-index: 9999;
  min-width: 200px;
  overflow: hidden;
  border: 1px solid var(--color-border-light);
  backdrop-filter: blur(10px);

  .context-menu-header {
    padding: var(--spacing-sm) var(--spacing-base);
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    color: white;
    font-size: 14px;
    font-weight: 600;
    line-height: 1.5;

    .menu-header-icon {
      font-size: 16px;
    }
  }

  .context-menu-body {
    padding: var(--spacing-sm);
  }

  .menu-item {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    padding: var(--spacing-sm) var(--spacing-base);
    cursor: pointer;
    font-size: 14px;
    color: var(--color-text-primary);
    border-radius: var(--radius-md);
    transition: all var(--transition-base);
    font-weight: 500;
    line-height: 1.5;
    min-height: 40px;

    .el-icon {
      font-size: 16px;
    }

    &:hover:not(.disabled) {
      background: var(--color-bg-tertiary);
      transform: translateX(4px);
    }

    &.success {
      &:hover {
        background: rgba(82, 196, 26, 0.1);
        color: var(--color-success);
      }
    }

    &.warning {
      &:hover {
        background: rgba(250, 173, 20, 0.1);
        color: var(--color-warning);
      }
    }

    &.danger {
      &:hover {
        background: rgba(255, 77, 79, 0.1);
        color: var(--color-danger);
      }
    }

    &.disabled {
      color: var(--color-text-disabled);
      cursor: not-allowed;
      opacity: 0.6;
    }
  }
}

.context-menu-fade-enter-active {
  animation: contextMenuFadeIn var(--transition-fast);
}

.context-menu-fade-leave-active {
  animation: contextMenuFadeOut 100ms;
}

@keyframes contextMenuFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

@keyframes contextMenuFadeOut {
  from {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
  to {
    opacity: 0;
    transform: scale(0.95) translateY(-10px);
  }
}

// ==========================================
// 图形统计信息
// ==========================================
.graph-stats {
  position: absolute;
  top: var(--spacing-base);
  left: var(--spacing-base);
  display: flex;
  gap: var(--spacing-sm);
  z-index: 10;

  .stat-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: var(--spacing-sm) var(--spacing-base);
    background: rgba(255, 255, 255, 0.95);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);
    backdrop-filter: blur(10px);
    border: 1px solid var(--color-border-light);

    .stat-label {
      font-size: 13px;
      color: var(--color-text-secondary);
      font-weight: 500;
      line-height: 1.5;
    }

    .stat-value {
      font-size: 14px;
      color: var(--color-primary);
      font-weight: 700;
      line-height: 1.5;
    }
  }
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 768px) {
  .modern-vue-flow {
    :deep(.vue-flow__controls) {
      padding: var(--spacing-xs);
      
      button {
        width: 28px;
        height: 28px;
        font-size: 13px;
      }
    }

    :deep(.vue-flow__minimap) {
      width: 120px;
      height: 80px;
    }
  }

  .graph-stats {
    top: var(--spacing-sm);
    left: var(--spacing-sm);
    gap: var(--spacing-xs);
    
    .stat-item {
      padding: var(--spacing-xs) var(--spacing-sm);
      gap: 4px;

      .stat-label {
        font-size: 12px;
      }

      .stat-value {
        font-size: 13px;
      }
    }
  }
}
</style>