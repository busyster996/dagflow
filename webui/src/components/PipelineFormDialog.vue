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
        :icon="isEdit ? Edit : DocumentAdd"
        :title="isEdit ? `编辑流水线: ${pipelineData?.name}` : '创建新流水线'"
        :show-close="false"
      >
        <template #actions>
          <el-button size="default" @click="handleClose">
            <el-icon><Close /></el-icon>
            取消
          </el-button>
          <el-button type="primary" size="default" @click="handleSubmit" :loading="loading">
            <el-icon><Check /></el-icon>
            {{ isEdit ? '保存更改' : '创建流水线' }}
          </el-button>
        </template>
      </DialogHeader>
    </template>

    <div class="dialog-body">
      <!-- Monaco 编辑器 -->
      <MonacoEditor
        ref="editorRef"
        v-model="yamlContent"
        language="yaml"
        theme="vs-dark"
        title="YAML 配置编辑器"
        :toolbar-icon="Edit"
        height="100%"
      />
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { DocumentAdd, Edit, Close, Check } from '@element-plus/icons-vue'
import { createPipeline, updatePipeline } from '@/api/pipeline'
import { TASK_YAML_TEMPLATE } from '@/config'
import { ElMessage } from 'element-plus'
import { DialogHeader, MonacoEditor } from '@/components/base'

const props = defineProps({
  modelValue: Boolean,
  pipelineData: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = ref(false)
const editorRef = ref(null)
const loading = ref(false)
const yamlContent = ref('')

/**
 * 是否为编辑模式
 */
const isEdit = computed(() => !!props.pipelineData)

/**
 * 构建初始内容
 */
const buildInitialContent = () => {
  if (isEdit.value && props.pipelineData) {
    // 编辑模式：使用现有流水线数据
    const desc = props.pipelineData.desc || ''
    const tplType = props.pipelineData.tplType || 'jinja2'
    const content = props.pipelineData.content || ''
    
    return `# 流水线描述\ndesc: |-\n  ${desc.replace(/\n/g, '\n  ')}\n\n# 模板类型\ntplType: ${tplType}\n\n# 流水线配置内容\ncontent: |-\n  ${content.replace(/\n/g, '\n  ')}`
  } else {
    // 创建模式：使用模板
    return `# 流水线名称, 必填\nname: 新流水线\n\n# 流水线描述, 可选\ndesc: |-\n  流水线描述信息\n\n# 模板类型, 默认 jinja2\ntplType: jinja2\n\n# 流水线配置内容\ncontent: |-\n  ${TASK_YAML_TEMPLATE.replace(/\n/g, '\n  ')}`
  }
}

watch(
  () => props.modelValue,
  (val) => {
    dialogVisible.value = val
    if (val) {
      // 每次打开时重新构建内容
      yamlContent.value = buildInitialContent()
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
 * 提交表单
 */
const handleSubmit = async () => {
  const content = yamlContent.value.trim()
  if (!content) {
    ElMessage.warning('请输入配置内容')
    return
  }

  loading.value = true
  try {
    if (isEdit.value && props.pipelineData) {
      // 编辑模式
      await updatePipeline(props.pipelineData.name, content)
      ElMessage.success('流水线更新成功')
    } else {
      // 创建模式
      await createPipeline(content)
      ElMessage.success('流水线创建成功')
    }
    emit('success')
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

  .editor-hint {
    flex-shrink: 0;

    .el-alert {
      border-radius: var(--radius-md);
      padding: var(--spacing-sm) var(--spacing-base);
      font-size: 14px;
      line-height: 1.6;
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
    
    .editor-hint {
      .el-alert {
        font-size: 13px;
      }
    }
  }
}
</style>