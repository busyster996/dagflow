import { request, requestYaml } from '@/utils/request'
import { API_ENDPOINTS } from '@/config'
import { ElMessage, ElMessageBox } from 'element-plus'

// 获取任务列表（通过WebSocket）
export function getTaskList() {
  // 这个通过WebSocket实现，在组件中处理
}

// 获取任务详情
export function getTaskDetail(taskName) {
  return request(`${API_ENDPOINTS.task}/${taskName}`)
}

// 创建任务
export function createTask(yamlContent) {
  return requestYaml(API_ENDPOINTS.task, yamlContent, 'POST')
}

// 任务操作（kill等）
export async function taskAction(taskName, action) {
  const confirmed = await ElMessageBox.confirm(
    `确定要${action} "${taskName}"?`,
    '确认操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)

  if (!confirmed) return

  try {
    const result = await request(`${API_ENDPOINTS.task}/${taskName}?action=${action}`, {
      method: 'PUT',
    })
    ElMessage.success(`${action}操作成功`)
    return result
  } catch (error) {
    throw error
  }
}

// 删除任务
export async function deleteTask(taskName) {
  const confirmed = await ElMessageBox.confirm(
    `确定要删除任务 "${taskName}"?`,
    '确认删除',
    {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
      confirmButtonClass: 'el-button--danger',
    }
  ).catch(() => false)

  if (!confirmed) return

  try {
    await request(`${API_ENDPOINTS.task}/${taskName}`, {
      method: 'DELETE',
    })
    ElMessage.success('任务删除成功')
  } catch (error) {
    throw error
  }
}

// 导出任务配置
export async function dumpTask(taskName) {
  try {
    const response = await request(`${API_ENDPOINTS.task}/${taskName}/dump`)
    const blob = new Blob([response.data], { type: 'application/yaml' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${taskName}.yaml`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('任务配置导出成功')
  } catch (error) {
    throw error
  }
}

// 步骤操作
export async function stepAction(taskName, stepName, action) {
  const confirmed = await ElMessageBox.confirm(
    `确定要${action} "${stepName}"?`,
    '确认操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)

  if (!confirmed) return

  try {
    const result = await request(
      `${API_ENDPOINTS.task}/${taskName}/step/${stepName}?action=${action}`,
      {
        method: 'PUT',
      }
    )
    ElMessage.success(`${action}操作成功`)
    return result
  } catch (error) {
    throw error
  }
}

// 获取步骤详情
export function getStepDetail(taskName, stepName) {
  return request(`${API_ENDPOINTS.task}/${taskName}/step/${stepName}`)
}