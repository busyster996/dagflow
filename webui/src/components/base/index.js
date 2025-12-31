/**
 * 基础组件统一导出
 */

export { default as StatusTag } from './StatusTag.vue'
export { default as StatCard } from './StatCard.vue'
export { default as PageHeader } from './PageHeader.vue'
export { default as EmptyState } from './EmptyState.vue'
export { default as DialogHeader } from './DialogHeader.vue'
export { default as InfoItem } from './InfoItem.vue'
export { default as CodeMirrorEditor } from './CodeMirrorEditor.vue'
export { default as SectionHeader } from './SectionHeader.vue'
export { default as CardGrid } from './CardGrid.vue'
export { default as PageContainer } from './PageContainer.vue'

// 向后兼容别名（推荐使用 CodeMirrorEditor）
export { default as MonacoEditor } from './CodeMirrorEditor.vue'