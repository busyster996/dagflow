/**
 * API配置管理
 * 运行时动态获取API地址，与 internal/server/router/static 实现保持一致
 * 
 * 生产环境：使用window.APP_CONFIG（由index.html注入）
 * 开发环境：使用Vite代理
 */

/**
 * 获取生产环境配置
 * @returns {Object} API配置对象
 */
const getProductionConfig = () => {
  if (!window.APP_CONFIG) {
    console.error('❌ APP_CONFIG未初始化！请检查index.html中的配置脚本');
    throw new Error('APP_CONFIG未初始化');
  }
  
  const { baseUrl, wsBaseUrl, api } = window.APP_CONFIG;
  
  return {
    baseUrl,
    wsBaseUrl,
    apiEndpoints: {
      task: api.task,
      pipeline: api.pipeline,
      event: api.event
    }
  };
};

/**
 * 获取开发环境配置
 * @returns {Object} API配置对象
 */
const getDevelopmentConfig = () => {
  const protocol = window.location.protocol;
  const hostname = window.location.hostname;
  const port = 3000; // Vite开发服务器端口
  
  return {
    baseUrl: '', // 使用相对路径，由Vite代理处理
    wsBaseUrl: `${protocol === 'https:' ? 'wss' : 'ws'}://${hostname}:${port}`,
    apiEndpoints: {
      task: '/api/v1/task',
      pipeline: '/api/v1/pipeline',
      event: '/api/v1/event'
    }
  };
};

/**
 * 导出API配置
 * 根据环境自动选择配置
 */
export const apiConfig = import.meta.env.DEV 
  ? getDevelopmentConfig() 
  : getProductionConfig();

/**
 * 获取完整的API URL
 * @param {string} endpoint - API端点路径
 * @returns {string} 完整的API URL
 */
export const getApiUrl = (endpoint) => {
  return `${apiConfig.baseUrl}${endpoint}`;
};

/**
 * 获取完整的WebSocket URL
 * @param {string} endpoint - WebSocket端点路径
 * @returns {string} 完整的WebSocket URL
 */
export const getWsUrl = (endpoint) => {
  return `${apiConfig.wsBaseUrl}${endpoint}`;
};

/**
 * 获取SSE URL（用于事件流）
 * @param {string} endpoint - SSE端点路径
 * @returns {string} 完整的SSE URL
 */
export const getSseUrl = (endpoint) => {
  return `${apiConfig.baseUrl}${endpoint}`;
};

/**
 * 导出API端点常量
 */
export const API_ENDPOINTS = {
  // 任务相关
  TASK_LIST: apiConfig.apiEndpoints.task,
  TASK_DETAIL: (name) => `${apiConfig.apiEndpoints.task}/${name}`,
  TASK_DUMP: (name) => `${apiConfig.apiEndpoints.task}/${name}/dump`,
  TASK_STEP: (name) => `${apiConfig.apiEndpoints.task}/${name}/step`,
  TASK_STEP_DETAIL: (taskName, stepName) => `${apiConfig.apiEndpoints.task}/${taskName}/step/${stepName}`,
  TASK_STEP_LOG: (taskName, stepName) => `${apiConfig.apiEndpoints.task}/${taskName}/step/${stepName}/log`,
  
  // 流水线相关
  PIPELINE_LIST: apiConfig.apiEndpoints.pipeline,
  PIPELINE_DETAIL: (name) => `${apiConfig.apiEndpoints.pipeline}/${name}`,
  PIPELINE_BUILD: (name) => `${apiConfig.apiEndpoints.pipeline}/${name}/build`,
  
  // 事件流
  EVENT_STREAM: apiConfig.apiEndpoints.event
};

/**
 * 导出WebSocket端点常量
 */
export const WS_ENDPOINTS = {
  TASK_LIST: getWsUrl(apiConfig.apiEndpoints.task),
  TASK_STEP: (taskName) => getWsUrl(`${apiConfig.apiEndpoints.task}/${taskName}/step`),
  TASK_STEP_LOG: (taskName, stepName) => getWsUrl(`${apiConfig.apiEndpoints.task}/${taskName}/step/${stepName}/log`),
  PIPELINE_LIST: getWsUrl(apiConfig.apiEndpoints.pipeline),
  PIPELINE_BUILD: (pipelineName) => getWsUrl(`${apiConfig.apiEndpoints.pipeline}/${pipelineName}/build`)
};

// 开发环境下打印配置信息
if (import.meta.env.DEV) {
  console.log('📡 API配置已加载:', {
    mode: '开发环境',
    baseUrl: apiConfig.baseUrl || '(使用代理)',
    wsBaseUrl: apiConfig.wsBaseUrl,
    endpoints: apiConfig.apiEndpoints
  });
}

export default apiConfig;