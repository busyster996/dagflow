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
        :icon="DocumentAdd"
        title="创建新任务"
        :show-close="false"
      >
        <template #actions>
          <el-button size="default" @click="handleClose">
            <el-icon><Close /></el-icon>
            取消
          </el-button>
          <el-button type="primary" size="default" @click="handleCreate" :loading="loading">
            <el-icon><Check /></el-icon>
            创建任务
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
import { ref, watch } from 'vue'
import { DocumentAdd, Close, Check, Edit } from '@element-plus/icons-vue'
import { createTask } from '@/api/task'
import { TASK_YAML_TEMPLATE } from '@/config'
import { ElMessage } from 'element-plus'
import { DialogHeader, MonacoEditor } from '@/components/base'

const props = defineProps({
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = ref(false)
const editorRef = ref(null)
const loading = ref(false)
const yamlContent = ref(`# 任务名称, 可选, 默认自动生成\nname: 测试任务\n${TASK_YAML_TEMPLATE}`)

watch(
  () => props.modelValue,
  (val) => {
    dialogVisible.value = val
  }
)

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
  if (!val) {
    handleClose()
  }
})

/**
 * 创建任务
 */
const handleCreate = async () => {
  const content = yamlContent.value.trim()
  if (!content) {
    ElMessage.warning('请输入配置内容')
    return
  }

  loading.value = true
  try {
    const result = await createTask(content)
    ElMessage.success('任务创建成功')
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
  // 重置内容
  yamlContent.value = `# 任务名称, 可选, 默认自动生成\nname: 测试任务\n${TASK_YAML_TEMPLATE}`
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