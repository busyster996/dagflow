/**
 * 对话框状态管理 Composable
 */

import { ref, watch } from 'vue'

/**
 * 对话框状态管理
 * @param {Object} options - 配置选项
 * @param {Function} options.onOpen - 打开时的回调
 * @param {Function} options.onClose - 关闭时的回调
 * @returns {Object} 对话框管理对象
 */
export function useDialog(options = {}) {
  const { onOpen, onClose } = options
  
  const visible = ref(false)
  const loading = ref(false)
  
  /**
   * 打开对话框
   */
  const open = () => {
    visible.value = true
    if (onOpen) {
      onOpen()
    }
  }
  
  /**
   * 关闭对话框
   */
  const close = () => {
    visible.value = false
    loading.value = false
    if (onClose) {
      onClose()
    }
  }
  
  /**
   * 切换对话框状态
   */
  const toggle = () => {
    if (visible.value) {
      close()
    } else {
      open()
    }
  }
  
  /**
   * 开始加载
   */
  const startLoading = () => {
    loading.value = true
  }
  
  /**
   * 结束加载
   */
  const stopLoading = () => {
    loading.value = false
  }
  
  return {
    visible,
    loading,
    open,
    close,
    toggle,
    startLoading,
    stopLoading,
  }
}

/**
 * 多对话框管理
 * @returns {Object} 多对话框管理对象
 */
export function useMultipleDialogs() {
  const dialogs = ref(new Map())
  
  /**
   * 创建对话框
   * @param {string} name - 对话框名称
   * @param {Object} options - 配置选项
   */
  const createDialog = (name, options = {}) => {
    if (!dialogs.value.has(name)) {
      dialogs.value.set(name, useDialog(options))
    }
    return dialogs.value.get(name)
  }
  
  /**
   * 获取对话框
   * @param {string} name - 对话框名称
   */
  const getDialog = (name) => {
    return dialogs.value.get(name)
  }
  
  /**
   * 打开对话框
   * @param {string} name - 对话框名称
   */
  const openDialog = (name) => {
    const dialog = getDialog(name)
    if (dialog) {
      dialog.open()
    }
  }
  
  /**
   * 关闭对话框
   * @param {string} name - 对话框名称
   */
  const closeDialog = (name) => {
    const dialog = getDialog(name)
    if (dialog) {
      dialog.close()
    }
  }
  
  /**
   * 关闭所有对话框
   */
  const closeAllDialogs = () => {
    dialogs.value.forEach(dialog => {
      dialog.close()
    })
  }
  
  return {
    dialogs,
    createDialog,
    getDialog,
    openDialog,
    closeDialog,
    closeAllDialogs,
  }
}