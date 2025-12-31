/**
 * 状态管理工具函数
 */

import { STATUS_COLORS, STATUS_TYPES, STATUS_TEXTS } from '@/constants/design-tokens'

/**
 * 获取状态颜色
 * @param {string} status - 状态值
 * @returns {string} 颜色值
 */
export function getStatusColor(status) {
  return STATUS_COLORS[status] || STATUS_COLORS.unknown
}

/**
 * 获取状态类型（Element Plus）
 * @param {string} status - 状态值
 * @returns {string} Element Plus 类型
 */
export function getStatusType(status) {
  return STATUS_TYPES[status] || STATUS_TYPES.unknown
}

/**
 * 获取状态文本
 * @param {string} status - 状态值
 * @returns {string} 状态文本
 */
export function getStatusText(status) {
  return STATUS_TEXTS[status] || status
}

/**
 * 判断是否为运行中状态
 * @param {string} status - 状态值
 * @returns {boolean} 是否为运行中
 */
export function isRunningStatus(status) {
  return status === 'running'
}

/**
 * 判断是否为完成状态
 * @param {string} status - 状态值
 * @returns {boolean} 是否已完成
 */
export function isCompletedStatus(status) {
  return ['stopped', 'success', 'failed', 'timeout', 'canceled', 'skipped', 'killed'].includes(status)
}

/**
 * 判断是否为成功状态
 * @param {string} status - 状态值
 * @returns {boolean} 是否成功
 */
export function isSuccessStatus(status) {
  return ['stopped', 'success'].includes(status)
}

/**
 * 判断是否为失败状态
 * @param {string} status - 状态值
 * @returns {boolean} 是否失败
 */
export function isFailedStatus(status) {
  return ['failed', 'timeout', 'killed'].includes(status)
}

/**
 * 判断是否可以执行操作
 * @param {string} status - 状态值
 * @param {string} action - 操作类型（kill, pause, resume）
 * @returns {boolean} 是否可以执行
 */
export function canPerformAction(status, action) {
  const actionMap = {
    kill: ['running', 'paused', 'pending'],
    pause: ['running', 'pending'],
    resume: ['paused'],
    restart: ['stopped', 'failed', 'timeout', 'canceled', 'killed'],
  }
  
  return actionMap[action]?.includes(status) || false
}

/**
 * 获取状态优先级（用于排序）
 * @param {string} status - 状态值
 * @returns {number} 优先级数值（越小越优先）
 */
export function getStatusPriority(status) {
  const priorityMap = {
    running: 1,
    pending: 2,
    paused: 3,
    failed: 4,
    timeout: 5,
    killed: 6,
    canceled: 7,
    stopped: 8,
    success: 9,
    skipped: 10,
    unknown: 11,
  }
  
  return priorityMap[status] || 999
}

/**
 * 获取状态码样式类名
 * @param {number} code - 状态码
 * @returns {string} 样式类名
 */
export function getCodeClass(code) {
  if (code === 0) return 'success'
  if (code > 0) return 'error'
  return 'default'
}

/**
 * 判断状态码是否为成功
 * @param {number} code - 状态码
 * @returns {boolean} 是否成功
 */
export function isSuccessCode(code) {
  return code === 0
}