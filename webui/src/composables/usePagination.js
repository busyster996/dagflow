/**
 * 分页管理 Composable
 */

import { ref, computed } from 'vue'
import { PAGINATION } from '@/constants/design-tokens'

/**
 * 分页状态管理
 * @param {Object} options - 配置选项
 * @param {number} options.initialPage - 初始页码
 * @param {number} options.initialSize - 初始每页数量
 * @param {Array<number>} options.pageSizes - 可选每页数量
 * @param {Function} options.onPageChange - 页码变化回调
 * @param {Function} options.onSizeChange - 每页数量变化回调
 * @returns {Object} 分页管理对象
 */
export function usePagination(options = {}) {
  const {
    initialPage = PAGINATION.defaultPage,
    initialSize = PAGINATION.defaultSize,
    pageSizes = PAGINATION.pageSizes,
    onPageChange,
    onSizeChange,
  } = options
  
  const currentPage = ref(initialPage)
  const pageSize = ref(initialSize)
  const total = ref(0)
  
  /**
   * 总页数
   */
  const totalPages = computed(() => {
    return Math.ceil(total.value / pageSize.value)
  })
  
  /**
   * 是否有上一页
   */
  const hasPrev = computed(() => {
    return currentPage.value > 1
  })
  
  /**
   * 是否有下一页
   */
  const hasNext = computed(() => {
    return currentPage.value < totalPages.value
  })
  
  /**
   * 当前页的数据范围
   */
  const dataRange = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value + 1
    const end = Math.min(currentPage.value * pageSize.value, total.value)
    return { start, end }
  })
  
  /**
   * 设置总数
   * @param {number} newTotal - 新的总数
   */
  const setTotal = (newTotal) => {
    total.value = newTotal
  }
  
  /**
   * 跳转到指定页
   * @param {number} page - 页码
   */
  const goToPage = (page) => {
    if (page < 1 || page > totalPages.value) return
    
    currentPage.value = page
    if (onPageChange) {
      onPageChange(page)
    }
  }
  
  /**
   * 上一页
   */
  const prevPage = () => {
    if (hasPrev.value) {
      goToPage(currentPage.value - 1)
    }
  }
  
  /**
   * 下一页
   */
  const nextPage = () => {
    if (hasNext.value) {
      goToPage(currentPage.value + 1)
    }
  }
  
  /**
   * 跳转到首页
   */
  const firstPage = () => {
    goToPage(1)
  }
  
  /**
   * 跳转到末页
   */
  const lastPage = () => {
    goToPage(totalPages.value)
  }
  
  /**
   * 改变每页数量
   * @param {number} size - 新的每页数量
   */
  const changePageSize = (size) => {
    pageSize.value = size
    currentPage.value = 1 // 重置到第一页
    if (onSizeChange) {
      onSizeChange(size)
    }
  }
  
  /**
   * 重置分页
   */
  const reset = () => {
    currentPage.value = initialPage
    pageSize.value = initialSize
    total.value = 0
  }
  
  /**
   * 根据后端响应更新分页信息
   * @param {Object} pageInfo - 分页信息
   */
  const updateFromResponse = (pageInfo) => {
    if (!pageInfo) return
    
    if (pageInfo.current !== undefined) {
      currentPage.value = pageInfo.current
    }
    if (pageInfo.size !== undefined) {
      pageSize.value = pageInfo.size
    }
    if (pageInfo.total !== undefined) {
      // 如果后端返回的是总页数，需要转换为总记录数
      total.value = pageInfo.total * pageSize.value
    }
  }
  
  return {
    currentPage,
    pageSize,
    total,
    pageSizes,
    totalPages,
    hasPrev,
    hasNext,
    dataRange,
    setTotal,
    goToPage,
    prevPage,
    nextPage,
    firstPage,
    lastPage,
    changePageSize,
    reset,
    updateFromResponse,
  }
}