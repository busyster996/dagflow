import { request, requestYaml } from '@/utils/request'
import { API_ENDPOINTS } from '@/config'
import { ElMessage, ElMessageBox } from 'element-plus'

// 获取流水线详情
export function getPipelineDetail(pipelineName) {
  return request(`${API_ENDPOINTS.pipeline}/${pipelineName}`)
}

// 创建流水线
export function createPipeline(yamlContent) {
  return requestYaml(API_ENDPOINTS.pipeline, yamlContent, 'POST')
}

// 更新流水线
export function updatePipeline(pipelineName, yamlContent) {
  return requestYaml(`${API_ENDPOINTS.pipeline}/${pipelineName}`, yamlContent, 'POST')
}

// 删除流水线
export async function deletePipeline(pipelineName) {
  const confirmed = await ElMessageBox.confirm(
    `确定要删除流水线 "${pipelineName}"?`,
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
    await request(`${API_ENDPOINTS.pipeline}/${pipelineName}`, {
      method: 'DELETE',
    })
    ElMessage.success('流水线删除成功')
  } catch (error) {
    throw error
  }
}

// 运行流水线
export function runPipeline(pipelineName, paramsYaml) {
  return requestYaml(
    `${API_ENDPOINTS.pipeline}/${pipelineName}/build`,
    paramsYaml,
    'POST'
  )
}