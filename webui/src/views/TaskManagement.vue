<template>
  <PageContainer>
    <template #toolbar>
      <div class="page-toolbar">
        <div class="toolbar-left">
          <el-button type="primary" @click="showCreateDialog = true" size="large">
            <el-icon><Plus /></el-icon>
            新建任务
          </el-button>
        </div>
        
        <div class="toolbar-right">
          <el-radio-group v-model="viewMode" size="default">
            <el-radio-button value="grid">
              <el-icon><Grid /></el-icon>
              卡片
            </el-radio-button>
            <el-radio-button value="table">
              <el-icon><List /></el-icon>
              列表
            </el-radio-button>
          </el-radio-group>
        </div>
      </div>
    </template>

    <template #stats>
      <div class="stats-grid">
        <StatCard
          :icon="VideoPlay"
          label="运行中"
          :value="stats.running"
          variant="running"
          clickable
        />
        <StatCard
          :icon="CircleCheckFilled"
          label="已完成"
          :value="stats.stopped"
          variant="success"
          clickable
        />
        <StatCard
          :icon="CircleCloseFilled"
          label="已失败"
          :value="stats.failed"
          variant="failed"
          clickable
        />
      </div>
    </template>

    <!-- 卡片视图 -->
    <div v-if="viewMode === 'grid'" v-loading="loading">
      <EmptyState v-if="items.length === 0 && !loading" description="暂无任务">
        <template #actions>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            创建第一个任务
          </el-button>
        </template>
      </EmptyState>
      <CardGrid v-else :min-column-width="340">
        <div
          v-for="task in items"
          :key="task.name"
          class="task-card"
          @click="handleViewDetail(task)"
        >
          <div class="card-header">
            <div class="task-status-indicator" :class="`status-${task.state}`"></div>
            <h3 class="task-name">{{ task.name }}</h3>
          </div>

          <div class="card-body">
            <div class="task-meta">
              <InfoItem
                :icon="Connection"
                label="步骤数"
                :value="task.count"
                layout="horizontal"
                hoverable
              />
              <InfoItem
                :icon="Clock"
                label="开始时间"
                :value="formatTime(task.time?.start)"
                value-class="time"
                layout="horizontal"
                hoverable
              />
              <InfoItem
                :icon="Calendar"
                label="结束时间"
                :value="formatTime(task.time?.end)"
                value-class="time"
                layout="horizontal"
                hoverable
              />
            </div>

            <div class="task-status-section">
              <StatusTag :status="task.state" size="large" effect="plain" show-icon />
              <div class="action-buttons">
                <el-button size="small" @click.stop="handleViewDetail(task)">
                  <el-icon><View /></el-icon>
                </el-button>
                <el-button size="small" @click.stop="handleCommand('dump', task)">
                  <el-icon><Download /></el-icon>
                </el-button>
                <el-button v-if="task.state === 'running'" size="small" type="warning" @click.stop="handleCommand('kill', task)">
                  <el-icon><VideoPause /></el-icon>
                </el-button>
                <el-button size="small" type="danger" @click.stop="handleCommand('delete', task)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </CardGrid>
    </div>

    <!-- 表格视图 -->
    <div v-else class="table-view">
      <el-table
        :data="items"
        stripe
        v-loading="loading"
        class="modern-table"
      >
        <el-table-column prop="name" label="任务名称" min-width="200" fixed>
          <template #default="{ row }">
            <div class="task-name-cell">
              <div class="status-dot" :class="`status-${row.state}`"></div>
              <span class="name-text">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="count" label="步骤数" width="100" align="center" />
        <el-table-column label="开始时间" min-width="180">
          <template #default="{ row }">
            {{ formatTime(row.time?.start) }}
          </template>
        </el-table-column>
        <el-table-column label="结束时间" min-width="180">
          <template #default="{ row }">
            {{ formatTime(row.time?.end) }}
          </template>
        </el-table-column>
        <el-table-column label="执行状态" width="140" align="center">
          <template #default="{ row }">
            <StatusTag :status="row.state" effect="plain" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button size="small" @click="handleViewDetail(row)">
                <el-icon><View /></el-icon>
                详情
              </el-button>
              <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
                <el-button size="small">
                  更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="dump">导出配置</el-dropdown-item>
                    <el-dropdown-item v-if="row.state === 'running'" command="kill" divided>
                      强制停止
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除任务</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <template #footer>
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :page-sizes="[15, 30, 45, 60]"
          :total="pagination.total * pagination.size"
          :background="true"
          layout="sizes, prev, pager, next, jumper"
          @size-change="changePageSize"
          @current-change="changePage"
        />
      </div>
    </template>

    <!-- 创建任务对话框 -->
    <TaskFormDialog v-model="showCreateDialog" @success="handleTaskCreated" />

    <!-- 任务详情对话框 -->
    <TaskDetailDialog
      v-model="showDetailDialog"
      :task-name="selectedTaskName"
    />
  </PageContainer>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  VideoPlay,
  CircleCheckFilled,
  CircleCloseFilled,
  Plus,
  Grid,
  List,
  View,
  Download,
  VideoPause,
  Delete,
  ArrowDown,
  Connection,
  Clock,
  Calendar,
} from '@element-plus/icons-vue'
import { API_ENDPOINTS } from '@/config'
import { taskAction, deleteTask, dumpTask } from '@/api/task'
import { formatTime } from '@/utils/format'
import { useWebSocketList } from '@/composables/useWebSocket'
import { useTaskStats } from '@/composables/useStats'
import { useViewMode } from '@/composables/useViewMode'
import { PageContainer, StatCard, CardGrid, StatusTag, EmptyState, InfoItem } from '@/components/base'
import TaskFormDialog from '@/components/TaskFormDialog.vue'
import TaskDetailDialog from '@/components/TaskDetailDialog.vue'

// 视图模式管理
const { viewMode } = useViewMode('grid', 'task-view-mode')

// 对话框状态
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedTaskName = ref('')

// WebSocket 列表管理
const {
  items,
  loading,
  pagination,
  refresh,
  changePage,
  changePageSize,
  connect,
} = useWebSocketList(API_ENDPOINTS.task)

// 统计数据
const { stats } = useTaskStats(items)

/**
 * 查看任务详情
 */
const handleViewDetail = (task) => {
  selectedTaskName.value = task.name
  showDetailDialog.value = true
}

/**
 * 执行命令操作
 */
const handleCommand = async (command, row) => {
  switch (command) {
    case 'detail':
      handleViewDetail(row)
      break
    case 'dump':
      await dumpTask(row.name)
      break
    case 'kill':
      await taskAction(row.name, 'kill')
      break
    case 'delete':
      await deleteTask(row.name)
      break
  }
}

/**
 * 任务创建成功回调
 */
const handleTaskCreated = (taskName) => {
  selectedTaskName.value = taskName
  showDetailDialog.value = true
}

// 初始化WebSocket连接
onMounted(() => {
  connect()
})
</script>

<style lang="scss" scoped>
// ==========================================
// 工具栏样式
// ==========================================
.page-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-base) var(--spacing-lg);
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border-light);
  min-height: 64px;

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }

  .toolbar-right {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }
}

// ==========================================
// 统计卡片网格
// ==========================================
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-base);
}

// ==========================================
// 任务卡片样式
// ==========================================
.task-card {
  background: white;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
  cursor: pointer;
  border: 1px solid var(--color-border-light);

  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-lg);
    border-color: var(--color-primary);
  }

  .card-header {
    padding: var(--spacing-base) var(--spacing-lg);
    background: linear-gradient(135deg, #f6f8fb 0%, #fafbfc 100%);
    border-bottom: 1px solid var(--color-border-light);
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    min-height: 56px;

    .task-status-indicator {
      width: 8px;
      height: 8px;
      border-radius: var(--radius-full);
      flex-shrink: 0;

      &.status-running {
        background: var(--color-primary);
        box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
        animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
      }

      &.status-stopped {
        background: var(--color-success);
      }

      &.status-failed {
        background: var(--color-danger);
      }

      &.status-pending {
        background: var(--color-warning);
      }
    }

    .task-name {
      flex: 1;
      margin: 0;
      font-size: 15px;
      font-weight: 600;
      color: var(--color-text-primary);
      line-height: 1.5;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .card-body {
    padding: var(--spacing-base) var(--spacing-lg) var(--spacing-lg);

    .task-meta {
      display: flex;
      flex-direction: column;
      gap: var(--spacing-xs);
      margin-bottom: var(--spacing-base);
    }

    .task-status-section {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-top: var(--spacing-base);
      border-top: 1px solid var(--color-border-light);

      .action-buttons {
        display: flex;
        gap: var(--spacing-xs);

        .el-button {
          transition: all var(--transition-base);

          &:hover {
            transform: translateY(-2px);
          }
        }
      }
    }
  }
}

// ==========================================
// 表格视图样式
// ==========================================
.table-view {
  background: white;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border-light);

  .modern-table {
    :deep(.el-table__header-wrapper) {
      th {
        background: linear-gradient(135deg, #f6f8fb 0%, #fafbfc 100%);
        color: var(--color-text-primary);
        font-weight: 600;
        font-size: 13px;
        text-transform: uppercase;
        letter-spacing: 0.5px;
      }
    }

    :deep(.el-table__body-wrapper) {
      tr {
        transition: all var(--transition-base);

        &:hover > td {
          background: #f0f8ff !important;
        }
      }
    }
  }

  .task-name-cell {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);

    .status-dot {
      width: 8px;
      height: 8px;
      border-radius: var(--radius-full);
      flex-shrink: 0;

      &.status-running {
        background: var(--color-primary);
        box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
        animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
      }

      &.status-stopped {
        background: var(--color-success);
      }

      &.status-failed {
        background: var(--color-danger);
      }

      &.status-pending {
        background: var(--color-warning);
      }
    }

    .name-text {
      font-weight: 600;
      color: var(--color-text-primary);
      font-size: 14px;
    }
  }

  .action-buttons {
    display: flex;
    gap: var(--spacing-xs);
  }
}

// 分页区域
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: var(--spacing-base) var(--spacing-lg);
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border-light);
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 768px) {
  .page-toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-sm);
    padding: var(--spacing-sm) var(--spacing-base);
    min-height: auto;

    .toolbar-left,
    .toolbar-right {
      width: 100%;
    }

    .toolbar-left {
      justify-content: space-between;
    }
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-sm);
  }

  .task-card {
    .card-header {
      padding: var(--spacing-sm) var(--spacing-base);
      gap: var(--spacing-xs);
      min-height: 52px;
      
      .task-status-indicator {
        width: 6px;
        height: 6px;
      }
      
      .task-name {
        font-size: 14px;
      }
    }

    .card-body {
      padding: var(--spacing-sm) var(--spacing-base) var(--spacing-base);
      
      .task-meta {
        gap: var(--spacing-xs);
        margin-bottom: var(--spacing-sm);
      }
      
      .task-status-section {
        padding-top: var(--spacing-sm);
        
        .action-buttons {
          gap: var(--spacing-xs);
        }
      }
    }
  }

  .pagination-wrapper {
    padding: var(--spacing-sm) var(--spacing-base);

    :deep(.el-pagination) {
      justify-content: center;
      
      .el-pagination__sizes,
      .el-pagination__jump {
        display: none;
      }
    }
  }
}
</style>