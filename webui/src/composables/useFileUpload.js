/**
 * 文件上传管理 Composable
 */

import { ref, computed } from 'vue'
import * as tus from 'tus-js-client'
import { API_BASE_URL, API_ENDPOINTS } from '@/config'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatFileSize, formatDateTime } from '@/utils/format'

/**
 * TUS 文件上传管理
 * @param {Object} config - 上传配置
 * @returns {Object} 上传管理对象
 */
export function useFileUpload(config = {}) {
  const uploads = ref(new Map())
  const history = ref([])
  const activeUploads = new Set()
  let uploadIdCounter = 0
  let isGlobalPaused = false
  
  const uploadConfig = ref({
    taskId: config.taskId || 'task-123456',
    chunkSizeMB: config.chunkSizeMB || 5,
    parallelUploads: config.parallelUploads || 6,
    queueConcurrency: config.queueConcurrency || 3,
  })
  
  /**
   * 状态统计
   */
  const statusCount = computed(() => {
    const count = { waiting: 0, uploading: 0, paused: 0, error: 0 }
    uploads.value.forEach(upload => {
      count[upload.status]++
    })
    return count
  })
  
  /**
   * 总上传进度
   */
  const totalProgress = computed(() => {
    if (uploads.value.size === 0) return 0
    
    let totalBytes = 0
    let uploadedBytes = 0
    
    uploads.value.forEach(upload => {
      totalBytes += upload.file.size
      uploadedBytes += (upload.file.size * upload.progress) / 100
    })
    
    return totalBytes > 0 ? (uploadedBytes / totalBytes) * 100 : 0
  })
  
  /**
   * 添加文件到队列
   * @param {File} file - 文件对象
   */
  const addFile = (file) => {
    const id = uploadIdCounter++
    uploads.value.set(id, {
      id,
      file,
      upload: null,
      status: 'waiting',
      progress: 0,
      speed: '',
      startTime: null,
    })
    
    // 自动开始处理队列
    processQueue()
  }
  
  /**
   * 创建 TUS 上传实例
   * @param {Object} uploadInfo - 上传信息
   * @returns {tus.Upload} TUS 上传实例
   */
  const createTusUpload = (uploadInfo) => {
    const chunkSize = uploadConfig.value.chunkSizeMB * 1024 * 1024
    
    return new tus.Upload(uploadInfo.file, {
      endpoint: `${API_BASE_URL}${API_ENDPOINTS.files}/`,
      chunkSize,
      addRequestId: true,
      uploadDataDuringCreation: true,
      removeFingerprintOnSuccess: true,
      retryDelays: [0, 1000, 3000, 5000],
      parallelUploads: uploadConfig.value.parallelUploads,
      metadata: {
        filename: uploadInfo.file.name,
        filetype: uploadInfo.file.type,
        task_id: uploadConfig.value.taskId.toString(),
      },
      metadataForPartialUploads: {
        task_id: uploadConfig.value.taskId.toString(),
      },
      onError: (error) => {
        console.error('上传失败:', error)
        uploadInfo.status = 'error'
        activeUploads.delete(uploadInfo.id)
        processQueue()
      },
      onProgress: (bytesUploaded, bytesTotal) => {
        const percentage = (bytesUploaded / bytesTotal) * 100
        uploadInfo.progress = parseFloat(percentage.toFixed(1))
        
        if (uploadInfo.startTime) {
          const elapsed = (Date.now() - uploadInfo.startTime) / 1000
          const speed = bytesUploaded / elapsed
          uploadInfo.speed = formatFileSize(speed) + '/s'
        }
      },
      onSuccess: () => {
        ElMessage.success(`${uploadInfo.file.name} 上传成功`)
        
        // 添加到历史记录
        history.value.unshift({
          id: Date.now(),
          fileName: uploadInfo.file.name,
          size: uploadInfo.file.size,
          time: formatDateTime(new Date()),
          url: uploadInfo.upload.url,
        })
        
        // 从队列中移除
        uploads.value.delete(uploadInfo.id)
        activeUploads.delete(uploadInfo.id)
        processQueue()
      },
    })
  }
  
  /**
   * 开始上传
   * @param {number} id - 上传ID
   */
  const startUpload = async (id) => {
    const uploadInfo = uploads.value.get(id)
    if (!uploadInfo || activeUploads.size >= uploadConfig.value.queueConcurrency) {
      return
    }
    
    activeUploads.add(id)
    uploadInfo.status = 'uploading'
    uploadInfo.startTime = Date.now()
    
    const upload = createTusUpload(uploadInfo)
    uploadInfo.upload = upload
    
    try {
      const previousUploads = await upload.findPreviousUploads()
      if (previousUploads.length > 0) {
        upload.resumeFromPreviousUpload(previousUploads[0])
      }
      upload.start()
    } catch (error) {
      console.error('启动上传失败:', error)
      uploadInfo.status = 'error'
      activeUploads.delete(id)
    }
  }
  
  /**
   * 暂停上传
   * @param {number} id - 上传ID
   */
  const pauseUpload = (id) => {
    const uploadInfo = uploads.value.get(id)
    if (!uploadInfo || !uploadInfo.upload) return
    
    uploadInfo.upload.abort()
    uploadInfo.status = 'paused'
    activeUploads.delete(id)
  }
  
  /**
   * 恢复上传
   * @param {number} id - 上传ID
   */
  const resumeUpload = (id) => {
    const uploadInfo = uploads.value.get(id)
    if (!uploadInfo || !uploadInfo.upload) return
    
    if (activeUploads.size >= uploadConfig.value.queueConcurrency) {
      ElMessage.warning('当前上传数量已达上限')
      return
    }
    
    activeUploads.add(id)
    uploadInfo.status = 'uploading'
    uploadInfo.upload.start()
  }
  
  /**
   * 移除上传
   * @param {number} id - 上传ID
   */
  const removeUpload = (id) => {
    const uploadInfo = uploads.value.get(id)
    if (uploadInfo && uploadInfo.upload) {
      uploadInfo.upload.abort()
    }
    activeUploads.delete(id)
    uploads.value.delete(id)
    processQueue()
  }
  
  /**
   * 处理队列
   */
  const processQueue = () => {
    if (isGlobalPaused) return
    
    const waitingUploads = Array.from(uploads.value.values())
      .filter(u => u.status === 'waiting')
      .slice(0, uploadConfig.value.queueConcurrency - activeUploads.size)
    
    waitingUploads.forEach(upload => {
      startUpload(upload.id)
    })
  }
  
  /**
   * 开始全部上传
   */
  const startAll = () => {
    isGlobalPaused = false
    processQueue()
  }
  
  /**
   * 暂停全部上传
   */
  const pauseAll = () => {
    isGlobalPaused = true
    activeUploads.forEach(id => {
      pauseUpload(id)
    })
  }
  
  /**
   * 恢复全部上传
   */
  const resumeAll = () => {
    isGlobalPaused = false
    const pausedUploads = Array.from(uploads.value.values())
      .filter(u => u.status === 'paused')
    pausedUploads.forEach(upload => {
      resumeUpload(upload.id)
    })
    processQueue()
  }
  
  /**
   * 清空队列
   */
  const clearAll = async () => {
    const confirmed = await ElMessageBox.confirm(
      '确定要清空所有文件吗?',
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    ).catch(() => false)
    
    if (!confirmed) return
    
    activeUploads.forEach(id => {
      const uploadInfo = uploads.value.get(id)
      if (uploadInfo && uploadInfo.upload) {
        uploadInfo.upload.abort()
      }
    })
    uploads.value.clear()
    activeUploads.clear()
    isGlobalPaused = false
  }
  
  /**
   * 复制链接
   * @param {string} url - 文件URL
   */
  const copyLink = (url) => {
    navigator.clipboard.writeText(url).then(() => {
      ElMessage.success('链接已复制到剪贴板')
    }).catch(() => {
      ElMessage.error('复制失败')
    })
  }
  
  /**
   * 更新配置
   * @param {Object} newConfig - 新配置
   */
  const updateConfig = (newConfig) => {
    uploadConfig.value = { ...uploadConfig.value, ...newConfig }
  }
  
  return {
    uploads,
    history,
    uploadConfig,
    statusCount,
    totalProgress,
    addFile,
    startUpload,
    pauseUpload,
    resumeUpload,
    removeUpload,
    startAll,
    pauseAll,
    resumeAll,
    clearAll,
    copyLink,
    updateConfig,
  }
}