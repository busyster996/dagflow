<template>
  <PageContainer>
    <template #toolbar v-if="uploads.size > 0">
      <div class="page-toolbar">
        <div class="toolbar-left">
          <el-button-group>
            <el-button type="primary" @click="startAll">
              <el-icon><VideoPlay /></el-icon>
              开始全部
            </el-button>
            <el-button @click="pauseAll">
              <el-icon><VideoPause /></el-icon>
              暂停全部
            </el-button>
            <el-button @click="resumeAll">
              <el-icon><Refresh /></el-icon>
              恢复全部
            </el-button>
            <el-button type="danger" @click="clearAll">
              <el-icon><Delete /></el-icon>
              清空队列
            </el-button>
          </el-button-group>
        </div>
        
        <div class="toolbar-right">
          <div class="stats-badges">
            <el-badge :value="statusCount.waiting" :max="99" type="info">
              <el-tag>排队</el-tag>
            </el-badge>
            <el-badge :value="statusCount.uploading" :max="99" type="success">
              <el-tag type="success">上传中</el-tag>
            </el-badge>
            <el-badge :value="statusCount.paused" :max="99" type="warning">
              <el-tag type="warning">暂停</el-tag>
            </el-badge>
            <el-badge :value="statusCount.error" :max="99" type="danger">
              <el-tag type="danger">失败</el-tag>
            </el-badge>
          </div>
        </div>
      </div>
    </template>

    <div class="content-layout">
      <!-- 左侧配置面板 -->
      <div class="config-panel">
        <el-card shadow="hover" class="config-card">
          <template #header>
            <SectionHeader
              :icon="Setting"
              title="上传配置"
            />
          </template>
          
          <el-form label-width="auto" size="default" label-position="top">
            <el-form-item label="任务 ID">
              <el-input v-model="uploadConfig.taskId" placeholder="task-123456">
                <template #prefix>
                  <el-icon><Key /></el-icon>
                </template>
              </el-input>
              <template #help>
                <span class="form-help">用于标识上传文件所属的任务</span>
              </template>
            </el-form-item>
            
            <el-form-item label="分块大小 (MB)">
              <el-input-number
                v-model="uploadConfig.chunkSizeMB"
                :min="1"
                :max="100"
                :step="1"
                controls-position="right"
                style="width: 100%"
              />
            </el-form-item>
            
            <el-form-item label="单文件并行数">
              <el-input-number
                v-model="uploadConfig.parallelUploads"
                :min="1"
                :max="10"
                :step="1"
                controls-position="right"
                style="width: 100%"
              />
              <template #help>
                <span class="form-help">每个文件同时上传的分块数量</span>
              </template>
            </el-form-item>
            
            <el-form-item label="队列并发数">
              <el-input-number
                v-model="uploadConfig.queueConcurrency"
                :min="1"
                :max="10"
                :step="1"
                controls-position="right"
                style="width: 100%"
              />
              <template #help>
                <span class="form-help">同时上传的文件数量</span>
              </template>
            </el-form-item>
          </el-form>

          <el-divider />

          <div class="config-summary">
            <h4 class="summary-title">当前配置</h4>
            <div class="summary-items">
              <InfoItem
                label="分块大小"
                :value="`${uploadConfig.chunkSizeMB} MB`"
                layout="horizontal"
              />
              <InfoItem
                label="单文件并行"
                :value="`${uploadConfig.parallelUploads} 块`"
                layout="horizontal"
              />
              <InfoItem
                label="队列并发"
                :value="`${uploadConfig.queueConcurrency} 个`"
                layout="horizontal"
              />
            </div>
          </div>
        </el-card>
      </div>

      <!-- 右侧内容区 -->
      <div class="main-panel">
        <el-card shadow="hover" class="main-card">
          <el-tabs v-model="activeTab" type="card" class="upload-tabs">
            <!-- 上传队列标签页 -->
            <el-tab-pane name="queue">
              <template #label>
                <div class="tab-label">
                  <el-icon><Upload /></el-icon>
                  <span>文件上传</span>
                  <el-badge :value="uploads.size" :max="99" v-if="uploads.size > 0" class="tab-badge" />
                </div>
              </template>
              
              <div class="tab-content">
                <!-- 上传拖拽区域 -->
                <div class="upload-zone">
                  <el-upload
                    ref="uploadRef"
                    drag
                    multiple
                    :auto-upload="false"
                    :show-file-list="false"
                    :on-change="handleFileSelect"
                    class="upload-dragger"
                  >
                    <div class="upload-dragger-content">
                      <el-icon class="upload-icon">
                        <UploadFilled />
                      </el-icon>
                      <h3 class="upload-title">拖放文件到这里</h3>
                      <p class="upload-subtitle">或点击选择文件</p>
                      <div class="upload-features">
                        <div class="feature-item">
                          <el-icon><Check /></el-icon>
                          <span>支持多文件</span>
                        </div>
                        <div class="feature-item">
                          <el-icon><Check /></el-icon>
                          <span>断点续传</span>
                        </div>
                        <div class="feature-item">
                          <el-icon><Check /></el-icon>
                          <span>任意格式</span>
                        </div>
                      </div>
                    </div>
                  </el-upload>
                </div>

                <!-- 文件队列列表 -->
                <div v-if="uploads.size > 0" class="queue-container">
                  <SectionHeader
                    :icon="FolderOpened"
                    :title="`上传队列 (${uploads.size})`"
                  />
                  
                  <el-scrollbar class="queue-scrollbar">
                    <div class="file-queue">
                      <div
                        v-for="upload in Array.from(uploads.values())"
                        :key="upload.id"
                        class="file-item"
                        :class="`status-${upload.status}`"
                      >
                        <div class="file-main">
                          <div class="file-icon-wrapper">
                            <el-icon class="file-icon" :size="40"><Document /></el-icon>
                            <div class="status-indicator" :class="upload.status"></div>
                          </div>
                          
                          <div class="file-details">
                            <div class="file-name-row">
                              <h4 class="file-name" :title="upload.file.name">{{ upload.file.name }}</h4>
                              <StatusTag :status="upload.status" size="small" effect="plain" />
                            </div>
                            
                            <div class="file-metadata">
                              <span class="meta-size">{{ formatFileSize(upload.file.size) }}</span>
                              <span class="meta-divider">•</span>
                              <span class="meta-id">ID: {{ upload.id }}</span>
                            </div>

                            <!-- 上传进度 -->
                            <div v-if="upload.status === 'uploading'" class="progress-container">
                              <div class="progress-header">
                                <span class="progress-percentage">{{ upload.progress.toFixed(1) }}%</span>
                                <span class="progress-speed">{{ upload.speed }}</span>
                              </div>
                              <el-progress 
                                :percentage="upload.progress" 
                                :stroke-width="6" 
                                :show-text="false"
                                :color="customColors"
                              />
                            </div>
                          </div>

                          <div class="file-actions">
                            <el-button-group size="small">
                              <el-button v-if="upload.status === 'waiting'" type="primary" @click="startUpload(upload.id)">
                                <el-icon><VideoPlay /></el-icon>
                              </el-button>
                              <el-button v-if="upload.status === 'uploading'" @click="pauseUpload(upload.id)">
                                <el-icon><VideoPause /></el-icon>
                              </el-button>
                              <el-button v-if="upload.status === 'paused'" type="success" @click="resumeUpload(upload.id)">
                                <el-icon><Refresh /></el-icon>
                              </el-button>
                              <el-button v-if="upload.status === 'error'" type="warning" @click="startUpload(upload.id)">
                                <el-icon><RefreshRight /></el-icon>
                              </el-button>
                              <el-button type="danger" @click="removeUpload(upload.id)">
                                <el-icon><Delete /></el-icon>
                              </el-button>
                            </el-button-group>
                          </div>
                        </div>
                      </div>
                    </div>
                  </el-scrollbar>
                </div>

                <!-- 空状态提示 -->
                <EmptyState v-else :icon="FolderOpened" description="队列为空，请添加文件开始上传" />
              </div>
            </el-tab-pane>

            <!-- 上传历史标签页 -->
            <el-tab-pane name="history">
              <template #label>
                <div class="tab-label">
                  <el-icon><Clock /></el-icon>
                  <span>上传记录</span>
                  <el-badge :value="history.length" :max="99" v-if="history.length > 0" class="tab-badge" />
                </div>
              </template>
              
              <div class="tab-content">
                <el-scrollbar v-if="history.length > 0" class="history-scrollbar">
                  <div class="history-list">
                    <div
                      v-for="item in history"
                      :key="item.id"
                      class="history-item"
                    >
                      <div class="history-main">
                        <div class="history-icon-wrapper">
                          <el-icon class="success-icon" :size="40"><CircleCheckFilled /></el-icon>
                        </div>
                        
                        <div class="history-details">
                          <h4 class="history-name" :title="item.fileName">{{ item.fileName }}</h4>
                          <div class="history-metadata">
                            <span class="meta-size">{{ formatFileSize(item.size) }}</span>
                            <span class="meta-divider">•</span>
                            <span class="meta-time">
                              <el-icon><Clock /></el-icon>
                              {{ item.time }}
                            </span>
                          </div>
                        </div>

                        <el-button type="primary" size="small" @click="copyLink(item.url)">
                          <el-icon><CopyDocument /></el-icon>
                          复制链接
                        </el-button>
                      </div>
                    </div>
                  </div>
                </el-scrollbar>
                
                <EmptyState v-else :icon="Clock" description="暂无上传记录" />
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </div>
    </div>
  </PageContainer>
</template>

<script setup>
import { ref } from 'vue'
import {
  VideoPlay,
  VideoPause,
  Refresh,
  Delete,
  Upload,
  UploadFilled,
  Check,
  FolderOpened,
  Clock,
  Document,
  CircleCheckFilled,
  CopyDocument,
  Setting,
  Key,
  RefreshRight,
} from '@element-plus/icons-vue'
import { formatFileSize } from '@/utils/format'
import { useFileUpload } from '@/composables/useFileUpload'
import { useTabs } from '@/composables/useViewMode'
import { PageContainer, EmptyState, InfoItem, SectionHeader, StatusTag } from '@/components/base'

// 标签页管理
const { activeTab } = useTabs('queue')

const uploadRef = ref(null)

// 文件上传管理
const {
  uploads,
  history,
  uploadConfig,
  statusCount,
  totalProgress,
  addFile,
  startUpload,
  pauseUpload,
  resumeUpload,
  removeUpload,
  startAll,
  pauseAll,
  resumeAll,
  clearAll,
  copyLink,
} = useFileUpload({
  taskId: 'task-123456',
  chunkSizeMB: 5,
  parallelUploads: 6,
  queueConcurrency: 3,
})

// 进度条颜色配置
const customColors = [
  { color: '#667eea', percentage: 20 },
  { color: '#764ba2', percentage: 40 },
  { color: '#f093fb', percentage: 60 },
  { color: '#4facfe', percentage: 80 },
  { color: '#00f2fe', percentage: 100 }
]

/**
 * 处理文件选择
 */
const handleFileSelect = (file) => {
  addFile(file.raw)
}
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

  .stats-badges {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);

    .el-badge {
      :deep(.el-tag) {
        font-weight: 500;
      }
    }
  }
}

// ==========================================
// 主要布局
// ==========================================
.content-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: var(--spacing-base);
  height: 100%;
  overflow: hidden;
}

// 左侧配置面板
.config-panel {
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: hidden;

  .config-card {
    height: fit-content;

    :deep(.el-card__body) {
      padding: var(--spacing-lg);
    }

    .el-form {
      .el-form-item {
        margin-bottom: var(--spacing-base);
      }

      .form-help {
        font-size: 12px;
        color: var(--color-text-tertiary);
        line-height: 1.6;
        margin-top: 4px;
      }
    }

    .config-summary {
      .summary-title {
        margin: 0 0 var(--spacing-sm) 0;
        font-size: 14px;
        font-weight: 600;
        color: var(--color-text-primary);
        line-height: 1.5;
      }

      .summary-items {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xs);
      }
    }
  }
}

// 右侧主面板
.main-panel {
  overflow: hidden;
  display: flex;
  flex-direction: column;

  .main-card {
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
  }

  .upload-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    :deep(.el-tabs__header) {
      margin: 0;
      padding: 0 var(--spacing-lg);
      background: var(--color-bg-secondary);
      border-bottom: 1px solid var(--color-border-light);
    }

    :deep(.el-tabs__content) {
      flex: 1;
      overflow: hidden;
      padding: 0;
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
        margin-left: 4px;
      }
    }

    .tab-content {
      height: 100%;
      display: flex;
      flex-direction: column;
      padding: var(--spacing-lg);
      overflow: hidden;
      gap: var(--spacing-base);
    }
  }
}

// ==========================================
// 上传区域样式
// ==========================================
.upload-zone {
  flex-shrink: 0;
}

.upload-dragger {
  :deep(.el-upload) {
    width: 100%;
  }

  :deep(.el-upload-dragger) {
    width: 100%;
    padding: var(--spacing-2xl);
    border: 2px dashed var(--color-border-base);
    border-radius: var(--radius-lg);
    background: linear-gradient(135deg, #f6f8fb 0%, #ffffff 100%);
    transition: all var(--transition-base);

    &:hover {
      border-color: var(--color-primary);
      background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
      box-shadow: 0 0 0 4px rgba(24, 144, 255, 0.08);
    }
  }

  .upload-dragger-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--spacing-base);

    .upload-icon {
      font-size: 64px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .upload-title {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: var(--color-text-primary);
      line-height: 1.5;
    }

    .upload-subtitle {
      margin: 0;
      font-size: 14px;
      color: var(--color-text-secondary);
      line-height: 1.6;
    }

    .upload-features {
      display: flex;
      gap: var(--spacing-lg);
      margin-top: var(--spacing-sm);

      .feature-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: var(--color-text-secondary);

        .el-icon {
          color: var(--color-success);
          font-size: 16px;
        }
      }
    }
  }
}

// ==========================================
// 队列列表样式
// ==========================================
.queue-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-md);
  padding: var(--spacing-base);
  gap: var(--spacing-base);

  .queue-scrollbar {
    flex: 1;
  }

  .file-queue {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }
}

.file-item {
  background: white;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
  border-left: 3px solid var(--color-border-base);

  &.status-waiting {
    border-left-color: var(--color-info);
  }

  &.status-uploading {
    border-left-color: var(--color-success);
    background: linear-gradient(135deg, #f0f9ff 0%, #ffffff 100%);
  }

  &.status-paused {
    border-left-color: var(--color-warning);
  }

  &.status-error {
    border-left-color: var(--color-danger);
    background: linear-gradient(135deg, #fef2f2 0%, #ffffff 100%);
  }

  &:hover {
    box-shadow: var(--shadow-md);
    transform: translateX(4px);
  }

  .file-main {
    padding: var(--spacing-base) var(--spacing-lg);
    display: flex;
    align-items: flex-start;
    gap: var(--spacing-base);
    min-height: 80px;

    .file-icon-wrapper {
      position: relative;
      flex-shrink: 0;

      .file-icon {
        color: var(--color-primary);
        font-size: 40px;
      }

      .status-indicator {
        position: absolute;
        bottom: 0;
        right: 0;
        width: 12px;
        height: 12px;
        border-radius: var(--radius-full);
        border: 2px solid white;

        &.waiting {
          background: var(--color-info);
        }

        &.uploading {
          background: var(--color-success);
          animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
        }

        &.paused {
          background: var(--color-warning);
        }

        &.error {
          background: var(--color-danger);
        }
      }
    }

    .file-details {
      flex: 1;
      min-width: 0;
      display: flex;
      flex-direction: column;
      gap: var(--spacing-xs);

      .file-name-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: var(--spacing-sm);

        .file-name {
          margin: 0;
          font-size: 15px;
          font-weight: 600;
          color: var(--color-text-primary);
          line-height: 1.5;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          flex: 1;
        }
      }

      .file-metadata {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        font-size: 13px;
        color: var(--color-text-secondary);

        .meta-size {
          font-weight: 600;
          color: var(--color-primary);
        }

        .meta-divider {
          color: var(--color-border-dark);
        }

        .meta-id {
          font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        }
      }

      .progress-container {
        padding: var(--spacing-sm);
        background: var(--color-bg-secondary);
        border-radius: var(--radius-md);
        margin-top: var(--spacing-xs);

        .progress-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 6px;

          .progress-percentage {
            font-size: 14px;
            font-weight: 700;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
          }

          .progress-speed {
            font-size: 12px;
            color: var(--color-success);
            font-weight: 600;
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
          }
        }
        
        :deep(.el-progress) {
          .el-progress__text {
            font-size: 12px;
          }
        }
      }
    }

    .file-actions {
      flex-shrink: 0;
    }
  }
}

// ==========================================
// 历史记录样式
// ==========================================
.history-scrollbar {
  height: calc(100vh - 300px);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.history-item {
  background: white;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
  border-left: 3px solid var(--color-success);

  &:hover {
    box-shadow: var(--shadow-md);
    transform: translateX(4px);
  }

  .history-main {
    padding: var(--spacing-base) var(--spacing-lg);
    display: flex;
    align-items: center;
    gap: var(--spacing-base);
    min-height: 72px;

    .history-icon-wrapper {
      flex-shrink: 0;

      .success-icon {
        color: var(--color-success);
        font-size: 40px;
      }
    }

    .history-details {
      flex: 1;
      min-width: 0;

      .history-name {
        margin: 0 0 6px 0;
        font-size: 15px;
        font-weight: 600;
        color: var(--color-text-primary);
        line-height: 1.5;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .history-metadata {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        font-size: 13px;
        color: var(--color-text-secondary);

        .meta-size {
          font-weight: 600;
          color: var(--color-success);
        }

        .meta-divider {
          color: var(--color-border-dark);
        }

        .meta-time {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 1200px) {
  .content-layout {
    grid-template-columns: 240px 1fr;
    gap: var(--spacing-sm);
  }
}

@media (max-width: 768px) {
  .page-toolbar {
    flex-direction: column;
    gap: var(--spacing-sm);
    padding: var(--spacing-sm) var(--spacing-base);
    min-height: auto;

    .toolbar-left {
      width: 100%;
      
      .el-button-group {
        width: 100%;
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: var(--spacing-xs);
        
        .el-button {
          margin: 0 !important;
        }
      }
    }

    .toolbar-right {
      width: 100%;
    }

    .stats-badges {
      flex-wrap: wrap;
      justify-content: center;
      gap: var(--spacing-sm);
    }
  }

  .content-layout {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
    gap: var(--spacing-base);
  }
  
  .config-panel {
    .config-card {
      :deep(.el-card__body) {
        padding: var(--spacing-base);
      }
      
      .el-form {
        .el-form-item {
          margin-bottom: var(--spacing-sm);
        }
      }
      
      .config-summary {
        .summary-title {
          font-size: 13px;
        }
      }
    }
  }

  .upload-dragger {
    :deep(.el-upload-dragger) {
      padding: var(--spacing-xl);
    }

    .upload-dragger-content {
      gap: var(--spacing-sm);
      
      .upload-icon {
        font-size: 48px;
      }

      .upload-title {
        font-size: 16px;
      }
      
      .upload-subtitle {
        font-size: 13px;
      }

      .upload-features {
        flex-direction: column;
        gap: var(--spacing-xs);
        margin-top: var(--spacing-sm);
        
        .feature-item {
          font-size: 12px;
        }
      }
    }
  }

  .file-item {
    .file-main {
      flex-direction: column;
      align-items: stretch;
      padding: var(--spacing-sm) var(--spacing-base);
      min-height: auto;
      
      .file-icon-wrapper {
        .file-icon {
          font-size: 32px;
        }
        
        .status-indicator {
          width: 10px;
          height: 10px;
        }
      }
      
      .file-details {
        .file-name-row {
          .file-name {
            font-size: 14px;
          }
        }
        
        .file-metadata {
          font-size: 12px;
        }
      }

      .file-actions {
        width: 100%;

        .el-button-group {
          width: 100%;
          display: flex;

          .el-button {
            flex: 1;
          }
        }
      }
    }
  }
  
  .history-item {
    .history-main {
      padding: var(--spacing-sm) var(--spacing-base);
      min-height: 64px;
      
      .history-icon-wrapper {
        .success-icon {
          font-size: 32px;
        }
      }
      
      .history-details {
        .history-name {
          font-size: 14px;
          margin-bottom: 4px;
        }
        
        .history-metadata {
          font-size: 12px;
        }
      }
    }
  }
}
</style>