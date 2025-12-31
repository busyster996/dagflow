/**
 * 设计令牌 - 统一的设计系统变量
 * 与 styles/index.scss 中的 CSS 变量保持一致
 */

// 颜色系统
export const COLORS = {
  primary: '#1890ff',
  primaryLight: '#40a9ff',
  primaryLighter: '#69c0ff',
  primaryDark: '#096dd9',
  
  success: '#52c41a',
  warning: '#faad14',
  danger: '#ff4d4f',
  info: '#909399',
  
  textPrimary: '#1f2937',
  textSecondary: '#6b7280',
  textTertiary: '#9ca3af',
  textDisabled: '#d1d5db',
  
  bgPrimary: '#ffffff',
  bgSecondary: '#f9fafb',
  bgTertiary: '#f3f4f6',
  bgDark: '#111827',
  
  borderLight: '#f0f0f0',
  borderBase: '#e5e7eb',
  borderDark: '#d1d5db',
}

// 状态颜色映射
export const STATUS_COLORS = {
  running: '#409EFF',
  stopped: '#67C23A',
  failed: '#F56C6C',
  pending: '#E6A23C',
  timeout: '#FF5722',
  canceled: '#7C4DFF',
  skipped: '#909399',
  unknown: '#909399',
  paused: '#E6A23C',
  killed: '#F56C6C',
  success: '#67C23A',
}

// 状态类型映射（Element Plus）
export const STATUS_TYPES = {
  running: 'primary',
  stopped: 'success',
  failed: 'danger',
  pending: 'warning',
  timeout: 'danger',
  canceled: 'info',
  skipped: 'info',
  unknown: 'info',
  paused: 'warning',
  killed: 'danger',
  success: 'success',
}

// 状态文本映射
export const STATUS_TEXTS = {
  running: '运行中',
  stopped: '已停止',
  failed: '已失败',
  pending: '待执行',
  timeout: '超时',
  canceled: '已取消',
  skipped: '已跳过',
  unknown: '未知',
  paused: '已暂停',
  killed: '已终止',
  success: '成功',
}

// 间距系统 - 舒适的视觉呼吸感
export const SPACING = {
  xs: '8px',
  sm: '12px',
  base: '16px',
  md: '20px',
  lg: '24px',
  xl: '32px',
  '2xl': '40px',
  '3xl': '48px',
  '4xl': '64px',
}

// 圆角系统 - 现代友好设计
export const RADIUS = {
  sm: '4px',
  base: '6px',
  md: '8px',
  lg: '12px',
  xl: '16px',
  full: '9999px',
}

// 阴影系统 - 清晰的层次感
export const SHADOWS = {
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  base: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
}

// 过渡时间 - 流畅的动画体验
export const TRANSITIONS = {
  fast: '150ms',
  base: '200ms',
  slow: '300ms',
}

// Z-index层级
export const Z_INDEX = {
  dropdown: 1000,
  sticky: 1010,
  fixed: 1020,
  modalBackdrop: 1030,
  modal: 1040,
  popover: 1050,
  tooltip: 1060,
}

// 断点系统
export const BREAKPOINTS = {
  xs: 480,
  sm: 768,
  md: 1024,
  lg: 1280,
  xl: 1536,
  '2xl': 1920,
}

// 分页默认配置 - 增加每页显示数量
export const PAGINATION = {
  defaultPage: 1,
  defaultSize: 20,
  pageSizes: [20, 40, 60, 80],
}