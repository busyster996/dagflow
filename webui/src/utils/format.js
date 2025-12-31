/**
 * 格式化工具函数
 */

/**
 * 格式化文件大小
 * @param {number} bytes - 字节数
 * @returns {string} 格式化后的文件大小
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

/**
 * 格式化时间
 * @param {string|number} time - 时间字符串或时间戳
 * @returns {string} 格式化后的时间
 */
export function formatTime(time) {
  if (!time) return '---'
  
  // 如果是时间字符串且包含冒号，提取时分秒
  if (typeof time === 'string' && time.includes(':')) {
    const match = time.match(/(\d{2}):(\d{2}):(\d{2})/)
    if (match) {
      return `${match[1]}:${match[2]}:${match[3]}`
    }
    return time
  }
  
  return time
}

/**
 * 格式化日期时间
 * @param {Date|string|number} date - 日期对象、字符串或时间戳
 * @param {string} format - 格式化模板，默认 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的日期时间
 */
export function formatDateTime(date, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!date) return '---'
  
  const d = date instanceof Date ? date : new Date(date)
  
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化相对时间
 * @param {Date|string|number} date - 日期对象、字符串或时间戳
 * @returns {string} 相对时间描述
 */
export function formatRelativeTime(date) {
  if (!date) return '---'
  
  const d = date instanceof Date ? date : new Date(date)
  const now = new Date()
  const diff = now - d
  
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  if (days > 7) return formatDateTime(d, 'YYYY-MM-DD')
  if (days > 0) return `${days}天前`
  if (hours > 0) return `${hours}小时前`
  if (minutes > 0) return `${minutes}分钟前`
  if (seconds > 30) return `${seconds}秒前`
  return '刚刚'
}

/**
 * 格式化持续时间
 * @param {number} milliseconds - 毫秒数
 * @returns {string} 格式化后的持续时间
 */
export function formatDuration(milliseconds) {
  if (!milliseconds || milliseconds < 0) return '0秒'
  
  const seconds = Math.floor(milliseconds / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  const parts = []
  
  if (days > 0) parts.push(`${days}天`)
  if (hours % 24 > 0) parts.push(`${hours % 24}小时`)
  if (minutes % 60 > 0) parts.push(`${minutes % 60}分钟`)
  if (seconds % 60 > 0 && days === 0) parts.push(`${seconds % 60}秒`)
  
  return parts.join(' ') || '0秒'
}

/**
 * 截断文本
 * @param {string} text - 要截断的文本
 * @param {number} maxLength - 最大长度
 * @param {string} suffix - 后缀，默认 '...'
 * @returns {string} 截断后的文本
 */
export function truncateText(text, maxLength, suffix = '...') {
  if (!text || text.length <= maxLength) return text
  return text.substring(0, maxLength - suffix.length) + suffix
}

/**
 * 高亮搜索关键词
 * @param {string} text - 原始文本
 * @param {string} keyword - 搜索关键词
 * @param {string} className - 高亮样式类名
 * @returns {string} 带高亮标记的 HTML 字符串
 */
export function highlightKeyword(text, keyword, className = 'highlight') {
  if (!text || !keyword) return text
  
  const regex = new RegExp(`(${keyword})`, 'gi')
  return text.replace(regex, `<span class="${className}">$1</span>`)
}

/**
 * 数字千分位格式化
 * @param {number} num - 数字
 * @param {number} precision - 小数位数
 * @returns {string} 格式化后的数字字符串
 */
export function formatNumber(num, precision = 0) {
  if (num === null || num === undefined) return '0'
  
  const fixed = Number(num).toFixed(precision)
  const parts = fixed.split('.')
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  
  return parts.join('.')
}

/**
 * 百分比格式化
 * @param {number} value - 数值（0-1 或 0-100）
 * @param {number} precision - 小数位数
 * @param {boolean} isDecimal - 是否为小数格式（0-1）
 * @returns {string} 格式化后的百分比字符串
 */
export function formatPercent(value, precision = 1, isDecimal = false) {
  if (value === null || value === undefined) return '0%'
  
  const percent = isDecimal ? value * 100 : value
  return `${percent.toFixed(precision)}%`
}