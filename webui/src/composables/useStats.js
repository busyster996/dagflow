/**
 * 统计数据管理 Composable
 */

import { ref, computed } from 'vue'

/**
 * 通用统计数据管理
 * @param {Array} items - 数据项数组
 * @param {Function} groupBy - 分组函数
 * @returns {Object} 统计管理对象
 */
export function useStats(items, groupBy) {
  const stats = computed(() => {
    if (!items || !items.value || !Array.isArray(items.value)) {
      return {}
    }
    
    const result = {}
    
    items.value.forEach(item => {
      const key = groupBy(item)
      result[key] = (result[key] || 0) + 1
    })
    
    return result
  })
  
  return {
    stats,
  }
}

/**
 * 任务状态统计
 * @param {Array} tasks - 任务列表
 * @returns {Object} 任务统计对象
 */
export function useTaskStats(tasks) {
  const stats = computed(() => {
    const result = {
      total: 0,
      running: 0,
      stopped: 0,
      failed: 0,
      pending: 0,
      timeout: 0,
      canceled: 0,
      skipped: 0,
      paused: 0,
      killed: 0,
    }
    
    if (!tasks || !tasks.value || !Array.isArray(tasks.value)) {
      return result
    }
    
    result.total = tasks.value.length
    
    tasks.value.forEach(task => {
      const state = task.state || 'unknown'
      if (result[state] !== undefined) {
        result[state]++
      }
    })
    
    return result
  })
  
  /**
   * 获取成功率
   */
  const successRate = computed(() => {
    if (stats.value.total === 0) return 0
    return ((stats.value.stopped / stats.value.total) * 100).toFixed(1)
  })
  
  /**
   * 获取失败率
   */
  const failureRate = computed(() => {
    if (stats.value.total === 0) return 0
    const failed = stats.value.failed + stats.value.timeout + stats.value.killed
    return ((failed / stats.value.total) * 100).toFixed(1)
  })
  
  /**
   * 获取活跃任务数
   */
  const activeCount = computed(() => {
    return stats.value.running + stats.value.pending + stats.value.paused
  })
  
  /**
   * 获取完成任务数
   */
  const completedCount = computed(() => {
    return stats.value.stopped + stats.value.failed + stats.value.timeout + 
           stats.value.canceled + stats.value.skipped + stats.value.killed
  })
  
  return {
    stats,
    successRate,
    failureRate,
    activeCount,
    completedCount,
  }
}

/**
 * 流水线状态统计
 * @param {Array} pipelines - 流水线列表
 * @returns {Object} 流水线统计对象
 */
export function usePipelineStats(pipelines) {
  const stats = computed(() => {
    const result = {
      total: 0,
      enabled: 0,
      disabled: 0,
      jinja: 0,
      yaml: 0,
    }
    
    if (!pipelines || !pipelines.value || !Array.isArray(pipelines.value)) {
      return result
    }
    
    result.total = pipelines.value.length
    
    pipelines.value.forEach(pipeline => {
      // 统计启用/禁用状态
      if (pipeline.disable) {
        result.disabled++
      } else {
        result.enabled++
      }
      
      // 统计模板类型
      if (pipeline.tplType === 'jinja2') {
        result.jinja++
      } else {
        result.yaml++
      }
    })
    
    return result
  })
  
  /**
   * 获取启用率
   */
  const enabledRate = computed(() => {
    if (stats.value.total === 0) return 0
    return ((stats.value.enabled / stats.value.total) * 100).toFixed(1)
  })
  
  return {
    stats,
    enabledRate,
  }
}

/**
 * 步骤状态统计
 * @param {Array} steps - 步骤列表
 * @returns {Object} 步骤统计对象
 */
export function useStepStats(steps) {
  const stats = computed(() => {
    const result = {
      total: 0,
      running: 0,
      stopped: 0,
      failed: 0,
      pending: 0,
      paused: 0,
      skipped: 0,
    }
    
    if (!steps || !steps.value || !Array.isArray(steps.value)) {
      return result
    }
    
    result.total = steps.value.length
    
    steps.value.forEach(step => {
      const state = step.state || 'unknown'
      if (result[state] !== undefined) {
        result[state]++
      }
    })
    
    return result
  })
  
  /**
   * 获取进度百分比
   */
  const progressPercent = computed(() => {
    if (stats.value.total === 0) return 0
    const completed = stats.value.stopped + stats.value.failed + stats.value.skipped
    return ((completed / stats.value.total) * 100).toFixed(1)
  })
  
  /**
   * 是否全部完成
   */
  const isAllCompleted = computed(() => {
    return stats.value.running === 0 && stats.value.pending === 0 && stats.value.paused === 0
  })
  
  /**
   * 是否有失败步骤
   */
  const hasFailures = computed(() => {
    return stats.value.failed > 0
  })
  
  return {
    stats,
    progressPercent,
    isAllCompleted,
    hasFailures,
  }
}

/**
 * 文件上传统计
 * @param {Map} uploads - 上传文件Map
 * @returns {Object} 上传统计对象
 */
export function useUploadStats(uploads) {
  const stats = computed(() => {
    const result = {
      total: 0,
      waiting: 0,
      uploading: 0,
      paused: 0,
      error: 0,
      completed: 0,
    }
    
    if (!uploads || !uploads.value) {
      return result
    }
    
    result.total = uploads.value.size
    
    uploads.value.forEach(upload => {
      const status = upload.status || 'waiting'
      if (result[status] !== undefined) {
        result[status]++
      }
    })
    
    return result
  })
  
  /**
   * 总上传进度
   */
  const totalProgress = computed(() => {
    if (!uploads || !uploads.value || uploads.value.size === 0) {
      return 0
    }
    
    let totalBytes = 0
    let uploadedBytes = 0
    
    uploads.value.forEach(upload => {
      totalBytes += upload.file.size
      uploadedBytes += (upload.file.size * upload.progress) / 100
    })
    
    return totalBytes > 0 ? (uploadedBytes / totalBytes) * 100 : 0
  })
  
  /**
   * 是否全部完成
   */
  const isAllCompleted = computed(() => {
    return stats.value.total > 0 && 
           stats.value.waiting === 0 && 
           stats.value.uploading === 0 && 
           stats.value.paused === 0
  })
  
  /**
   * 是否有错误
   */
  const hasErrors = computed(() => {
    return stats.value.error > 0
  })
  
  return {
    stats,
    totalProgress,
    isAllCompleted,
    hasErrors,
  }
}