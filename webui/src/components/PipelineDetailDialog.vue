<template>
  <el-dialog
    v-model="dialogVisible"
    :close-on-click-modal="false"
    :show-close="false"
    @close="handleClose"
    fullscreen
    class="pipeline-detail-dialog"
  >
    <template #header>
      <DialogHeader
        :icon="Connection"
        :title="pipelineName"
        @close="dialogVisible = false"
      />
    </template>

    <div class="dialog-content">
      <el-tabs v-model="activeTab" type="border-card" class="content-tabs" v-loading="loading">
        <!-- 详情标签页 -->
        <el-tab-pane name="detail">
          <template #label>
            <div class="tab-label">
              <el-icon><Document /></el-icon>
              <span>流水线配置</span>
            </div>
          </template>
          
          <div class="tab-content detail-content">
            <!-- 基本信息卡片 -->
            <el-card shadow="hover" class="detail-info-card">
              <template #header>
                <SectionHeader
                  :icon="InfoFilled"
                  title="基本信息"
                />
              </template>
              <div class="info-grid" v-if="pipelineData">
                <InfoItem
                  label="流水线名称"
                  :value="pipelineData.name"
                  layout="vertical"
                  hoverable
                />
                <InfoItem label="模板类型" layout="vertical" hoverable>
                  <template #value>
                    <el-tag :type="pipelineData.tplType === 'jinja2' ? 'warning' : 'info'" size="small">
                      {{ pipelineData.tplType }}
                    </el-tag>
                  </template>
                </InfoItem>
                <InfoItem
                  label="流水线描述"
                  :value="pipelineData.desc || '暂无描述'"
                  value-class="desc"
                  layout="vertical"
                  hoverable
                />
                <InfoItem label="状态" layout="vertical" hoverable>
                  <template #value>
                    <el-tag :type="pipelineData.disable ? 'info' : 'success'" size="small">
                      {{ pipelineData.disable ? '已禁用' : '已启用' }}
                    </el-tag>
                  </template>
                </InfoItem>
              </div>
            </el-card>

            <!-- YAML内容编辑器 -->
            <el-card shadow="hover" class="editor-card">
              <template #header>
                <SectionHeader
                  :icon="EditPen"
                  title="YAML 配置内容"
                  tag="只读模式"
                  tag-type="info"
                />
              </template>
              <MonacoEditor
                v-if="pipelineData"
                v-model="editorContent"
                language="yaml"
                theme="vs-dark"
                :read-only="true"
                :show-toolbar="false"
                height="100%"
              />
            </el-card>
          </div>
        </el-tab-pane>

        <!-- 任务列表标签页 -->
        <el-tab-pane name="tasks">
          <template #label>
            <div class="tab-label">
              <el-icon><List /></el-icon>
              <span>执行历史</span>
              <el-badge :value="taskItems.length" :max="99" v-if="taskItems.length > 0" class="tab-badge" />
            </div>
          </template>
          
          <div class="tab-content tasks-content">
            <!-- 任务列表工具栏 -->
            <div class="tasks-toolbar">
              <div class="toolbar-left">
                <h3 class="toolbar-title">
                  <el-icon><Clock /></el-icon>
                  执行记录 ({{ taskPagination.total * taskPagination.size }})
                </h3>
              </div>
              <div class="toolbar-right">
                <el-pagination
                  v-model:current-page="taskPagination.current"
                  v-model:page-size="taskPagination.size"
                  :page-sizes="[10, 20, 30, 50]"
                  :total="taskPagination.total * taskPagination.size"
                  :background="true"
                  layout="sizes, prev, pager, next"
                  @size-change="handleTaskSizeChange"
                  @current-change="handleTaskPageChange"
                />
              </div>
            </div>

            <!-- 任务卡片列表 -->
            <el-scrollbar class="tasks-scrollbar">
              <EmptyState v-if="taskItems.length === 0" :icon="Clock" description="暂无执行记录" />
              <div v-else class="tasks-grid">
                <div
                  v-for="task in taskItems"
                  :key="task.taskName"
                  class="task-history-card"
                  @click="handleTaskClick(task.taskName)"
                >
                  <div class="card-left">
                    <div class="status-indicator" :class="`status-${task.state}`"></div>
                    <div class="task-info">
                      <h4 class="task-name">{{ task.taskName }}</h4>
                      <div class="task-meta">
                        <div class="meta-item">
                          <el-icon><Clock /></el-icon>
                          <span>开始: {{ formatTime(task.time?.start) }}</span>
                        </div>
                        <div class="meta-item">
                          <el-icon><Calendar /></el-icon>
                          <span>结束: {{ formatTime(task.time?.end) }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="card-right">
                    <StatusTag :status="task.state" effect="plain" size="large" show-icon />
                    <el-icon class="arrow-icon"><ArrowRight /></el-icon>
                  </div>
                </div>
              </div>
            </el-scrollbar>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 任务详情对话框 -->
    <TaskDetailDialog
      v-model="showTaskDetail"
      :task-name="selectedTaskName"
    />
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { 
  Connection, 
  Document, 
  InfoFilled, 
  EditPen, 
  List, 
  Clock,
  Calendar,
  ArrowRight,
} from '@element-plus/icons-vue'
import { getPipelineDetail } from '@/api/pipeline'
import { useWebSocketList } from '@/composables/useWebSocket'
import { useTabs } from '@/composables/useViewMode'
import { API_ENDPOINTS } from '@/config'
import { formatTime } from '@/utils/format'
import { DialogHeader, SectionHeader, InfoItem, MonacoEditor, EmptyState, StatusTag } from '@/components/base'
import TaskDetailDialog from './TaskDetailDialog.vue'

const props = defineProps({
  modelValue: Boolean,
  pipelineName: String,
})

const emit = defineEmits(['update:modelValue'])

const dialogVisible = ref(false)
const loading = ref(false)
const pipelineData = ref(null)
const editorContent = ref('')
const showTaskDetail = ref(false)
const selectedTaskName = ref('')

// 标签页管理
const { activeTab, switchTab } = useTabs('detail')

// 任务列表WebSocket管理
const {
  items: taskItems,
  pagination: taskPagination,
  refresh: refreshTasks,
  changePage: handleTaskPageChange,
  changePageSize: handleTaskSizeChange,
  connect: connectTasks,
  close: closeTasks,
} = useWebSocketList(`${API_ENDPOINTS.pipeline}/${props.pipelineName}/build`)

watch(
  () => props.modelValue,
  async (val) => {
    dialogVisible.value = val
    if (val && props.pipelineName) {
      await loadPipelineDetail()
    }
  }
)

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
  if (!val) {
    handleClose()
  }
})

watch(activeTab, (val) => {
  if (val === 'tasks') {
    connectTasks()
  }
})

/**
 * 加载流水线详情
 */
const loadPipelineDetail = async () => {
  loading.value = true
  try {
    const response = await getPipelineDetail(props.pipelineName)
    pipelineData.value = response.data
    editorContent.value = response.data.content || ''
  } catch (error) {
    console.error('获取流水线详情失败:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 处理任务点击
 */
const handleTaskClick = (taskName) => {
  selectedTaskName.value = taskName
  showTaskDetail.value = true
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  closeTasks()
  pipelineData.value = null
  editorContent.value = ''
  activeTab.value = 'detail'
}
</script>

<style lang="scss" scoped>
.pipeline-detail-dialog {
  :deep(.el-dialog) {
    background: var(--color-bg-secondary);
  }

  :deep(.el-dialog__header) {
    padding: 0 !important;
    margin: 0 !important;
    border: none !important;
  }

  :deep(.el-dialog__body) {
    padding: 0;
    height: calc(100vh - 68px);
    overflow: hidden;
  }
}

// ==========================================
// 对话框内容
// ==========================================
.dialog-content {
  height: 100%;
  padding: var(--spacing-lg);

  .content-tabs {
    height: 100%;
    display: flex;
    flex-direction: column;
    background: white;
    border-radius: var(--radius-md);
    overflow: hidden;
    box-shadow: var(--shadow-sm);

    :deep(.el-tabs__header) {
      margin: 0;
      padding: var(--spacing-sm) var(--spacing-lg);
      background: linear-gradient(135deg, #f6f8fb 0%, #fafbfc 100%);
      border-bottom: 1px solid var(--color-border-light);
      height: 56px;
    }

    :deep(.el-tabs__content) {
      flex: 1;
      overflow: hidden;
    }

    :deep(.el-tab-pane) {
      height: 100%;
    }

    .tab-label {
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);
      font-weight: 500;
      font-size: 14px;

      .tab-badge {
        margin-left: 6px;
      }
    }

    .tab-content {
      height: 100%;
      padding: var(--spacing-lg);
      overflow-y: auto;
    }
  }
}

// ==========================================
// 详情页内容
// ==========================================
.detail-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-base);

  .detail-info-card {
    flex-shrink: 0;

    .info-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: var(--spacing-xs);
    }
  }

  .editor-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    :deep(.el-card__body) {
      flex: 1;
      padding: 0;
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }

    :deep(.codemirror-editor-wrapper) {
      flex: 1;
      min-height: 0;
      border: none;
      box-shadow: none;
    }
  }
}

// ==========================================
// 任务历史内容
// ==========================================
.tasks-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-base);
  overflow: hidden;

  .tasks-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--spacing-base) var(--spacing-lg);
    background: var(--color-bg-secondary);
    border-radius: var(--radius-md);
    flex-shrink: 0;
    min-height: 56px;

    .toolbar-title {
      margin: 0;
      font-size: 15px;
      font-weight: 600;
      color: var(--color-text-primary);
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);
      line-height: 1.5;

      .el-icon {
        color: var(--color-primary);
        font-size: 18px;
      }
    }
  }

  .tasks-scrollbar {
    flex: 1;
  }

  .tasks-grid {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .task-history-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--spacing-base) var(--spacing-lg);
    background: white;
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-base);
    box-shadow: var(--shadow-sm);
    min-height: 80px;

    &:hover {
      box-shadow: var(--shadow-lg);
      transform: translateY(-2px);
      border-color: var(--color-primary);

      .arrow-icon {
        transform: translateX(4px);
      }
    }

    .card-left {
      flex: 1;
      display: flex;
      align-items: center;
      gap: var(--spacing-base);
      min-width: 0;

      .status-indicator {
        width: 4px;
        height: 56px;
        border-radius: var(--radius-base);
        flex-shrink: 0;

        &.status-running {
          background: linear-gradient(180deg, var(--color-primary) 0%, var(--color-primary-light) 100%);
          animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
        }

        &.status-stopped {
          background: linear-gradient(180deg, var(--color-success) 0%, #84fab0 100%);
        }

        &.status-failed {
          background: linear-gradient(180deg, var(--color-danger) 0%, #fa709a 100%);
        }

        &.status-pending {
          background: linear-gradient(180deg, var(--color-warning) 0%, #fed6e3 100%);
        }
      }

      .task-info {
        flex: 1;
        min-width: 0;

        .task-name {
          margin: 0 0 var(--spacing-xs) 0;
          font-size: 12px;
          font-weight: 600;
          color: var(--color-text-primary);
          line-height: 1.3;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .task-meta {
          display: flex;
          flex-direction: column;
          gap: 2px;

          .meta-item {
            display: flex;
            align-items: center;
            gap: 4px;
            font-size: 10px;
            color: var(--color-text-secondary);
            line-height: 1.3;

            .el-icon {
              font-size: 11px;
              color: var(--color-primary);
            }
          }
        }
      }
    }

    .card-right {
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);
      flex-shrink: 0;

      .arrow-icon {
        font-size: 16px;
        color: var(--color-text-tertiary);
        transition: transform var(--transition-fast);
      }
    }
  }
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 1200px) {
  .detail-info-card {
    .info-grid {
      grid-template-columns: 1fr !important;
      gap: 1px;
    }
  }
}

@media (max-width: 768px) {
  .dialog-content {
    padding: var(--spacing-xs);
  }

  .tasks-toolbar {
    flex-direction: column;
    gap: var(--spacing-sm);
    align-items: stretch;
    padding: var(--spacing-xs);

    .toolbar-title {
      font-size: 11px;
      
      .el-icon {
        font-size: 12px;
      }
    }

    .toolbar-right {
      :deep(.el-pagination) {
        justify-content: center;

        .el-pagination__sizes {
          display: none;
        }
      }
    }
  }
  
  .task-history-card {
    padding: var(--spacing-xs);
    
    .card-left {
      gap: var(--spacing-xs);
      
      .status-indicator {
        width: 6px;
        height: 36px;
      }
      
      .task-info {
        .task-name {
          font-size: 11px;
        }
        
        .task-meta {
          .meta-item {
            font-size: 9px;
            gap: 2px;
            
            .el-icon {
              font-size: 10px;
            }
          }
        }
      }
    }
  }
}
</style>