<template>
  <PageContainer>
    <template #toolbar>
      <div class="page-toolbar">
        <div class="toolbar-left">
          <el-button type="primary" @click="showCreateDialog = true" size="large">
            <el-icon><Plus /></el-icon>
            新建流水线
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
          :icon="Connection"
          label="启用中"
          :value="stats.enabled"
          variant="active"
          clickable
        />
        <StatCard
          :icon="CircleCloseFilled"
          label="已禁用"
          :value="stats.disabled"
          variant="disabled"
          clickable
        />
        <StatCard
          :icon="Edit"
          label="Jinja2模板"
          :value="stats.jinja"
          variant="jinja"
          clickable
        />
      </div>
    </template>

    <!-- 卡片视图 -->
    <div v-if="viewMode === 'grid'" v-loading="loading">
      <EmptyState v-if="items.length === 0 && !loading" description="暂无流水线">
        <template #actions>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            创建第一个流水线
          </el-button>
        </template>
      </EmptyState>
      <CardGrid v-else :min-column-width="360">
        <div
          v-for="pipeline in items"
          :key="pipeline.name"
          class="pipeline-card"
          :class="{ disabled: pipeline.disable }"
          @click="handleViewDetail(pipeline)"
        >
          <div class="card-header">
            <div class="pipeline-icon-wrapper">
              <el-icon class="pipeline-icon"><Connection /></el-icon>
            </div>
            <div class="header-info">
              <h3 class="pipeline-name">{{ pipeline.name }}</h3>
              <el-tag size="small" :type="pipeline.tplType === 'jinja2' ? 'warning' : 'info'">
                {{ pipeline.tplType }}
              </el-tag>
            </div>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, pipeline)">
              <el-icon class="more-icon" @click.stop><MoreFilled /></el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="detail">
                    <el-icon><View /></el-icon>
                    查看详情
                  </el-dropdown-item>
                  <el-dropdown-item command="edit">
                    <el-icon><Edit /></el-icon>
                    编辑配置
                  </el-dropdown-item>
                  <el-dropdown-item command="run" divided>
                    <el-icon><VideoPlay /></el-icon>
                    立即运行
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    <el-icon><Delete /></el-icon>
                    删除流水线
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>

          <div class="card-body">
            <div class="pipeline-desc">
              <el-icon class="desc-icon"><Document /></el-icon>
              <p class="desc-text">{{ pipeline.desc || '暂无描述' }}</p>
            </div>

            <div class="pipeline-status">
              <div class="status-item">
                <el-icon><Switch /></el-icon>
                <span>{{ pipeline.disable ? '已禁用' : '已启用' }}</span>
              </div>
            </div>

            <div class="card-actions">
              <el-button size="small" type="primary" plain @click.stop="handleRun(pipeline)">
                <el-icon><VideoPlay /></el-icon>
                运行
              </el-button>
              <el-button size="small" plain @click.stop="handleEdit(pipeline)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
            </div>
          </div>

          <!-- 禁用遮罩 -->
          <div v-if="pipeline.disable" class="disabled-overlay">
            <el-icon><Lock /></el-icon>
            <span>已禁用</span>
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
        <el-table-column prop="name" label="流水线名称" min-width="200" fixed>
          <template #default="{ row }">
            <div class="pipeline-name-cell">
              <el-icon class="pipeline-cell-icon"><Connection /></el-icon>
              <span class="name-text">{{ row.name }}</span>
              <el-tag v-if="row.disable" size="small" type="info">禁用</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="描述" min-width="200">
          <template #default="{ row }">
            <span class="desc-text">{{ row.desc || '暂无描述' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="模板类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.tplType === 'jinja2' ? 'warning' : 'info'">
              {{ row.tplType }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.disable ? 'info' : 'success'" size="small">
              {{ row.disable ? '已禁用' : '已启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button size="small" @click="handleViewDetail(row)">
                <el-icon><View /></el-icon>
                详情
              </el-button>
              <el-button size="small" type="primary" plain @click="handleRun(row)">
                <el-icon><VideoPlay /></el-icon>
                运行
              </el-button>
              <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
                <el-button size="small">
                  更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="edit">编辑配置</el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除流水线</el-dropdown-item>
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

    <!-- 创建流水线对话框 -->
    <PipelineFormDialog
      v-model="showCreateDialog"
      @success="handlePipelineCreated"
    />

    <!-- 编辑流水线对话框 -->
    <PipelineFormDialog
      v-model="showEditDialog"
      :pipeline-data="selectedPipeline"
      @success="handlePipelineUpdated"
    />

    <!-- 流水线详情对话框 -->
    <PipelineDetailDialog
      v-model="showDetailDialog"
      :pipeline-name="selectedPipelineName"
    />

    <!-- 运行流水线对话框 -->
    <RunPipelineDialog
      v-model="showRunDialog"
      :pipeline-name="selectedPipelineName"
      @success="handlePipelineRun"
    />
  </PageContainer>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  Connection,
  CircleCloseFilled,
  Edit,
  Plus,
  Grid,
  List,
  View,
  VideoPlay,
  Delete,
  ArrowDown,
  Document,
  Switch,
  Lock,
  MoreFilled,
} from '@element-plus/icons-vue'
import { API_ENDPOINTS } from '@/config'
import { deletePipeline, getPipelineDetail } from '@/api/pipeline'
import { useWebSocketList } from '@/composables/useWebSocket'
import { usePipelineStats } from '@/composables/useStats'
import { useViewMode } from '@/composables/useViewMode'
import { PageContainer, StatCard, CardGrid, EmptyState } from '@/components/base'
import PipelineFormDialog from '@/components/PipelineFormDialog.vue'
import PipelineDetailDialog from '@/components/PipelineDetailDialog.vue'
import RunPipelineDialog from '@/components/RunPipelineDialog.vue'

// 视图模式管理
const { viewMode } = useViewMode('grid', 'pipeline-view-mode')

// 对话框状态
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showDetailDialog = ref(false)
const showRunDialog = ref(false)
const selectedPipelineName = ref('')
const selectedPipeline = ref(null)

// WebSocket 列表管理
const {
  items,
  loading,
  pagination,
  refresh,
  changePage,
  changePageSize,
  connect,
} = useWebSocketList(API_ENDPOINTS.pipeline, {
  onMessage: (response) => {
    // 自定义消息处理，映射pipelines字段
    if (response.data && response.data.pipelines) {
      // 已在 useWebSocketList 中处理
    }
  }
})

// 统计数据
const { stats } = usePipelineStats(items)

/**
 * 查看流水线详情
 */
const handleViewDetail = (pipeline) => {
  selectedPipelineName.value = pipeline.name
  showDetailDialog.value = true
}

/**
 * 运行流水线
 */
const handleRun = (pipeline) => {
  selectedPipelineName.value = pipeline.name
  showRunDialog.value = true
}

/**
 * 编辑流水线
 */
const handleEdit = async (pipeline) => {
  await loadPipelineForEdit(pipeline.name)
}

/**
 * 执行命令操作
 */
const handleCommand = async (command, row) => {
  switch (command) {
    case 'detail':
      handleViewDetail(row)
      break
    case 'edit':
      await loadPipelineForEdit(row.name)
      break
    case 'run':
      handleRun(row)
      break
    case 'delete':
      await deletePipeline(row.name)
      break
  }
}

/**
 * 加载流水线数据用于编辑
 */
const loadPipelineForEdit = async (pipelineName) => {
  try {
    const response = await getPipelineDetail(pipelineName)
    selectedPipeline.value = response.data
    showEditDialog.value = true
  } catch (error) {
    console.error('获取流水线详情失败:', error)
  }
}

/**
 * 流水线创建成功回调
 */
const handlePipelineCreated = () => {
  // 列表会自动通过WebSocket更新
}

/**
 * 流水线更新成功回调
 */
const handlePipelineUpdated = () => {
  selectedPipeline.value = null
  // 列表会自动通过WebSocket更新
}

/**
 * 流水线运行成功回调
 */
const handlePipelineRun = (taskName) => {
  console.log('流水线已运行，任务名称:', taskName)
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
// 流水线卡片样式
// ==========================================
.pipeline-card {
  background: white;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
  cursor: pointer;
  border: 1px solid var(--color-border-light);
  position: relative;

  &:hover:not(.disabled) {
    transform: translateY(-2px);
    box-shadow: var(--shadow-lg);
    border-color: var(--color-primary);

    .more-icon {
      opacity: 1;
    }
  }

  &.disabled {
    opacity: 0.6;
    
    .card-header {
      background: linear-gradient(135deg, #e5e7eb 0%, #f3f4f6 100%);
    }
  }

  .card-header {
    padding: var(--spacing-base) var(--spacing-lg);
    background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
    border-bottom: 1px solid var(--color-border-light);
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    min-height: 64px;

    .pipeline-icon-wrapper {
      width: 40px;
      height: 40px;
      border-radius: var(--radius-md);
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: var(--shadow-sm);
      flex-shrink: 0;

      .pipeline-icon {
        font-size: 20px;
        color: white;
      }
    }

    .header-info {
      flex: 1;
      min-width: 0;
      display: flex;
      flex-direction: column;
      gap: 4px;

      .pipeline-name {
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

    .more-icon {
      font-size: 20px;
      color: var(--color-text-tertiary);
      cursor: pointer;
      transition: all var(--transition-base);
      opacity: 0;
      padding: 4px;
      border-radius: var(--radius-md);

      &:hover {
        color: var(--color-primary);
        background: rgba(24, 144, 255, 0.1);
      }
    }
  }

  .card-body {
    padding: var(--spacing-base) var(--spacing-lg) var(--spacing-lg);
    display: flex;
    flex-direction: column;
    gap: var(--spacing-base);

    .pipeline-desc {
      display: flex;
      align-items: flex-start;
      gap: var(--spacing-sm);
      padding: var(--spacing-sm) var(--spacing-base);
      background: var(--color-bg-tertiary);
      border-radius: var(--radius-md);
      min-height: 56px;

      .desc-icon {
        color: var(--color-primary);
        font-size: 16px;
        margin-top: 2px;
        flex-shrink: 0;
      }

      .desc-text {
        margin: 0;
        font-size: 13px;
        color: var(--color-text-secondary);
        line-height: 1.6;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
    }

    .pipeline-status {
      display: flex;
      gap: var(--spacing-sm);

      .status-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: var(--color-text-secondary);

        .el-icon {
          font-size: 16px;
          color: var(--color-primary);
        }
      }
    }

    .card-actions {
      display: flex;
      gap: var(--spacing-sm);
      padding-top: var(--spacing-base);
      border-top: 1px solid var(--color-border-light);

      .el-button {
        flex: 1;
      }
    }
  }

  .disabled-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.85);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--spacing-sm);
    font-size: 14px;
    font-weight: 600;
    color: var(--color-text-tertiary);
    pointer-events: none;
    backdrop-filter: blur(2px);

    .el-icon {
      font-size: 32px;
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

  .pipeline-name-cell {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);

    .pipeline-cell-icon {
      color: var(--color-primary);
      font-size: 18px;
    }

    .name-text {
      font-weight: 600;
      color: var(--color-text-primary);
      font-size: 14px;
      flex: 1;
    }
  }

  .desc-text {
    color: var(--color-text-secondary);
    font-size: 13px;
    line-height: 1.6;
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

  .pipeline-card {
    .card-header {
      padding: var(--spacing-sm) var(--spacing-base);
      gap: var(--spacing-xs);
      min-height: 60px;

      .pipeline-icon-wrapper {
        width: 36px;
        height: 36px;

        .pipeline-icon {
          font-size: 18px;
        }
      }

      .header-info {
        gap: 4px;
        
        .pipeline-name {
          font-size: 14px;
        }
      }
      
      .more-icon {
        font-size: 18px;
        opacity: 1;
        padding: 4px;
      }
    }

    .card-body {
      padding: var(--spacing-sm) var(--spacing-base) var(--spacing-base);
      gap: var(--spacing-sm);
      
      .pipeline-desc {
        padding: var(--spacing-sm);
        min-height: 52px;
        
        .desc-icon {
          font-size: 14px;
        }
        
        .desc-text {
          font-size: 13px;
        }
      }
      
      .pipeline-status {
        .status-item {
          font-size: 12px;
          gap: 4px;
          
          .el-icon {
            font-size: 14px;
          }
        }
      }
      
      .card-actions {
        gap: var(--spacing-xs);
        padding-top: var(--spacing-sm);
      }
    }
    
    .disabled-overlay {
      font-size: 13px;
      gap: var(--spacing-sm);
      
      .el-icon {
        font-size: 28px;
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