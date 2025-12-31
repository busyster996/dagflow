<template>
  <div v-loading="loading" class="step-content">
    <!-- 步骤基本信息 -->
    <el-card shadow="hover" class="section-card">
      <template #header>
        <SectionHeader
          :icon="InfoFilled"
          title="基本信息"
        />
      </template>
      <div class="info-grid" v-if="stepData">
        <InfoItem
          label="步骤名称"
          :value="stepData.name"
          hoverable
        />
        <InfoItem label="步骤类型" hoverable>
          <template #value>
            <el-tag type="info" size="small">{{ stepData.type }}</el-tag>
          </template>
        </InfoItem>
        <InfoItem label="步骤状态" hoverable>
          <template #value>
            <StatusTag :status="stepData.state" effect="plain" size="small" show-icon />
          </template>
        </InfoItem>
        <InfoItem
          v-if="stepData.code !== undefined"
          label="状态码"
          :value="stepData.code"
          :value-class="stepData.code === 0 ? 'success' : 'error'"
          hoverable
        />
        <InfoItem
          label="开始时间"
          :value="stepData.time?.start || '---'"
          value-class="time"
          hoverable
        />
        <InfoItem
          label="结束时间"
          :value="stepData.time?.end || '---'"
          value-class="time"
          hoverable
        />
      </div>
    </el-card>

    <!-- 环境变量 -->
    <el-card 
      v-if="stepData?.env && stepData.env.length > 0" 
      shadow="hover" 
      class="section-card collapsible-card"
    >
      <template #header>
        <SectionHeader
          :icon="Key"
          title="环境变量"
          :count="stepData.env.length"
          clickable
          :expandable="true"
          :is-expanded="expandedSections.env"
          @click="toggleSection('env')"
        />
      </template>
      <el-collapse-transition>
        <div v-show="expandedSections.env">
          <el-scrollbar max-height="300px">
            <div class="env-list">
              <div v-for="(env, index) in stepData.env" :key="index" class="env-item">
                <div class="env-header">
                  <el-icon class="env-icon"><Setting /></el-icon>
                  <span class="env-name">{{ env.name }}</span>
                </div>
                <div class="env-value-wrapper">
                  <code class="env-value">{{ env.value }}</code>
                </div>
              </div>
            </div>
          </el-scrollbar>
        </div>
      </el-collapse-transition>
    </el-card>

    <!-- 输入内容 -->
    <el-card shadow="hover" class="section-card collapsible-card">
      <template #header>
        <SectionHeader
          :icon="Document"
          title="输入内容"
          clickable
          :expandable="true"
          :is-expanded="expandedSections.input"
          @click="toggleSection('input')"
        />
      </template>
      <el-collapse-transition>
        <div v-show="expandedSections.input">
          <el-scrollbar max-height="300px">
            <pre class="code-content">{{ stepData?.content || '无内容' }}</pre>
          </el-scrollbar>
        </div>
      </el-collapse-transition>
    </el-card>

    <!-- 执行输出 -->
    <el-card shadow="hover" class="section-card output-card">
      <template #header>
        <SectionHeader
          :icon="Monitor"
          title="执行输出"
          tag="实时更新"
          tag-type="success"
        />
      </template>
      <el-scrollbar ref="outputScrollbar" class="output-scrollbar">
        <pre ref="outputRef" class="output-content">{{ output || '暂无输出' }}</pre>
      </el-scrollbar>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onUnmounted, nextTick } from 'vue'
import { InfoFilled, Key, Setting, Document, Monitor } from '@element-plus/icons-vue'
import { getStepDetail } from '@/api/task'
import { useWebSocketList } from '@/composables/useWebSocket'
import { API_ENDPOINTS } from '@/config'
import { SectionHeader, InfoItem, StatusTag, EmptyState } from '@/components/base'

const props = defineProps({
  taskName: String,
  stepName: String,
})

const loading = ref(false)
const stepData = ref(null)
const output = ref('')
const outputScrollbar = ref(null)
const outputRef = ref(null)
const expandedSections = reactive({
  env: false,
  input: false,
})

// WebSocket 管理日志输出
const {
  connect: connectLog,
  close: closeLog,
} = useWebSocketList(
  `${API_ENDPOINTS.task}/${props.taskName}/step/${props.stepName}/log`,
  {
    onMessage: (response) => {
      if (response.data && response.data.length > 0) {
        const newData = response.data.map(item => item.content).join('\n')
        output.value += newData + '\n'

        // 自动滚动到底部
        nextTick(() => {
          if (outputScrollbar.value) {
            const scrollEl = outputScrollbar.value.$refs.wrap$
            if (scrollEl) {
              scrollEl.scrollTop = scrollEl.scrollHeight
            }
          }
        })
      }
    },
    onError: () => {
      output.value = stepData.value?.message || '无输出'
    }
  }
)

/**
 * 切换折叠区域
 */
const toggleSection = (section) => {
  expandedSections[section] = !expandedSections[section]
}

/**
 * 加载步骤详情
 */
const loadStepDetail = async () => {
  loading.value = true
  try {
    const response = await getStepDetail(props.taskName, props.stepName)
    stepData.value = response.data
  } catch (error) {
    console.error('获取步骤详情失败:', error)
  } finally {
    loading.value = false
  }
}

watch(
  () => props.stepName,
  () => {
    if (props.stepName) {
      loadStepDetail()
      connectLog()
    }
  },
  { immediate: true }
)

onUnmounted(() => {
  closeLog()
})
</script>

<style lang="scss" scoped>
.step-content {
  padding: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-base);
  height: 100%;
  overflow-y: auto;
}

// ==========================================
// 卡片样式
// ==========================================
.section-card {
  flex-shrink: 0;
}

// 信息网格
.info-grid {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

// 环境变量列表
.env-list {
  padding: var(--spacing-sm);

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

    .env-header {
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);
      margin-bottom: var(--spacing-sm);

      .env-icon {
        font-size: 14px;
        color: var(--color-primary);
      }

      .env-name {
        font-weight: 600;
        color: var(--color-text-primary);
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        font-size: 13px;
        line-height: 1.5;
      }
    }

    .env-value-wrapper {
      padding-left: 22px;

      .env-value {
        display: block;
        color: var(--color-text-secondary);
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        font-size: 12px;
        word-break: break-all;
        background: var(--color-bg-primary);
        padding: var(--spacing-sm);
        border-radius: var(--radius-md);
        line-height: 1.6;
      }
    }
  }
}

// 代码内容
.code-content {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  margin: 0;
  padding: var(--spacing-base);
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
  border-radius: var(--radius-md);
  white-space: pre-wrap;
  word-break: break-word;
}

// 输出卡片
.output-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 400px;

  :deep(.el-card__body) {
    flex: 1;
    padding: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .output-scrollbar {
    flex: 1;
  }

  .output-content {
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.6;
    margin: 0;
    padding: var(--spacing-base);
    background: #1e1e1e;
    color: #d4d4d4;
    min-height: 100%;
    white-space: pre-wrap;
    word-break: break-word;
  }
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 768px) {
  .step-content {
    padding: var(--spacing-base);
    gap: var(--spacing-sm);
  }
  
  .env-list {
    padding: var(--spacing-xs);
    
    .env-item {
      padding: var(--spacing-xs) var(--spacing-sm);
      margin-bottom: var(--spacing-xs);
      
      .env-header {
        gap: var(--spacing-xs);
        margin-bottom: var(--spacing-xs);
        
        .env-icon {
          font-size: 12px;
        }
        
        .env-name {
          font-size: 12px;
        }
      }
      
      .env-value-wrapper {
        padding-left: 18px;
        
        .env-value {
          font-size: 11px;
          padding: var(--spacing-xs);
        }
      }
    }
  }
  
  .code-content {
    font-size: 12px;
    padding: var(--spacing-sm);
  }

  .output-card {
    min-height: 300px;
  }
  
  .output-content {
    font-size: 12px;
    line-height: 1.5;
    padding: var(--spacing-sm);
  }
}
</style>