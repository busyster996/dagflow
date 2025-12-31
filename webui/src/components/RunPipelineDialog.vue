<template>
  <el-dialog
    v-model="dialogVisible"
    :close-on-click-modal="false"
    :show-close="false"
    @close="handleClose"
    fullscreen
    class="form-dialog"
  >
    <template #header>
      <DialogHeader
        :icon="VideoPlay"
        :title="`运行流水线: ${pipelineName}`"
        :show-close="false"
      >
        <template #actions>
          <el-button size="default" @click="handleClose">
            <el-icon><Close /></el-icon>
            取消
          </el-button>
          <el-button type="primary" size="default" @click="handleRun" :loading="loading">
            <el-icon><VideoPlay /></el-icon>
            开始执行
          </el-button>
        </template>
      </DialogHeader>
    </template>

    <div class="dialog-body">
      <!-- 参数说明 -->
      <el-alert
        title="参数配置说明"
        type="info"
        :closable="false"
        show-icon
        class="param-hint"
      >
        <template #default>
          请使用 YAML 格式配置流水线执行参数,这些参数将在模板渲染时使用。
        </template>
      </el-alert>

      <!-- 编辑器 -->
      <MonacoEditor
        ref="editorRef"
        v-model="paramsContent"
        language="yaml"
        theme="vs-dark"
        title="参数配置编辑器"
        :toolbar-icon="Setting"
        height="100%"
      />
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { VideoPlay, Close, Setting } from '@element-plus/icons-vue'
import { runPipeline } from '@/api/pipeline'
import { PIPELINE_PARAMS_TEMPLATE } from '@/config'
import { ElMessage } from 'element-plus'
import { DialogHeader, MonacoEditor } from '@/components/base'

const props = defineProps({
  modelValue: Boolean,
  pipelineName: String,
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = ref(false)
const editorRef = ref(null)
const loading = ref(false)
const paramsContent = ref(PIPELINE_PARAMS_TEMPLATE)

watch(
  () => props.modelValue,
  (val) => {
    dialogVisible.value = val
    if (val) {
      // 重置为默认模板
      paramsContent.value = PIPELINE_PARAMS_TEMPLATE
    }
  }
)

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
})

/**
 * 运行流水线
 */
const handleRun = async () => {
  const content = paramsContent.value.trim() || 'params: {}'

  loading.value = true
  try {
    const result = await runPipeline(props.pipelineName, content)
    ElMessage.success('流水线执行成功')
    emit('success', result.data.name)
    handleClose()
  } catch (error) {
    // 错误已在API层处理
  } finally {
    loading.value = false
  }
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  dialogVisible.value = false
  paramsContent.value = PIPELINE_PARAMS_TEMPLATE
}
</script>

<style lang="scss" scoped>
.form-dialog {
  :deep(.el-dialog__header) {
    padding: 0;
    margin: 0;
  }

  :deep(.el-dialog__body) {
    padding: 0;
    height: calc(100vh - 68px);
    display: flex;
    flex-direction: column;
  }
}

.dialog-body {
  flex: 1;
  padding: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-base);
  overflow: hidden;
  
  .param-hint {
    :deep(.el-alert__content) {
      padding: var(--spacing-sm) 0;
    }
    
    :deep(.el-alert__description) {
      font-size: 14px;
      line-height: 1.6;
      margin-top: 4px;
    }
  }

  :deep(.codemirror-editor-wrapper) {
    flex: 1;
    min-height: 0;
  }
}

@media (max-width: 768px) {
  .form-dialog {
    :deep(.el-dialog__body) {
      height: calc(100vh - 80px);
    }
  }
  
  .dialog-body {
    padding: var(--spacing-base);
    
    .param-hint {
      :deep(.el-alert__description) {
        font-size: 13px;
      }
    }
  }
}
</style>