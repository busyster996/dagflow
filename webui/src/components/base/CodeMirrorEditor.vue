<template>
  <div class="codemirror-editor-wrapper">
    <div v-if="showToolbar" class="editor-toolbar">
      <div class="toolbar-left">
        <el-icon class="toolbar-icon"><component :is="toolbarIcon" /></el-icon>
        <span class="toolbar-title">{{ title }}</span>
      </div>
      <div class="toolbar-right">
        <el-tag size="small" :type="languageTagType">{{ languageLabel }}</el-tag>
        <slot name="toolbar-actions" />
      </div>
    </div>
    <div ref="containerRef" class="codemirror-editor-container" :style="containerStyle"></div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useCodeMirror } from '@/composables/useCodeMirror'
import { Edit } from '@element-plus/icons-vue'

const props = defineProps({
  // v-model 绑定值
  modelValue: {
    type: String,
    default: '',
  },
  // 编辑器语言
  language: {
    type: String,
    default: 'yaml',
  },
  // 编辑器主题 ('light' 或 'dark')
  theme: {
    type: String,
    default: 'dark',
  },
  // 是否只读
  readOnly: {
    type: Boolean,
    default: false,
  },
  // 编辑器高度
  height: {
    type: [String, Number],
    default: '100%',
  },
  // 显示工具栏
  showToolbar: {
    type: Boolean,
    default: true,
  },
  // 工具栏标题
  title: {
    type: String,
    default: '代码编辑器',
  },
  // 工具栏图标
  toolbarIcon: {
    type: Object,
    default: () => Edit,
  },
  // 额外的编辑器选项
  options: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:modelValue', 'change', 'ready'])

const containerRef = ref(null)

const {
  initEditor,
  getValue,
  setValue,
  onContentChange,
  dispose,
  getEditor,
  setReadOnly,
  setLanguage,
  setTheme,
  focus,
  formatDocument,
} = useCodeMirror({
  language: props.language,
  theme: props.theme,
  readOnly: props.readOnly,
  initialValue: props.modelValue,
  editorOptions: props.options,
})

/**
 * 容器样式
 */
const containerStyle = computed(() => {
  const height = typeof props.height === 'number' ? `${props.height}px` : props.height
  return {
    height,
  }
})

/**
 * 语言标签类型
 */
const languageTagType = computed(() => {
  const typeMap = {
    yaml: 'info',
    json: 'success',
    javascript: 'warning',
    typescript: 'primary',
  }
  return typeMap[props.language] || 'info'
})

/**
 * 语言标签文本
 */
const languageLabel = computed(() => {
  return props.language.toUpperCase()
})

/**
 * 监听外部值变化
 */
watch(
  () => props.modelValue,
  (newValue) => {
    const currentValue = getValue()
    if (newValue !== currentValue) {
      setValue(newValue)
    }
  }
)

/**
 * 监听语言变化
 */
watch(
  () => props.language,
  (newLang) => {
    setLanguage(newLang)
  }
)

/**
 * 监听主题变化
 */
watch(
  () => props.theme,
  (newTheme) => {
    setTheme(newTheme)
  }
)

/**
 * 监听只读状态变化
 */
watch(
  () => props.readOnly,
  (newReadOnly) => {
    setReadOnly(newReadOnly)
  }
)

/**
 * 初始化
 */
onMounted(async () => {
  if (containerRef.value) {
    await initEditor(containerRef.value, props.modelValue)
    
    // 监听内容变化
    onContentChange((newValue) => {
      emit('update:modelValue', newValue)
      emit('change', newValue)
    })
    
    emit('ready', getEditor())
  }
})

/**
 * 清理
 */
onUnmounted(() => {
  dispose()
})

/**
 * 暴露方法供父组件调用
 */
defineExpose({
  getValue,
  setValue,
  getEditor,
  focus: () => focus(),
  format: () => formatDocument(),
})
</script>

<style lang="scss" scoped>
.codemirror-editor-wrapper {
  display: flex;
  flex-direction: column;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-sm);
  background: white;
  
  .editor-toolbar {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--spacing-sm) var(--spacing-base);
    background: linear-gradient(135deg, #f6f8fb 0%, #fafbfc 100%);
    border-bottom: 1px solid var(--color-border-light);
    min-height: 48px;
    
    .toolbar-left {
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);
      
      .toolbar-icon {
        font-size: 16px;
        color: var(--color-primary);
      }
      
      .toolbar-title {
        font-size: 14px;
        font-weight: 600;
        color: var(--color-text-primary);
        line-height: 1.5;
      }
    }
    
    .toolbar-right {
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);
    }
  }
  
  .codemirror-editor-container {
    flex: 1;
    width: 100%;
    min-height: 0;
    position: relative;
    
    // CodeMirror 6 样式覆盖
    :deep(.cm-editor) {
      height: 100%;
      font-size: 14px;
      
      .cm-scroller {
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        line-height: 1.7;
        padding: var(--spacing-base);
      }
      
      .cm-gutters {
        border-right: 1px solid var(--color-border-light);
      }
      
      .cm-activeLineGutter {
        background-color: rgba(103, 126, 234, 0.05);
      }
      
      .cm-activeLine {
        background-color: rgba(103, 126, 234, 0.03);
      }
      
      // 光标样式
      .cm-cursor {
        border-left-width: 2px;
      }
      
      // 选中文本样式
      .cm-selectionBackground {
        background-color: rgba(103, 126, 234, 0.2) !important;
      }
      
      &.cm-focused {
        outline: none;
        
        .cm-selectionBackground {
          background-color: rgba(103, 126, 234, 0.3) !important;
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .codemirror-editor-wrapper {
    .editor-toolbar {
      padding: var(--spacing-xs) var(--spacing-sm);
      min-height: 44px;
      
      .toolbar-left {
        gap: var(--spacing-xs);
        
        .toolbar-icon {
          font-size: 14px;
        }
        
        .toolbar-title {
          font-size: 13px;
        }
      }
    }
    
    .codemirror-editor-container {
      :deep(.cm-editor) {
        font-size: 13px;
        
        .cm-scroller {
          padding: var(--spacing-sm);
        }
      }
    }
  }
}
</style>