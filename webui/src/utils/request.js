import { ElMessage } from 'element-plus'
import { API_BASE_URL } from '@/config'

// 通用请求函数
export async function request(url, options = {}) {
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    console.error('API请求失败:', error)
    ElMessage.error(`请求失败: ${error.message}`)
    throw error
  }
}

// YAML内容请求
export async function requestYaml(url, data, method = 'POST') {
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      method,
      headers: {
        'Content-Type': 'application/yaml',
      },
      body: data,
    })

    const result = await response.json()

    if (!response.ok || result.code !== 0) {
      throw new Error(result.message || '请求失败')
    }

    return result
  } catch (error) {
    console.error('YAML请求失败:', error)
    ElMessage.error(`请求失败: ${error.message}`)
    throw error
  }
}