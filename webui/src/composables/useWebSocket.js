/**
 * WebSocket 管理 Composable
 */

import { ref, onUnmounted } from 'vue'
import { WebSocketManager } from '@/utils/websocket'

/**
 * WebSocket 连接管理
 * @param {string} url - WebSocket URL
 * @param {Object} options - 配置选项
 * @param {Function} options.onMessage - 消息处理回调
 * @param {Function} options.onError - 错误处理回调
 * @param {Function} options.onOpen - 连接打开回调
 * @param {Function} options.onClose - 连接关闭回调
 * @param {number} options.reconnectInterval - 重连间隔
 * @param {number} options.maxReconnectAttempts - 最大重连次数
 * @returns {Object} WebSocket 管理对象
 */
export function useWebSocket(url, options = {}) {
  const {
    onMessage,
    onError,
    onOpen,
    onClose,
    reconnectInterval,
    maxReconnectAttempts,
  } = options
  
  const isConnected = ref(false)
  const isConnecting = ref(false)
  const error = ref(null)
  const data = ref(null)
  
  let wsManager = null
  
  /**
   * 连接WebSocket
   */
  const connect = () => {
    if (wsManager) {
      wsManager.close()
    }
    
    isConnecting.value = true
    error.value = null
    
    const handleMessage = (message) => {
      data.value = message
      if (onMessage) {
        onMessage(message)
      }
    }
    
    const handleError = (err) => {
      error.value = err
      isConnected.value = false
      isConnecting.value = false
      if (onError) {
        onError(err)
      }
    }
    
    wsManager = new WebSocketManager(
      url,
      handleMessage,
      handleError,
      {
        reconnectInterval,
        maxReconnectAttempts,
      }
    )
    
    // 监听连接状态
    const checkConnection = setInterval(() => {
      if (wsManager && wsManager.socket) {
        if (wsManager.socket.readyState === WebSocket.OPEN) {
          isConnected.value = true
          isConnecting.value = false
          if (onOpen && !isConnected.value) {
            onOpen()
          }
        } else if (wsManager.socket.readyState === WebSocket.CLOSED) {
          const wasConnected = isConnected.value
          isConnected.value = false
          isConnecting.value = false
          if (onClose && wasConnected) {
            onClose()
          }
        }
      }
    }, 100)
    
    // 清理定时器
    onUnmounted(() => {
      clearInterval(checkConnection)
    })
  }
  
  /**
   * 发送消息
   * @param {any} message - 要发送的消息
   */
  const send = (message) => {
    if (wsManager) {
      wsManager.send(message)
    }
  }
  
  /**
   * 关闭连接
   */
  const close = () => {
    if (wsManager) {
      wsManager.close()
      wsManager = null
      isConnected.value = false
      isConnecting.value = false
    }
  }
  
  /**
   * 重新连接
   */
  const reconnect = () => {
    close()
    connect()
  }
  
  /**
   * 获取连接状态
   * @returns {number} WebSocket 连接状态
   */
  const getReadyState = () => {
    if (!wsManager || !wsManager.socket) {
      return WebSocket.CLOSED
    }
    return wsManager.socket.readyState
  }
  
  // 组件卸载时自动关闭连接
  onUnmounted(() => {
    close()
  })
  
  return {
    isConnected,
    isConnecting,
    error,
    data,
    connect,
    send,
    close,
    reconnect,
    getReadyState,
  }
}

/**
 * 用于列表数据的WebSocket管理
 * @param {string} url - WebSocket URL
 * @param {Object} options - 配置选项
 * @returns {Object} 列表数据管理对象
 */
export function useWebSocketList(url, options = {}) {
  const items = ref([])
  const loading = ref(false)
  const pagination = ref({
    current: 1,
    size: 15,
    total: 1,
  })
  
  const handleMessage = (response) => {
    if (response.data) {
      // 处理列表数据
      if (Array.isArray(response.data)) {
        items.value = response.data
      } else if (response.data.list || response.data.items) {
        items.value = response.data.list || response.data.items
      } else if (response.data.tasks) {
        items.value = response.data.tasks
      } else if (response.data.pipelines) {
        items.value = response.data.pipelines
      }
      
      // 处理分页信息
      if (response.data.page) {
        pagination.value = {
          current: response.data.page.current || 1,
          size: response.data.page.size || 15,
          total: response.data.page.total || 1,
        }
      }
    } else {
      items.value = []
    }
    
    loading.value = false
    
    if (options.onMessage) {
      options.onMessage(response)
    }
  }
  
  const ws = useWebSocket(url, {
    ...options,
    onMessage: handleMessage,
  })
  
  /**
   * 刷新数据
   * @param {Object} params - 查询参数
   */
  const refresh = (params = {}) => {
    loading.value = true
    const queryParams = {
      page: pagination.value.current,
      size: pagination.value.size,
      ...params,
    }
    ws.send(queryParams)
  }
  
  /**
   * 改变页码
   * @param {number} page - 新页码
   */
  const changePage = (page) => {
    pagination.value.current = page
    refresh()
  }
  
  /**
   * 改变每页数量
   * @param {number} size - 新的每页数量
   */
  const changePageSize = (size) => {
    pagination.value.size = size
    pagination.value.current = 1
    refresh()
  }
  
  return {
    ...ws,
    items,
    loading,
    pagination,
    refresh,
    changePage,
    changePageSize,
  }
}