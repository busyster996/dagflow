import { WS_BASE_URL } from '@/config'

export class WebSocketManager {
  constructor(url, onMessage, onError = null, options = {}) {
    this.url = `${WS_BASE_URL}${url}`
    this.onMessage = onMessage
    this.onError = onError
    this.reconnectInterval = options.reconnectInterval || 5000
    this.maxReconnectAttempts = options.maxReconnectAttempts || 5
    this.socket = null
    this.isManuallyClosed = false
    this.reconnectAttempts = 0
    this.connect()
  }

  connect() {
    if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
      return
    }

    this.isManuallyClosed = false
    this.socket = new WebSocket(this.url)

    this.socket.onopen = () => {
      this.reconnectAttempts = 0
    }

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        this.onMessage(data)
      } catch (error) {
        console.error('解析WebSocket消息失败:', error)
      }
    }

    this.socket.onerror = (error) => {
      console.error('WebSocket错误:', error)
      if (this.onError) {
        this.onError(error)
      }
    }

    this.socket.onclose = (event) => {
      if (event.wasClean) {
      } else {
        console.log(`WebSocket意外关闭: ${event.code}, ${event.reason}`)
        this.handleReconnect()
      }
    }
  }

  handleReconnect() {
    if (this.isManuallyClosed || this.reconnectAttempts >= this.maxReconnectAttempts) {
      return
    }

    this.reconnectAttempts++
    const delay = Math.min(this.reconnectInterval * Math.pow(2, this.reconnectAttempts - 1), 30000)

    console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)
    setTimeout(() => {
      this.connect()
    }, delay)
  }

  send(data) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(data))
    } else {
      console.warn('WebSocket未连接，无法发送数据')
    }
  }

  close() {
    this.isManuallyClosed = true
    if (this.socket) {
      this.socket.close(1000, '正常关闭')
      this.socket = null
    }
  }
}