# DagFlow WebUI

DagFlow 的现代化 Web 前端界面，基于 Vue 3 + Vite 构建，提供直观的任务和流水线管理功能。

## ✨ 特性

- 🚀 **现代技术栈** - Vue 3 + Vite + Element Plus
- 📊 **流程可视化** - 基于 Vue Flow 的 DAG 图形编辑器
- 💻 **代码编辑** - 集成 CodeMirror 6，支持 YAML 语法高亮
- 📁 **文件上传** - 基于 TUS 协议的可断点续传文件上传
- 🎨 **响应式设计** - 完美支持桌面和移动端
- ⚡ **开发体验** - HMR、自动导入、TypeScript 支持
- 🔄 **实时更新** - WebSocket 实时任务状态推送
- 🌐 **灵活部署** - 支持一次构建，多环境部署

## 📦 技术栈

### 核心框架
- [Vue 3.5](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Vite 6](https://vitejs.dev/) - 下一代前端构建工具
- [Vue Router 4](https://router.vuejs.org/) - 官方路由管理器
- [Pinia 2](https://pinia.vuejs.org/) - Vue 状态管理库

### UI 组件库
- [Element Plus 2.8](https://element-plus.org/) - Vue 3 组件库
- [@element-plus/icons-vue](https://element-plus.org/zh-CN/component/icon.html) - Element Plus 图标库

### 可视化与编辑器
- [@vue-flow/core](https://vueflow.dev/) - Vue 流程图组件
- [CodeMirror 6](https://codemirror.net/) - 代码编辑器
- [@codemirror/lang-yaml](https://github.com/codemirror/lang-yaml) - YAML 语法支持

### 其他依赖
- [tus-js-client](https://github.com/tus/tus-js-client) - 可断点续传文件上传客户端
- [Sass](https://sass-lang.com/) - CSS 预处理器

### 开发工具
- [unplugin-auto-import](https://github.com/antfu/unplugin-auto-import) - API 自动导入
- [unplugin-vue-components](https://github.com/antfu/unplugin-vue-components) - 组件自动导入

## 🚀 快速开始

### 环境要求

- Node.js >= 18.0.0
- npm >= 8.0.0

### 安装依赖

```bash
npm install
```

### 开发模式

启动开发服务器（支持热重载）：

```bash
npm run dev
```

开发服务器将在 http://localhost:3000 启动，并自动代理 API 请求到 `http://localhost:2376`。

### 生产构建

```bash
npm run build
```

构建产物将输出到 `../internal/server/router/static` 目录，可直接被后端服务使用。

### 预览构建

```bash
npm run preview
```

## 📁 项目结构

```
webui/
├── src/
│   ├── api/                    # API 接口定义
│   │   ├── pipeline.js         # 流水线相关 API
│   │   └── task.js             # 任务相关 API
│   ├── components/             # Vue 组件
│   │   ├── base/               # 基础组件
│   │   │   ├── CardGrid.vue    # 卡片网格
│   │   │   ├── CodeMirrorEditor.vue  # 代码编辑器
│   │   │   ├── DialogHeader.vue      # 对话框头部
│   │   │   ├── EmptyState.vue        # 空状态
│   │   │   ├── InfoItem.vue          # 信息项
│   │   │   ├── PageContainer.vue     # 页面容器
│   │   │   ├── PageHeader.vue        # 页面头部
│   │   │   ├── SectionHeader.vue     # 章节头部
│   │   │   ├── StatCard.vue          # 统计卡片
│   │   │   └── StatusTag.vue         # 状态标签
│   │   ├── CustomNode.vue           # 自定义流程节点
│   │   ├── PipelineDetailDialog.vue # 流水线详情
│   │   ├── PipelineFormDialog.vue   # 流水线表单
│   │   ├── RunPipelineDialog.vue    # 运行流水线
│   │   ├── StepContent.vue          # 步骤内容
│   │   ├── TaskDetailDialog.vue     # 任务详情
│   │   ├── TaskFormDialog.vue       # 任务表单
│   │   └── VueFlowGraph.vue         # 流程图组件
│   ├── composables/            # 组合式函数
│   │   ├── useCodeMirror.js    # CodeMirror 钩子
│   │   ├── useDialog.js        # 对话框钩子
│   │   ├── useFileUpload.js    # 文件上传钩子
│   │   ├── usePagination.js    # 分页钩子
│   │   ├── useStats.js         # 统计钩子
│   │   ├── useViewMode.js      # 视图模式钩子
│   │   └── useWebSocket.js     # WebSocket 钩子
│   ├── App.vue                 # 根组件
│   ├── main.js                 # 应用入口
│   └── auto-imports.d.ts       # 自动导入类型声明
├── index.html                  # HTML 模板
├── vite.config.js              # Vite 配置
├── package.json                # 项目依赖
└── .env                        # 环境变量

```

## ⚙️ 配置说明

### 环境变量

开发环境配置文件 `.env`：

```bash
# API Base URL（开发环境）
VITE_API_BASE_URL=http://localhost:2376

# WebSocket Base URL（开发环境）
VITE_WS_BASE_URL=ws://localhost:2376
```

**注意**：这些环境变量仅在开发环境中使用。生产环境会使用运行时动态配置。

### 动态配置机制

生产环境采用运行时动态配置，无需重新构建即可适配不同部署环境：

- ✅ 一次构建，到处部署
- ✅ 自动适配不同域名/端口
- ✅ 支持 HTTP/HTTPS 自动切换
- ✅ 支持子路径部署

配置逻辑见 `index.html` 中的 `APP_CONFIG` 对象。

### Vite 配置

`vite.config.js` 主要配置项：

```javascript
{
  server: {
    port: 3000,                    // 开发服务器端口
    proxy: {
      '/api': {
        target: 'http://localhost:2376',  // API 代理目标
        changeOrigin: true,
        ws: true                   // WebSocket 代理
      }
    }
  },
  build: {
    outDir: '../internal/server/router/static',  // 构建输出目录
    chunkSizeWarningLimit: 1500,  // 代码块大小警告限制
    rollupOptions: {
      output: {
        manualChunks: {
          'element-plus': ['element-plus'],
          'codemirror': ['codemirror', '@codemirror/lang-yaml', ...],
          'vue-flow': ['@vue-flow/core', '@vue-flow/background', ...]
        }
      }
    }
  }
}
```

## 🎯 核心功能

### 任务管理
- 创建、编辑、删除任务
- 查看任务详情和执行日志
- 实时任务状态更新
- 任务步骤可视化

### 流水线管理
- 可视化流水线编辑器
- 拖拽式节点编排
- 流水线执行和监控
- 参数化构建支持

### 文件上传
- 支持断点续传（基于 TUS 协议）
- 大文件分片上传
- 上传进度实时显示
- 文件工作区管理

### 代码编辑
- YAML 语法高亮
- 代码折叠和格式化
- 深色主题支持
- 实时预览

## 🔧 开发指南

### 代码规范

项目使用 ESLint + Prettier 进行代码规范检查：

- 使用组合式 API（Composition API）
- 采用 `<script setup>` 语法
- 组件命名采用 PascalCase
- 文件命名采用 kebab-case

### 自动导入

项目配置了自动导入功能，无需手动导入常用 API：

```javascript
// ❌ 不需要手动导入
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

// ✅ 直接使用
const count = ref(0)
const router = useRouter()
```

### 组件开发

基础组件位于 `src/components/base/` 目录，已全局注册可直接使用：

```vue
<template>
  <PageContainer>
    <PageHeader title="页面标题" />
    <StatCard label="统计" :value="100" />
  </PageContainer>
</template>
```

## 🚢 部署指南

### 构建生产版本

```bash
npm run build
```

构建产物会自动输出到后端静态文件目录 `../internal/server/router/static`。

### 部署注意事项

1. **单页应用路由**：确保服务器配置了 SPA 路由回退
2. **API 代理**：生产环境无需代理，使用相对路径 `/api`
3. **WebSocket 支持**：确保服务器支持 WebSocket 连接
4. **静态资源**：构建产物包含代码分割，确保正确配置 MIME 类型

### Docker 部署

DagFlow 后端已包含前端静态资源，直接部署后端即可：

```bash
docker run -p 2376:2376 dagflow/dagflow:latest
```

访问 http://localhost:2376 即可使用 Web UI。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 提交规范

提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建工具或辅助工具的变动

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](../LICENSE) 文件

## 🔗 相关链接

- [DagFlow 主项目](https://github.com/busyster996/dagflow)
- [在线文档](https://github.com/busyster996/dagflow/docs)
- [问题反馈](https://github.com/busyster996/dagflow/issues)

## 📮 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 [Issue](https://github.com/busyster996/dagflow/issues)
- 发送邮件至项目维护者

---

**Made with ❤️ by DagFlow Team**