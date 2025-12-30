// ==================== 配置管理 ====================
class ConfigManager {
    static get API_ENDPOINTS() {
        return {
            task: taskUrl,
            pipeline: pipelineUrl,
            event: eventUrl
        };
    }

    static get WS_ENDPOINTS() {
        return {
            task: `${wsBaseUrl}${taskUrl}`,
            pipeline: `${wsBaseUrl}${pipelineUrl}`
        };
    }

    static get TEMPLATES() {
        return {
            taskYaml: `# 类型. 空: 串行无策略, strategy: 串行支持策略, dag: 自定义编排
kind: dag
# 描述
desc: 这是一段任务描述
# 延后执行
#delayed: 2025-04-27T15:04:05Z
# 允许节点, 可选, 默认为当前节点
#node: node-01
# 禁用, 可选, 默认false
#disable: false
# 超时时间, 可选, 默认48小时
timeout: 2m
# 全局环境变量, 可选
env:
  - name: Test
    value: "test_env"
# 步骤列表, 不能为空
step:
    # 步骤名称, 唯一, 可选[当自定义编排是必须设置], 默认自动生成
  - name: 步骤2
    # 描述
    desc: 这是一段步骤描述
    # 超时时间, 可选, 默认任务级超时时间
    timeout: 2m
    # 禁用, 可选, 默认false
    #disable: false
    # 依赖步骤, 可选[自定义编排时用到]
    depends:
      - 步骤1
    # 局部环境变量, 会覆盖同名的全局变量
    env:
      - name: Test
        value: "test_env"
    # 类型
    type: sh
    # 内容
    content: |-
      ping 1.1.1.1
  - name: 步骤1
    # 描述
    desc: 这是一段步骤描述
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    retryPolicy: # 重试策略, 可选
      maxAttempts: 0 # 最大重试次数, 默认0(不重试)
      interval: 1s # 重试间隔, 默认1秒
      maxInterval: 30s # 最大重试间隔, 默认30秒
      multiplier: 2.0  # 重试间隔倍数, 默认2.0
    type: sh
    content: |-
      ping -c 4 1.1.1.1
  - name: python脚本
    desc: 测试python脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    type: python
    content: |-
      #!/usr/bin/env python3
      import sys
      import subprocess
      import platform
      
      # 输出 Hello World
      print("Hello World")
      
      # 获取并打印命令行参数
      print(f"\\n脚本路径: {sys.argv[0]}")
      print(f"参数数量: {len(sys.argv) - 1}")
      if len(sys.argv) > 1:
          print("参数列表:")
          for i, arg in enumerate(sys.argv[1:], 1):
              print(f"  参数 {i}: {arg}")
      else:
          print("没有传递任何参数")
      print()
      # 从标准输入读取每一行（支持管道和非交互模式）
      try:
          while True:
              line = sys.stdin.readline()
              # 如果读到EOF或空行则退出
              if not line or line.strip() == '':
                  break
              # 去掉末尾的换行符并输出
              print(f"Line: {line.rstrip()}")
      except (EOFError, KeyboardInterrupt):
          # 如果没有输入流或用户中断，静默处理
          pass
      
      # 持续 ping
      try:
          # 根据操作系统选择合适的 ping 参数
          if platform.system().lower() == 'windows':
              # Windows: -t 表示持续 ping
              subprocess.run(['ping', '-t', '1.1.1.1'])
          else:
              # Linux/Mac: 不需要 -t 参数，默认就是持续 ping
              subprocess.run(['ping', '1.1.1.1'])
      except KeyboardInterrupt:
          print("\\nPing 已停止")
      
  - name: shell脚本参数
    desc: 测试python脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    type: sh
    content: |-
      echo "Hello World"
      echo "脚本路径: $0"
      echo "参数1: $1"
      echo "参数2: $2"
      echo "所有参数: $@"
      echo "参数数量: $#"
      while read line; do
        echo "Line: $line"
      done
      exit 123
      # 持续ping
      ping 1.1.1.1`,
            pipelineParams: `params:\n  # 在这里添加流水线运行参数\n  imageTag: latest`
        };
    }

    static get STATUS_COLORS() {
        return {
            'running': '#00AEEF',
            'stopped': '#00C853',
            'failed': '#D50000',
            'pending': '#FF9800',
            'timeout': '#FF5722',
            'canceled': '#7C4DFF',
            'skipped': '#3F51B5',
            'unknown': '#9E9E9E'
        };
    }

    static get PAGINATION_DEFAULTS() {
        return {
            page: 1,
            size: 15,
            sizes: [15, 25, 35, 50]
        };
    }
}

// ==================== 工具类 ====================
class Utils {
    static removeElementById(id) {
        const element = document.getElementById(id);
        if (element) element.remove();
    }

    static getStatusColor(status) {
        return ConfigManager.STATUS_COLORS[status] || ConfigManager.STATUS_COLORS.unknown;
    }

    static escapeHTML(html) {
        const div = document.createElement('div');
        div.textContent = html;
        return div.innerHTML;
    }

    static showToast(message, type = 'info') {
        const toast = document.createElement('div');
        toast.className = 'toast';
        toast.style.cssText = `
            position: fixed;
            top: 20px;
            left: 50%;
            padding: 12px 20px;
            border-radius: 8px;
            color: white;
            font-weight: 500;
            z-index: 10000;
            opacity: 0;
            transform: translateX(-50%) translateY(-100%);
            transition: all 0.3s ease;
        `;

        const colors = {
            success: '#10b981',
            error: '#ef4444',
            warning: '#f59e0b',
            info: '#06b6d4'
        };

        toast.style.background = colors[type] || colors.info;
        toast.textContent = message;

        document.body.appendChild(toast);

        setTimeout(() => {
            toast.style.opacity = '1';
            toast.style.transform = 'translateX(-50%) translateY(0)';
        }, 10);

        setTimeout(() => {
            toast.style.opacity = '0';
            toast.style.transform = 'translateX(-50%) translateY(-100%)';
            setTimeout(() => toast.remove(), 300);
        }, 1000);
    }

    static async initializeEditor(elementID, data, readOnly = false) {
        return new Promise((resolve, reject) => {
            // 确保monaco已加载
            if (!window.monaco) {
                reject(new Error('Monaco Editor 未加载'));
                return;
            }

            try {
                const editor = window.monaco.editor.create(document.getElementById(elementID), {
                    value: data,
                    language: 'yaml',
                    theme: 'vs-dark',
                    readOnly: readOnly,
                    autoIndent: 'full',
                    automaticLayout: true,
                    overviewRulerBorder: false,
                    foldingStrategy: 'indentation',
                    lineNumbers: 'on',
                    minimap: { enabled: false },
                    tabSize: 2,
                    mouseWheelZoom: true,
                    formatOnType: true,
                    formatOnPaste: true,
                    cursorStyle: 'line',
                    fontSize: 13,
                    lineHeight: 20,
                    scrollBeyondLastLine: true,
                    wordWrap: 'on',
                });
                resolve(editor);
            } catch (error) {
                reject(error);
            }
        });
    }

    static getNestedValue(obj, path) {
        return path.split('.').reduce((current, key) => current?.[key], obj);
    }
}

// ==================== API管理器 ====================
class APIManager {
    static async request(url, options = {}) {
        try {
            const response = await fetch(url, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                ...options
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API请求失败:', error);
            throw error;
        }
    }

    static async taskAction(taskName, action) {
        const confirmed = confirm(`确定要${action} "${taskName}"?`);
        if (!confirmed) return;

        try {
            const result = await this.request(`${baseUrl}${ConfigManager.API_ENDPOINTS.task}/${taskName}?action=${action}`, {
                method: 'PUT'
            });
            Utils.showToast(`${action}操作成功`, 'success');
            return result;
        } catch (error) {
            Utils.showToast(`${action}操作失败: ${error.message}`, 'error');
            throw error;
        }
    }

    static async stepAction(taskName, stepName, action) {
        const confirmed = confirm(`确定要${action} "${stepName}"?`);
        if (!confirmed) return;

        try {
            const result = await this.request(`${baseUrl}${ConfigManager.API_ENDPOINTS.task}/${taskName}/step/${stepName}?action=${action}`, {
                method: 'PUT'
            });
            Utils.showToast(`${action}操作成功`, 'success');
            return result;
        } catch (error) {
            Utils.showToast(`${action}操作失败: ${error.message}`, 'error');
            throw error;
        }
    }

    static async deleteResource(url, name, type = '资源') {
        const confirmed = confirm(`确定要删除${type} "${name}"?`);
        if (!confirmed) return;

        try {
            await this.request(url, { method: 'DELETE' });
            Utils.showToast(`${type}删除成功`, 'success');
        } catch (error) {
            Utils.showToast(`${type}删除失败: ${error.message}`, 'error');
            throw error;
        }
    }

    static async createResource(url, data, type = '资源') {
        try {
            const result = await this.request(url, {
                method: 'POST',
                headers: { 'Content-Type': 'application/yaml' },
                body: data
            });

            if (result.code === 0) {
                Utils.showToast(`${type}创建成功`, 'success');
                return result;
            } else {
                throw new Error(result.message);
            }
        } catch (error) {
            Utils.showToast(`${type}创建失败: ${error.message}`, 'error');
            throw error;
        }
    }
}

// ==================== WebSocket管理器 ====================
class WebSocketManager {
    constructor(url, onMessage, onError = null, options = {}) {
        this.url = url;
        this.onMessage = onMessage;
        this.onError = onError;
        this.reconnectInterval = options.reconnectInterval || 5000;
        this.maxReconnectAttempts = options.maxReconnectAttempts || 5;
        this.socket = null;
        this.isManuallyClosed = false;
        this.reconnectAttempts = 0;
        this.connect();
    }

    connect() {
        if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
            return;
        }

        this.isManuallyClosed = false;
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            console.log("WebSocket连接已建立");
            this.reconnectAttempts = 0;
        };

        this.socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.onMessage(data);
            } catch (error) {
                console.error('解析WebSocket消息失败:', error);
            }
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket错误:', error);
            if (this.onError) this.onError(error);
        };

        this.socket.onclose = (event) => {
            if (event.wasClean) {
                console.log("WebSocket正常关闭");
            } else {
                console.log(`WebSocket意外关闭: ${event.code}, ${event.reason}`);
                this.handleReconnect();
            }
        };
    }

    handleReconnect() {
        if (this.isManuallyClosed || this.reconnectAttempts >= this.maxReconnectAttempts) return;

        this.reconnectAttempts++;
        const delay = Math.min(this.reconnectInterval * Math.pow(2, this.reconnectAttempts - 1), 30000);

        setTimeout(() => {
            console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
            this.connect();
        }, delay);
    }

    send(data) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        } else {
            console.warn("WebSocket未连接，无法发送数据");
        }
    }

    close() {
        this.isManuallyClosed = true;
        if (this.socket) {
            this.socket.close(1000, "正常关闭");
        }
    }
}

// ==================== 基础模态框类 ====================
class BaseModal {
    constructor(options = {}) {
        this.options = {
            id: options.id || 'modal',
            title: options.title || '模态框',
            size: options.size || 'default', // default, large, small
            ...options
        };
        this.element = null;
        this.overlay = null;
        this.keydownHandler = null;
    }

    create() {
        this.cleanup();
        this.createElement();
        this.addEventListeners();
        this.show();
    }

    createElement() {
        // 创建遮罩
        this.overlay = document.createElement('div');
        this.overlay.id = `${this.options.id}-modal-overlay`;
        this.overlay.className = 'modal-overlay';

        // 创建模态框
        this.element = document.createElement('div');
        this.element.id = `${this.options.id}-card`;
        this.element.className = 'card-one';

        // 根据尺寸设置样式
        this.applySizeStyles();

        this.element.innerHTML = this.getContent();

        document.body.appendChild(this.overlay);
        document.body.appendChild(this.element);
    }

    applySizeStyles() {
        const sizeStyles = {
            small: { left: '30%', right: '30%' },
            default: { left: '25%', right: '15%' },
            large: { left: '0%', right: '0%' }
        };

        const style = sizeStyles[this.options.size] || sizeStyles.default;
        Object.assign(this.element.style, style);
    }

    getContent() {
        return `
            <div class="card-header">
                <h5>${this.options.title}</h5>
                <div style="display: flex; gap: 8px;">
                    ${this.getHeaderButtons()}
                    <span class="card-close" id="${this.options.id}-close">&times;</span>
                </div>
            </div>
            <div class="card-body-content">
                ${this.getBodyContent()}
            </div>
        `;
    }

    getHeaderButtons() {
        return '';
    }

    getBodyContent() {
        return '<p>模态框内容</p>';
    }

    addEventListeners() {
        // 关闭按钮
        document.getElementById(`${this.options.id}-close`)
            ?.addEventListener('click', () => this.close());

        // 遮罩点击关闭
        this.overlay?.addEventListener('click', () => this.close());

        // ESC键关闭
        this.keydownHandler = (event) => {
            if (event.key === 'Escape') this.close();
        };
        window.addEventListener('keydown', this.keydownHandler);

        // 添加自定义事件监听器
        this.bindCustomEvents();
    }

    bindCustomEvents() {
        // 子类重写此方法
    }

    show() {
        setTimeout(() => {
            this.overlay?.classList.add('show');
            this.element?.classList.add('show');
        }, 10);
    }

    close() {
        this.element?.classList.remove('show');
        this.overlay?.classList.remove('show');

        setTimeout(() => {
            this.cleanup();
        }, 300);
    }

    cleanup() {
        Utils.removeElementById(`${this.options.id}-card`);
        Utils.removeElementById(`${this.options.id}-modal-overlay`);

        if (this.keydownHandler) {
            window.removeEventListener('keydown', this.keydownHandler);
            this.keydownHandler = null;
        }
    }
}

// ==================== 基础表格类 ====================
class BaseTable {
    constructor(options = {}) {
        this.options = options;
        this.webSocketManager = null;
        this.currentPage = ConfigManager.PAGINATION_DEFAULTS.page;
        this.rowsPerPage = ConfigManager.PAGINATION_DEFAULTS.size;
        this.data = [];
        this.totalPage = 0;
        this.previousData = [];
        this.init();
    }

    init() {
        this.setupWebSocket();
        this.setupEventListeners();
    }

    setupWebSocket() {
        if (!this.getWebSocketUrl()) return;

        this.webSocketManager = new WebSocketManager(
            this.getWebSocketUrl(),
            this.handleWebSocketData.bind(this)
        );
    }

    getWebSocketUrl() {
        // 子类实现
        return null;
    }

    handleWebSocketData(res) {
        if (res.data && this.getDataFromResponse(res)) {
            this.data = this.getDataFromResponse(res);
            this.totalPage = res.data.page?.total || 1;
            this.currentPage = res.data.page?.current || 1;
            this.rowsPerPage = res.data.page?.size || ConfigManager.PAGINATION_DEFAULTS.size;
            this.updateTable();
            this.updatePagination();
        } else {
            this.data = [];
            this.totalPage = 1;
            this.updateTable();
            this.updatePagination();
        }
    }

    getDataFromResponse(res) {
        // 子类实现
        return [];
    }

    updateTable() {
        const tableBody = this.getTableBody();
        if (!tableBody) return;

        if (this.data.length === 0) {
            tableBody.innerHTML = this.getEmptyContent();
            return;
        }

        this.renderRows(tableBody);
        this.previousData = JSON.parse(JSON.stringify(this.data));
    }

    getTableBody() {
        // 子类实现
        return null;
    }

    getEmptyContent() {
        return `
            <tr>
                <td colspan="${this.getColumnCount()}" style="text-align: center; padding: 40px; color: var(--gray-500);">
                    暂无数据
                </td>
            </tr>
        `;
    }

    getColumnCount() {
        // 子类实现
        return 1;
    }

    renderRows(tableBody) {
        // 子类实现
    }

    updatePagination() {
        const prevBtn = document.getElementById(this.getPrevButtonId());
        const nextBtn = document.getElementById(this.getNextButtonId());

        if (prevBtn) prevBtn.disabled = this.currentPage === 1;
        if (nextBtn) nextBtn.disabled = this.currentPage === this.totalPage;

        const pageInfo = document.getElementById(this.getPageInfoId());
        if (pageInfo) {
            pageInfo.textContent = `第${this.currentPage}页 / 共${this.totalPage}页`;
        }
    }

    getPrevButtonId() {
        // 子类实现
        return '';
    }

    getNextButtonId() {
        // 子类实现
        return '';
    }

    getPageInfoId() {
        // 子类实现
        return '';
    }

    setupEventListeners() {
        const prevBtn = document.getElementById(this.getPrevButtonId());
        const nextBtn = document.getElementById(this.getNextButtonId());
        const sizeSelect = document.getElementById(this.getPageSizeSelectId());

        prevBtn?.addEventListener("click", () => {
            if (this.currentPage > 1) {
                this.currentPage--;
                this.fetchData();
            }
        });

        nextBtn?.addEventListener("click", () => {
            if (this.currentPage < this.totalPage) {
                this.currentPage++;
                this.fetchData();
            }
        });

        sizeSelect?.addEventListener("change", (event) => {
            this.rowsPerPage = parseInt(event.target.value);
            this.currentPage = 1;
            this.fetchData();
        });
    }

    getPageSizeSelectId() {
        // 子类实现
        return '';
    }

    fetchData() {
        const request = {
            page: this.currentPage,
            size: this.rowsPerPage,
        };
        this.webSocketManager?.send(request);
    }
}

// ==================== G6图形组件 ====================
class G6GraphManager {
    constructor() {
        this.setupCustomNode();
    }

    setupCustomNode() {
        // 确保G6已加载
        if (!window.G6) {
            console.error('G6 未加载，无法注册自定义节点');
            return;
        }

        // 检查是否已注册过，避免重复注册
        if (window.G6._customNodeRegistered) {
            return;
        }

        // 注册自定义节点
        window.G6.registerNode('custom-node', {
            drawShape: function drawShape(cfg, group) {
                const color = Utils.getStatusColor(cfg.step.state || 'unknown');
                const r = 8;
                // 主容器
                group.addShape('rect', {
                    attrs: {
                        x: 0, y: 0, width: 240, height: 80,
                        stroke: color, fill: "#FFFFFF",
                        radius: r, lineWidth: 2,
                    },
                    name: 'main-box',
                    draggable: true,
                });
                // 标题栏
                group.addShape('rect', {
                    attrs: {
                        x: 0, y: 0, width: 240, height: 24,
                        fill: color, radius: [r, r, 0, 0],
                    },
                    name: 'title-box',
                    draggable: true,
                });
                // 图标
                group.addShape('image', {
                    attrs: {
                        x: 6, y: 4, height: 16, width: 16,
                        cursor: 'pointer',
                        img: `${baseUrl}/img/node-icon.png`,
                    },
                    name: 'node-icon',
                });
                // 标题文字
                group.addShape('text', {
                    attrs: {
                        textBaseline: 'top', y: 8, x: 28,
                        lineHeight: 16, text: cfg.step.name,
                        fill: '#fff', fontWeight: 500,
                    },
                    name: 'title-text',
                });
                // 状态码
                if (cfg.step.code !== undefined) {
                    group.addShape('text', {
                        attrs: {
                            textBaseline: 'top', y: 32, x: 8,
                            lineHeight: 16, text: "状态码: " + cfg.step.code,
                            fill: 'rgba(0,0,0, 0.6)', fontSize: 11,
                        },
                        name: 'title-code',
                    });
                }
                // 状态
                group.addShape('text', {
                    attrs: {
                        textBaseline: 'top', y: 32, x: 120,
                        lineHeight: 16, text: "状态: " + cfg.step.state,
                        fill: 'rgba(0,0,0, 0.6)', fontSize: 11,
                    },
                    name: 'title-state',
                });
                // 时间信息
                group.addShape('text', {
                    attrs: {
                        textBaseline: 'top', y: 48, x: 8,
                        lineHeight: 16,
                        text: cfg.step.time.start ? "开始: " + cfg.step.time.start || '---' : '开始: ---',
                        fill: 'rgba(0,0,0, 0.6)', fontSize: 10,
                    },
                    name: 'title-time-start',
                });
                group.addShape('text', {
                    attrs: {
                        textBaseline: 'top', y: 64, x: 8,
                        lineHeight: 16,
                        text: cfg.step.time.end ? "结束: " + cfg.step.time.end || '---' : '结束: ---',
                        fill: 'rgba(0,0,0, 0.6)', fontSize: 10,
                    },
                    name: 'title-time-end',
                });
                return group;},
        });
        window.G6._customNodeRegistered = true;
    }

    createGraph(containerId, data, options = {}) {
        const container = document.getElementById(containerId);
        const width = container.scrollWidth;
        const height = container.scrollHeight || 500;

        const grid = new G6.Grid();
        const menu = new G6.Menu({
            itemTypes: ['node'],
            getContent(e) {
                const code = e.item.getModel().code;
                const step = e.item.getModel().step;

                if (!step || code === 0 || code === 1002 || code === 1003) {
                    return '<div style="padding: 8px; color: #6b7280;">无操作可选</div>';
                }

                let menuContent = '';
                if (step.state === 'running') {
                    menuContent += '<a href="#" id="kill-step">强杀</a>';
                } else if (step.state === 'paused') {
                    menuContent += '<a href="#" id="resume-step">解挂</a>';
                    menuContent += '<a href="#" id="kill-step">强杀</a>';
                } else if (step.state === 'pending') {
                    menuContent += '<a href="#" id="pause-step">挂起</a>';
                    menuContent += '<a href="#" id="kill-step">强杀</a>';
                }

                return menuContent || '<div style="padding: 8px; color: #6b7280;">无操作可选</div>';
            },
            handleMenuClick(target, item) {
                const taskName = item.getModel().taskName;
                const step = item.getModel().step;

                if (target.id === 'kill-step') {
                    APIManager.stepAction(taskName, step.name, 'kill');
                } else if (target.id === 'pause-step') {
                    APIManager.stepAction(taskName, step.name, 'pause');
                } else if (target.id === 'resume-step') {
                    APIManager.stepAction(taskName, step.name, 'resume');
                }
            },
        });

        const toolbar = new G6.ToolBar({
            getContent: () => `
                <ul class="g6-component-toolbar">
                    <li code="zoomOut" title="放大">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
                            <path d="m12 10h-2v2H9v-2H7V9h2V7h1v2h2v1z"/>
                        </svg>
                    </li>
                    <li code="zoomIn" title="缩小">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
                            <path d="M7 9h5v1H7z"/>
                        </svg>
                    </li>
                    <li code="autoZoom" title="适应画布">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                        </svg>
                    </li>
                </ul>
            `,
            position: { x: width - 120, y: 10 }
        });

        const graph = new G6.Graph({
            container: containerId,
            height: height,
            width: width,
            layout: {
                type: 'dagre',
                rankdir: 'LR',
                align: 'UL',
                nodesep: 32,
                ranksep: 96,
            },
            modes: {
                default: ['drag-canvas', 'drag-node', 'zoom-canvas', 'activate-relations'],
            },
            defaultNode: {
                type: 'custom-node',
                anchorPoints: [[0, 0.5], [1, 0.5]],
            },
            defaultEdge: {
                type: 'polyline',
                style: {
                    endArrow: {
                        path: G6.Arrow.triangle(8, 8, 12),
                        d: 12
                    },
                    stroke: '#2563eb',
                    lineWidth: 2,
                },
            },
            plugins: [grid, menu, toolbar],
        });

        if (options.onNodeClick) {
            graph.on('node:click', options.onNodeClick);
        }

        if (options.onCanvasClick) {
            graph.on('canvas:click', options.onCanvasClick);
        }

        graph.data(data);
        graph.render();

        return graph;
    }
}

// ==================== 任务相关模态框 ====================
class TaskModal extends BaseModal {
    constructor(taskName) {
        super({
            id: 'task',
            title: `任务: ${taskName}`,
            size: 'large'
        });

        this.taskName = taskName;
        this.task = null;
        this.graph = null;
        this.stepModals = [];
        this.webSocketManager = null;
        this.graphManager = new G6GraphManager();

        this.init();
    }

    async init() {
        try {
            const response = await APIManager.request(`${baseUrl}${ConfigManager.API_ENDPOINTS.task}/${this.taskName}`);
            this.task = response.data;
            this.create();
            this.startWebSocket();
        } catch (error) {
            Utils.showToast(`获取任务详情失败: ${error.message}`, 'error');
            throw error;
        }
    }

    startWebSocket() {
        this.webSocketManager = new WebSocketManager(
            `${ConfigManager.WS_ENDPOINTS.task}/${this.taskName}/step`,
            this.updateGraphData.bind(this),
            (error) => {
                Utils.showToast('WebSocket连接错误', 'error');
                this.close();
            }
        );
    }

    getBodyContent() {
        return `
            <div class="card-body">
                ${this.task.env ? `
                    <div class="card-body-left">
                        <h6 style="margin-bottom: 12px; color: var(--gray-700);">环境变量</h6>
                        <div class="env">${this.task.env.map(env => `${env.name}: ${env.value}`).join('\n')}</div>
                    </div>
                ` : ''}
                <div class="card-body-right" id="task-card-left"></div>
            </div>
        `;
    }

    updateGraphData = (res) => {
        if (!res.data) return;

        if (this.graph) {
            res.data.forEach(step => {
                let node = {
                    id: step.name,
                    step: step,
                    code: res.code,
                    taskName: this.taskName,
                };

                const item = this.graph.findById(node.id);
                if (item) {
                    this.graph.updateItem(node.id, node);
                } else {
                    this.graph.addItem('node', node);
                }

                if (step.depends) {
                    step.depends.forEach(depend => {
                        let edge = {
                            id: depend + '-' + step.name,
                            source: depend,
                            target: step.name,
                        };

                        if (!this.graph.findById(edge.id)) {
                            this.graph.addItem('edge', edge);
                        }
                    });
                }
            });
            return;
        }

        // 创建新图表
        let data = { nodes: [], edges: [] };

        res.data.forEach(step => {
            data.nodes.push({
                id: step.name,
                step: step,
                code: res.code,
                taskName: this.taskName,
            });

            if (step.depends) {
                step.depends.forEach(depend => {
                    data.edges.push({
                        id: depend + '-' + step.name,
                        source: depend,
                        target: step.name,
                    });
                });
            }
        });

        this.graph = this.graphManager.createGraph("task-card-left", data, {
            onNodeClick: (evt) => {
                const step = evt.item.getModel().step;
                this.openStepModal(step.name);
            },
            onCanvasClick: () => {
                this.closeAllStepModals();
            }
        });
    }

    openStepModal(stepName) {
        const stepModal = new StepModal(this.task.name, stepName);
        this.stepModals.push(stepModal);
    }

    closeAllStepModals() {
        this.stepModals.forEach(stepModal => stepModal.close());
        this.stepModals = [];
    }

    close() {
        if (this.webSocketManager) {
            this.webSocketManager.close();
            this.webSocketManager = null;
        }
        this.closeAllStepModals();
        super.close();
    }
}

class StepModal {
    constructor(taskName, stepName) {
        this.taskName = taskName;
        this.stepName = stepName;
        this.step = null;
        this.webSocketManager = null;
        this.isDragging = false;
        this.offsetX = 0;
        this.offsetY = 0;
        this.init();
    }

    async init() {
        try {
            const response = await APIManager.request(`${baseUrl}${ConfigManager.API_ENDPOINTS.task}/${this.taskName}/step/${this.stepName}`);
            this.step = response.data;
            this.create();
            this.startWebSocket();
        } catch (error) {
            Utils.showToast(`获取步骤详情失败: ${error.message}`, 'error');
            throw error;
        }
    }

    startWebSocket() {
        this.webSocketManager = new WebSocketManager(
            `${ConfigManager.WS_ENDPOINTS.task}/${this.taskName}/step/${this.stepName}/log`,
            this.updateStepOutput.bind(this),
            () => {
                const outputElement = document.getElementById(this.step.name + '-step-output-text');
                if (outputElement) {
                    outputElement.textContent = this.step.message || '无输出';
                }
            }
        );
    }

    create() {
        const existingCard = document.getElementById(this.step.name + "-step-card");
        if (existingCard) {
            existingCard.style.zIndex = ++window.highestZIndex || 1000;
            return;
        }

        const card = document.createElement('div');
        card.setAttribute("id", this.step.name + "-step-card");
        card.className = 'step-card';
        card.style.zIndex = ++window.highestZIndex || 1000;
        card.innerHTML = `
            <div id="${this.step.name + '-step-card-header'}" class="card-header">
                <div>
                    <h6 style="margin: 0; margin-bottom: 4px;">步骤: ${this.step.name}</h6>
                    <span style="font-size: 12px; color: var(--gray-500);">类型: ${this.step.type}</span>
                </div>
                <span class="card-close" id="${this.step.name + '-close-step-card'}">&times;</span>
            </div>
            <div id="${this.step.name + '-step-card-body'}" class="card-body">
                ${this.step.env ? `
                    <div class="card-body-left">
                        <h6 style="margin-bottom: 12px; color: var(--gray-700);">环境变量</h6>
                        <div class="env">${this.step.env.map(env => `${env.name}: ${env.value}`).join('\n')}</div>
                    </div>
                ` : ''}
                <div class="card-body-right">
                    <h6 style="margin-bottom: 8px; color: var(--gray-700);">输入内容</h6>
                    <div class="step-card-input">
                        <pre class="step-card-code">${this.step.content}</pre>
                    </div>
                    <h6 style="margin-bottom: 8px; color: var(--gray-700);">执行输出</h6>
                    <div class="step-card-output">
                        <pre id="${this.step.name + '-step-output-text'}" class="step-card-code"></pre>
                    </div>
                </div>
            </div>
        `;

        document.body.appendChild(card);
        this.addEventListeners(card);
    }

    addEventListeners(card) {
        document.getElementById(this.step.name + '-close-step-card')
            ?.addEventListener('click', () => this.close());

        const header = document.getElementById(this.step.name + '-step-card-header');
        this.addDragEventListeners(header, card);

        card.addEventListener('mousedown', () => {
            card.style.zIndex = ++window.highestZIndex || 1000;
        });
    }

    addDragEventListeners(header, card) {
        header.style.cursor = 'move';

        const handleMouseDown = (e) => {
            this.isDragging = true;
            this.offsetX = e.clientX - card.offsetLeft;
            this.offsetY = e.clientY - card.offsetTop;
            document.addEventListener('mousemove', this.handleMouseMove);
            document.addEventListener('mouseup', this.handleMouseUp);
        };

        header.addEventListener('mousedown', handleMouseDown);
    }

    handleMouseMove = (e) => {
        if (this.isDragging) {
            const card = document.getElementById(this.step.name + "-step-card");
            const newLeft = Math.max(0, Math.min(window.innerWidth - card.offsetWidth, e.clientX - this.offsetX));
            const newTop = Math.max(0, Math.min(window.innerHeight - card.offsetHeight, e.clientY - this.offsetY));

            card.style.left = newLeft + 'px';
            card.style.top = newTop + 'px';
        }
    }

    handleMouseUp = () => {
        this.isDragging = false;
        document.removeEventListener('mousemove', this.handleMouseMove);
        document.removeEventListener('mouseup', this.handleMouseUp);
    }

    updateStepOutput = (res) => {
        const outputElement = document.getElementById(this.step.name + '-step-output-text');
        if (!outputElement || !res.data) return;

        const scrollTop = outputElement.scrollTop;
        const scrollHeight = outputElement.scrollHeight;
        const clientHeight = outputElement.clientHeight;
        const isScrolledToBottom = scrollTop + clientHeight >= scrollHeight - 5;

        const data = res.data.map(item => item.content).join("\n");
        outputElement.textContent += data + "\n";

        if (isScrolledToBottom) {
            outputElement.scrollTop = outputElement.scrollHeight;
        }
    }

    close() {
        if (this.webSocketManager) {
            this.webSocketManager.close();
            this.webSocketManager = null;
        }
        Utils.removeElementById(this.step.name + "-step-card");
    }
}

class TaskFormModal extends BaseModal {
    constructor() {
        super({
            id: 'add-task',
            title: '创建新任务',
            size: 'default'
        });

        this.editor = null;
        this.create();
    }

    getHeaderButtons() {
        return `
            <button id="create-task" class="button-sure">创建任务</button>
            <button id="cancel-task" class="button-cancel">取消</button>
        `;
    }

    getBodyContent() {
        return `
            <div class="create-content">
                <label for="yaml-editor">任务配置 (YAML)</label>
                <div id="yaml-editor"></div>
            </div>
        `;
    }

    bindCustomEvents() {
        Utils.initializeEditor('yaml-editor', "# 任务名称, 可选, 默认自动生成\nname: 测试任务\n" + ConfigManager.TEMPLATES.taskYaml)
            .then(editor => this.editor = editor);

        document.getElementById("create-task")?.addEventListener("click", () => this.createTask());
        document.getElementById("cancel-task")?.addEventListener("click", () => this.close());
    }

    async createTask() {
        const yamlContent = this.editor.getValue().trim();

        if (!yamlContent) {
            Utils.showToast("请输入YAML内容", 'warning');
            return;
        }

        try {
            const result = await APIManager.createResource(`${baseUrl}${ConfigManager.API_ENDPOINTS.task}`, yamlContent, '任务');
            new TaskModal(result.data.name);
            this.close();
        } catch (error) {
            // 错误已在APIManager中处理
        }
    }
}

// ==================== 流水线相关模态框 ====================
class PipelineFormModal extends BaseModal {
    constructor(pipelineData = null) {
        super({
            id: pipelineData ? 'edit-pipeline' : 'add-pipeline',
            title: pipelineData ? `编辑流水线: ${pipelineData.name}` : '创建新流水线',
            size: 'default'
        });

        this.pipelineData = pipelineData;
        this.editor = null;
        this.isEdit = !!pipelineData;

        if (this.isEdit) {
            this.create();
        } else {
            this.create();
        }
    }

    getHeaderButtons() {
        const action = this.isEdit ? '保存更改' : '创建流水线';
        return `
            <button id="${this.isEdit ? 'edit' : 'create'}-pipeline" class="button-sure">${action}</button>
            <button id="cancel-${this.isEdit ? 'edit' : 'create'}-pipeline" class="button-cancel">取消</button>
        `;
    }

    getBodyContent() {
        return `
            <div class="create-content">
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 16px;">
                    <div>
                        <label for="name">流水线名称</label>
                        <input type="text" id="name" name="name" placeholder="请输入流水线名称" 
                               value="${this.pipelineData?.name || ''}" ${this.isEdit ? 'readonly style="background: var(--gray-100);"' : ''}>
                    </div>
                    <div>
                        <label for="tplType">模板类型</label>
                        <select id="tplType" name="tplType">
                            <option value="jinja2" ${this.pipelineData?.tplType === 'jinja2' ? 'selected' : ''}>jinja2</option>
                        </select>
                    </div>
                </div>
                <div style="margin-bottom: 16px;">
                    <label for="description">流水线描述</label>
                    <textarea id="description" name="description" placeholder="请输入流水线描述...">${this.pipelineData?.desc || ''}</textarea>
                </div>
                <div style="flex: 1;display: contents;bottom: 0;top: 0;position: fixed;flex-direction: column;">
                    <label for="yaml-editor">流水线内容 (YAML)</label>
                    <div id="yaml-editor" style="flex: 1; min-height: 400px;"></div>
                </div>
            </div>
        `;
    }

    bindCustomEvents() {
        const content = this.pipelineData?.content || ConfigManager.TEMPLATES.taskYaml;
        Utils.initializeEditor('yaml-editor', content)
            .then(editor => this.editor = editor);

        const actionButton = document.getElementById(`${this.isEdit ? 'edit' : 'create'}-pipeline`);
        const cancelButton = document.getElementById(`cancel-${this.isEdit ? 'edit' : 'create'}-pipeline`);

        actionButton?.addEventListener("click", () => this.submitPipeline());
        cancelButton?.addEventListener("click", () => this.close());
    }

    async submitPipeline() {
        const name = document.getElementById("name").value.trim();
        const description = document.getElementById("description").value.trim();
        const tplType = document.getElementById("tplType").value;
        const content = this.editor.getValue().trim();

        if (!name && !this.isEdit) {
            Utils.showToast("流水线名称不能为空", 'warning');
            return;
        }

        if (!content) {
            Utils.showToast("请输入YAML内容", 'warning');
            return;
        }

        const yamlData = this.isEdit
            ? `desc: |-\n    ${description.replace(/\n/g, '\n    ')}\ntplType: ${tplType}\ncontent: |-\n    ${content.replace(/\n/g, '\n    ')}`
            : `name: ${name}\ndesc: |-\n    ${description.replace(/\n/g, '\n    ')}\ntplType: ${tplType}\ncontent: |-\n    ${content.replace(/\n/g, '\n    ')}`;

        try {
            const url = this.isEdit
                ? `${baseUrl}${ConfigManager.API_ENDPOINTS.pipeline}/${this.pipelineData.name}`
                : `${baseUrl}${ConfigManager.API_ENDPOINTS.pipeline}`;

            await APIManager.createResource(url, yamlData, '流水线');
            this.close();
        } catch (error) {
            // 错误已在APIManager中处理
        }
    }
}

class RunPipelineModal extends BaseModal {
    constructor(pipelineName) {
        super({
            id: 'run-pipeline',
            title: `运行流水线: ${pipelineName}`,
            size: 'default'
        });

        this.pipelineName = pipelineName;
        this.editor = null;
        this.create();
    }

    getHeaderButtons() {
        return `
            <button id="run-start-pipeline" class="button-sure">开始执行</button>
            <button id="cancel-run-pipeline" class="button-cancel">取消</button>
        `;
    }

    getBodyContent() {
        return `
            <div class="create-content">
                <label for="yaml-editor">运行参数 (YAML)</label>
                <div id="yaml-editor"></div>
            </div>
        `;
    }

    bindCustomEvents() {
        Utils.initializeEditor('yaml-editor', ConfigManager.TEMPLATES.pipelineParams)
            .then(editor => this.editor = editor);

        document.getElementById("run-start-pipeline")?.addEventListener("click", () => this.runPipeline());
        document.getElementById("cancel-run-pipeline")?.addEventListener("click", () => this.close());
    }

    async runPipeline() {
        const content = this.editor.getValue().trim();

        try {
            const result = await APIManager.createResource(
                `${baseUrl}${ConfigManager.API_ENDPOINTS.pipeline}/${this.pipelineName}/build`,
                content || 'params: {}',
                '流水线执行'
            );
            new TaskModal(result.data.name);
            this.close();
        } catch (error) {
            // 错误已在APIManager中处理
        }
    }
}

class PipelineDetailModal extends BaseModal {
    constructor(pipelineName) {
        super({
            id: 'pipeline-detail',
            title: `流水线: ${pipelineName}`,
            size: 'default'
        });

        this.pipelineName = pipelineName;
        this.pipeline = null;
        this.webSocketManager = null;
        this.currentPage = ConfigManager.PAGINATION_DEFAULTS.page;
        this.rowsPerPage = ConfigManager.PAGINATION_DEFAULTS.size;
        this.totalPage = 0;
        this.tasks = [];
        this.editor = null;

        this.init();
    }

    async init() {
        try {
            const response = await APIManager.request(`${baseUrl}${ConfigManager.API_ENDPOINTS.pipeline}/${this.pipelineName}`);
            this.pipeline = response.data;
            this.create();
            this.startWebSocket();
        } catch (error) {
            Utils.showToast(`获取流水线详情失败: ${error.message}`, 'error');
            throw error;
        }
    }

    startWebSocket() {
        this.webSocketManager = new WebSocketManager(
            `${ConfigManager.WS_ENDPOINTS.pipeline}/${this.pipelineName}/build`,
            this.updateTaskData.bind(this),
            (error) => {
                Utils.showToast('WebSocket连接错误', 'error');
                this.close();
            }
        );
    }

    getHeaderButtons() {
        return `
            <button id="pipeline-view-detail-bt" class="button-sure">详情</button>
            <button id="pipeline-view-list-bt" class="button-sure">任务列表</button>
        `;
    }

    getBodyContent() {
        return `
            <div id="pipeline-view-detail" style="position: fixed;inset: 81px 24px 24px;overflow: hidden;">
                <div style="margin-bottom: 24px;">
                    <label for="name">流水线名称</label>
                    <div style="font-weight: 500;">${this.pipelineName}</div>
                </div>
                <div style="margin-bottom: 24px;">
                    <label for="description">流水线描述</label>
                    <div class="env" style="min-height: 60px;">${this.pipeline?.desc || '暂无描述'}</div>
                </div>
                <div style="flex: 1;display: contents;bottom: 0;top: 0;position: fixed;flex-direction: column;">
                    <label for="yaml-editor">流水线内容(YAML)</label>
                    <div id="yaml-editor" style="flex: 1; min-height: 400px;"></div>
                </div>
            </div>
            <div id="pipeline-view-list" style="display: none; position: fixed; inset: 81px 24px 24px; overflow: hidden; flex-direction: column;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; flex-shrink: 0;">
                    <label>执行记录</label>
                    <div class="pagination">
                        <button id="pipeline-task-prev-page" class="button-sure">上一页</button>
                        <span id="pipeline-task-page-info">第1页 / 共1页</span>
                        <button id="pipeline-task-next-page" class="button-sure">下一页</button>
                        <select id="pipeline-task-page-size" class="page-size">
                            <option value="15">15</option>
                            <option value="25">25</option>
                            <option value="35">35</option>
                        </select>
                    </div>
                </div>
                <div id="pipeline-task-list" style="flex: 1; overflow-y: auto; overflow-x: hidden;"></div>
            </div>


        `;
    }

    bindCustomEvents() {
        if (this.pipeline) {
            Utils.initializeEditor('yaml-editor', this.pipeline.content, true)
                .then(editor => this.editor = editor);
        }

        document.getElementById("pipeline-view-detail-bt")?.addEventListener("click", () => this.showDetailView());
        document.getElementById("pipeline-view-list-bt")?.addEventListener("click", () => this.showTaskView());

        // 分页事件
        document.getElementById("pipeline-task-prev-page")?.addEventListener("click", () => this.prevPage());
        document.getElementById("pipeline-task-next-page")?.addEventListener("click", () => this.nextPage());
        document.getElementById("pipeline-task-page-size")?.addEventListener("change", (e) => this.changePageSize(e));
    }

    showDetailView() {
        document.getElementById("pipeline-view-list").style.display = "none";
        document.getElementById("pipeline-view-detail").style.display = "block";
    }

    showTaskView() {
        document.getElementById("pipeline-view-detail").style.display = "none";
        document.getElementById("pipeline-view-list").style.display = "flex";
        this.fetchTasks();
    }

    updateTaskData(res) {
        if (res.data && res.data.tasks) {
            this.tasks = res.data.tasks;
            this.totalPage = res.data.page.total;
            this.currentPage = res.data.page.current;
            this.rowsPerPage = res.data.page.size;
            this.renderTasks();
            this.updatePagination();
        } else {
            this.tasks = [];
            this.totalPage = 1;
            this.renderTasks();
            this.updatePagination();
        }
    }

    renderTasks() {
        const taskContainer = document.getElementById("pipeline-task-list");
        if (!taskContainer) return;

        taskContainer.innerHTML = '';

        if (this.tasks.length === 0) {
            taskContainer.innerHTML = '<div style="text-align: center; color: var(--gray-500); padding: 40px;">暂无任务记录</div>';
            return;
        }

        this.tasks.forEach(task => {
            const taskItem = document.createElement('div');
            taskItem.className = 'task-item';
            taskItem.style.cssText = `
                display: flex;
                align-items: center;
                justify-content: space-between;
                padding: 12px 16px;
                border: 1px solid var(--gray-200);
                border-radius: var(--border-radius-sm);
                margin-bottom: 8px;
                cursor: pointer;
                transition: all 0.2s ease;
                background: white;
            `;

            const statusColor = Utils.getStatusColor(task.state || 'unknown');

            taskItem.innerHTML = `
                <div style="flex: 1;">
                    <div style="font-weight: 500; margin-bottom: 4px;">${task.taskName}</div>
                    <div style="font-size: 12px; color: var(--gray-500);">
                        开始时间: ${task.time.start || '---'}
                        结束时间: ${task.time.end || '---'}
                    </div>
                </div>
                <div class="status-indicator" style="background-color: ${statusColor};">
                    ${task.state || 'unknown'}
                </div>
            `;

            taskItem.addEventListener('mouseenter', () => {
                taskItem.style.boxShadow = 'var(--shadow-md)';
                taskItem.style.transform = 'translateY(-1px)';
            });

            taskItem.addEventListener('mouseleave', () => {
                taskItem.style.boxShadow = 'none';
                taskItem.style.transform = 'translateY(0)';
            });

            taskItem.addEventListener('click', () => {
                new TaskModal(task.taskName);
                this.close();
            });

            taskContainer.appendChild(taskItem);
        });
    }

    updatePagination() {
        const prevBtn = document.getElementById("pipeline-task-prev-page");
        const nextBtn = document.getElementById("pipeline-task-next-page");
        const pageInfo = document.getElementById("pipeline-task-page-info");

        if (prevBtn) prevBtn.disabled = this.currentPage === 1;
        if (nextBtn) nextBtn.disabled = this.currentPage === this.totalPage;
        if (pageInfo) pageInfo.textContent = `第${this.currentPage}页 / 共${this.totalPage}页`;
    }

    prevPage() {
        if (this.currentPage > 1) {
            this.currentPage--;
            this.fetchTasks();
        }
    }

    nextPage() {
        if (this.currentPage < this.totalPage) {
            this.currentPage++;
            this.fetchTasks();
        }
    }

    changePageSize(event) {
        this.rowsPerPage = parseInt(event.target.value);
        this.currentPage = 1;
        this.fetchTasks();
    }

    fetchTasks() {
        const request = {
            page: this.currentPage,
            size: this.rowsPerPage,
        };
        this.webSocketManager?.send(request);
    }

    close() {
        if (this.webSocketManager) {
            this.webSocketManager.close();
            this.webSocketManager = null;
        }
        super.close();
    }
}

// ==================== 表格管理 ====================
class TaskTable extends BaseTable {
    constructor() {
        super();
        // 添加节流控制
        this.throttledRenderRows = this.throttle((tableBody) => {
            this.renderRows(tableBody);
        }, 100);
    }

    getWebSocketUrl() {
        return ConfigManager.WS_ENDPOINTS.task;
    }

    getDataFromResponse(res) {
        return res.data?.tasks || [];
    }

    getTableBody() {
        return document.querySelector("#task-table tbody");
    }

    getColumnCount() {
        return 7;
    }

    getPrevButtonId() {
        return "task-prev-page";
    }

    getNextButtonId() {
        return "task-next-page";
    }

    getPageInfoId() {
        return "task-page-info";
    }

    getPageSizeSelectId() {
        return "task-page-size";
    }

    renderRows(tableBody) {
        const oldMap = new Map(this.previousData.map(t => [t.name, t]));
        const newMap = new Map(this.data.map(t => [t.name, t]));

        // 收集需要新增的行（准备插入到首行）
        const newRows = [];

        // 删除已移除的行
        this.previousData.forEach(task => {
            if (!newMap.has(task.name)) {
                const row = tableBody.querySelector(`tr[data-name="${task.name}"]`);
                if (row) row.remove();
            }
        });

        // 收集需要新增的行
        this.data.forEach(task => {
            if (!oldMap.has(task.name)) {
                const row = this.createRowElement(task);
                newRows.push({ row, task });
            }
        });

        // 批量插入新行到表格首行位置
        if (newRows.length > 0) {
            // 倒序插入确保顺序正确（最新的在最上面）
            newRows.reverse().forEach(({ row, task }) => {
                if (tableBody.firstChild) {
                    tableBody.insertBefore(row, tableBody.firstChild);
                } else {
                    tableBody.appendChild(row);
                }
                // 绑定事件
                this.bindRowEvents(row, task);
            });
        }

        // 批量更新已有行的变更字段
        this.batchUpdateExistingRows(tableBody, oldMap, newMap);
    }

    // 优化后的批量更新方法
    batchUpdateExistingRows(tableBody, oldMap, newMap) {
        // 使用requestAnimationFrame优化DOM更新时机
        requestAnimationFrame(() => {
            const updates = [];

            this.data.forEach(task => {
                if (!oldMap.has(task.name)) return;

                const old = oldMap.get(task.name);
                const row = tableBody.querySelector(`tr[data-name="${task.name}"]`);
                if (!row) return;

                // 如果状态发生变化，标记为需要重新创建
                if (task.state !== old.state) {
                    updates.push({ type: 'recreate', row, task });
                    return;
                }

                // 收集字段更新
                const fieldUpdates = this.getFieldUpdates(row, task, old);
                if (fieldUpdates.length > 0) {
                    updates.push({ type: 'update', row, task, fieldUpdates });
                }
            });

            // 批量执行更新，减少重排重绘
            this.executeBatchUpdates(updates, tableBody);
        });
    }

    getFieldUpdates(row, task, old) {
        const updateConfigs = [
            ['count', '.col-count'],
            ['time.start', '.col-start', (val) => val || '---'],
            ['time.end', '.col-end', (val) => val || '---'],
            ['message', '.col-message', (val, element) => {
                const full = val || "";
                const display = full.length > 150 ? full.slice(0, 150) + '...' : full;
                return { text: display, title: full };
            }]
        ];

        const updates = [];

        updateConfigs.forEach(([path, selector, formatter]) => {
            const newVal = Utils.getNestedValue(task, path);
            const oldVal = Utils.getNestedValue(old, path);

            if (newVal !== oldVal) {
                const element = row.querySelector(selector);
                if (element) {
                    updates.push({
                        element,
                        newVal,
                        formatter
                    });
                }
            }
        });

        return updates;
    }

    executeBatchUpdates(updates, tableBody) {
        updates.forEach(update => {
            if (update.type === 'recreate') {
                // 保存位置信息，在原位置重新创建
                const nextSibling = update.row.nextSibling;
                const parentNode = update.row.parentNode;

                update.row.remove();
                const newRow = this.createRowElement(update.task);

                if (nextSibling) {
                    parentNode.insertBefore(newRow, nextSibling);
                } else {
                    parentNode.appendChild(newRow);
                }

                this.bindRowEvents(newRow, update.task);
            } else if (update.type === 'update') {
                // 批量更新字段，减少重排
                update.fieldUpdates.forEach(({ element, newVal, formatter }) => {
                    if (formatter) {
                        if (typeof formatter === 'function' && formatter.length > 1) {
                            formatter(newVal, element);
                        } else {
                            const result = formatter(newVal);
                            if (typeof result === 'object' && result.text !== undefined) {
                                element.textContent = result.text;
                                if (result.title) element.title = result.title;
                            } else {
                                element.textContent = result;
                            }
                        }
                    } else {
                        element.textContent = newVal;
                    }
                });
            }
        });
    }

    // 重构createRow为createRowElement，返回元素而不直接添加到DOM
    createRowElement(task) {
        const row = document.createElement("tr");
        row.setAttribute("data-name", task.name);

        const color = Utils.getStatusColor(task.state || 'unknown');
        const fullMsg = task.message || "";
        const dispMsg = fullMsg.length > 150 ? fullMsg.slice(0, 150) + '...' : fullMsg;

        row.innerHTML = `
            <td class="col-name" style="font-weight: 500;">${task.name}</td>
            <td class="col-count">${task.count}</td>
            <td class="col-message message" title="${Utils.escapeHTML(fullMsg)}">${Utils.escapeHTML(dispMsg)}</td>
            <td class="col-start">${task.time.start || '---'}</td>
            <td class="col-end">${task.time.end || '---'}</td>
            <td class="col-state">
                <div class="status-indicator" style="background-color: ${color};">
                    ${task.state}
                </div>
            </td>
            <td class="col-actions">
                <div class="dropdown">
                    <button class="dropbtn">操作</button>
                    <div class="dropdown-content">
                        <a class="detail-task" href="#">查看详情</a>
                        <a class="dump-task" href="#">导出配置</a>
                        ${task.state === 'running' ? '<a class="kill-task" href="#">强制停止</a>' : ''}
                        <a class="delete-task" href="#">删除任务</a>
                    </div>
                </div>
            </td>
        `;

        return row;
    }

    bindRowEvents(row, task) {
        row.querySelector(".detail-task")?.addEventListener("click", (e) => {
            e.preventDefault();
            new TaskModal(task.name);
        });

        row.querySelector(".dump-task")?.addEventListener("click", (e) => {
            e.preventDefault();
            this.dumpTask(task);
        });

        const killBtn = row.querySelector(".kill-task");
        if (killBtn) {
            killBtn.addEventListener("click", (e) => {
                e.preventDefault();
                APIManager.taskAction(task.name, 'kill');
            });
        }

        row.querySelector(".delete-task")?.addEventListener("click", (e) => {
            e.preventDefault();
            APIManager.deleteResource(`${baseUrl}${ConfigManager.API_ENDPOINTS.task}/${task.name}`, task.name, '任务');
        });
    }

    async dumpTask(task) {
        try {
            const response = await APIManager.request(`${baseUrl}${ConfigManager.API_ENDPOINTS.task}/${task.name}/dump`);

            const blob = new Blob([response.data], { type: 'application/yaml' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${task.name}.yaml`;
            a.click();
            window.URL.revokeObjectURL(url);

            Utils.showToast("任务配置导出成功", 'success');
        } catch (error) {
            Utils.showToast(`任务导出失败: ${error.message}`, 'error');
        }
    }

    // 节流工具方法
    throttle(func, limit) {
        let inThrottle;
        return function() {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        }
    }

    // 如果需要在高频更新场景使用节流版本
    renderRowsThrottled(tableBody) {
        this.throttledRenderRows(tableBody);
    }

    // 新增：支持自定义插入位置
    insertRowAtPosition(task, tableBody, position = 'top') {
        const row = this.createRowElement(task);

        switch(position) {
            case 'top':
                if (tableBody.firstChild) {
                    tableBody.insertBefore(row, tableBody.firstChild);
                } else {
                    tableBody.appendChild(row);
                }
                break;
            case 'bottom':
                tableBody.appendChild(row);
                break;
            default:
                // 指定索引位置
                const children = tableBody.children;
                if (position < children.length) {
                    tableBody.insertBefore(row, children[position]);
                } else {
                    tableBody.appendChild(row);
                }
        }

        this.bindRowEvents(row, task);
        return row;
    }
}

class PipelineTable extends BaseTable {
    getWebSocketUrl() {
        return ConfigManager.WS_ENDPOINTS.pipeline;
    }

    getDataFromResponse(res) {
        return res.data?.pipelines || [];
    }

    getTableBody() {
        return document.querySelector("#pipeline-table tbody");
    }

    getColumnCount() {
        return 4;
    }

    getPrevButtonId() {
        return "pipeline-prev-page";
    }

    getNextButtonId() {
        return "pipeline-next-page";
    }

    getPageInfoId() {
        return "pipeline-page-info";
    }

    getPageSizeSelectId() {
        return "pipeline-page-size";
    }

    renderRows(tableBody) {
        tableBody.innerHTML = "";

        this.data.forEach(pipeline => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td style="font-weight: 500;">${pipeline.name}</td>
                <td>
                    <span style="background: var(--gray-100); padding: 2px 8px; border-radius: 12px; font-size: 12px;">
                        ${pipeline.tplType}
                    </span>
                </td>
                <td>${pipeline.disable ? '是' : '否'}</td>
                <td>
                    <div class="dropdown">
                        <button class="dropbtn">操作</button>
                        <div class="dropdown-content">
                            <a class="detail-pipeline" href="#">查看详情</a>
                            <a class="edit-pipeline" href="#">编辑配置</a>
                            <a class="run-pipeline" href="#">立即运行</a>
                            <a class="delete-pipeline" href="#">删除流水线</a>
                        </div>
                    </div>
                </td>
            `;

            tableBody.appendChild(row);
            this.bindRowEvents(row, pipeline);
        });
    }

    bindRowEvents(row, pipeline) {
        row.querySelector(".detail-pipeline")?.addEventListener("click", (e) => {
            e.preventDefault();
            new PipelineDetailModal(pipeline.name);
        });

        row.querySelector(".edit-pipeline")?.addEventListener("click", (e) => {
            e.preventDefault();
            this.editPipeline(pipeline.name);
        });

        row.querySelector(".run-pipeline")?.addEventListener("click", (e) => {
            e.preventDefault();
            new RunPipelineModal(pipeline.name);
        });

        row.querySelector(".delete-pipeline")?.addEventListener("click", (e) => {
            e.preventDefault();
            APIManager.deleteResource(`${baseUrl}${ConfigManager.API_ENDPOINTS.pipeline}/${pipeline.name}`, pipeline.name, '流水线');
        });
    }

    async editPipeline(pipelineName) {
        try {
            const response = await APIManager.request(`${baseUrl}${ConfigManager.API_ENDPOINTS.pipeline}/${pipelineName}`);
            new PipelineFormModal(response.data);
        } catch (error) {
            Utils.showToast(`获取流水线详情失败: ${error.message}`, 'error');
        }
    }
}

// ==================== 主应用类 ====================
class Main {
    constructor() {
        this.taskTable = null;
        this.pipelineTable = null;
        this.currentView = 'task';
        this.graphManager = new G6GraphManager();

        window.addEventListener("resize", () => this.handleResize());
        this.createMainContent();
        this.showTaskView();
        this.initTables();
        this.addEventListeners();
        new EventListener();
    }

    handleResize() {
        clearTimeout(this.resizeTimeout);
        this.resizeTimeout = setTimeout(() => {
            location.reload();
        }, 300);
    }

    createMainContent() {
        const mainDiv = document.createElement('div');
        mainDiv.id = 'main';

        mainDiv.innerHTML = `
            <div id="container" class="container">
                <div class="header">
                    <div class="menu">
                        <button id="task-list" class="button-sure">任务管理</button>
                        <button id="pipeline-list" class="button-sure">流水线管理</button>
                        <button id="file-upload-btn" class="button-sure">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style="margin-right: 4px;">
                                <path d="M9 16h6v-6h4l-7-7-7 7h4v6zm-4 2h14v2H5v-2z"/>
                            </svg>
                            文件上传
                        </button>
                    </div>
                    <div id="options">
                        <button id="add-task" class="button-sure" style="display: block">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style="margin-right: 4px;">
                                <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                            </svg>
                            新建任务
                        </button>
                        <button id="add-pipeline" class="button-sure" style="display: none">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style="margin-right: 4px;">
                                <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                            </svg>
                            新建流水线
                        </button>
                        
                        <div style="display: flex; align-items: center; gap: 12px; margin-left: 16px;">
                            <div style="width: 1px; height: 20px; background: var(--gray-300);"></div>
                            <span id="title" style="font-size: 16px; font-weight: 600; color: var(--gray-700);">任务管理</span>
                        </div>
                        
                        <div id="task-table-pagination" class="pagination" style="display: flex;">
                            <button id="task-prev-page" class="button-sure">上一页</button>
                            <span id="task-page-info">第1页 / 共1页</span>
                            <button id="task-next-page" class="button-sure">下一页</button>
                            <select id="task-page-size" class="page-size">
                                <option value="15">15条/页</option>
                                <option value="25">25条/页</option>
                                <option value="35">35条/页</option>
                                <option value="50">50条/页</option>
                            </select>
                        </div>
                        
                        <div id="pipeline-table-pagination" class="pagination" style="display: none;">
                            <button id="pipeline-prev-page" class="button-sure">上一页</button>
                            <span id="pipeline-page-info">第1页 / 共1页</span>
                            <button id="pipeline-next-page" class="button-sure">下一页</button>
                            <select id="pipeline-page-size" class="page-size">
                                <option value="15">15条/页</option>
                                <option value="25">25条/页</option>
                                <option value="35">35条/页</option>
                                <option value="50">50条/页</option>
                            </select>
                        </div>
                    </div>
                </div>
                
                <div id="task-table-container" class="table-container" style="display: block;">
                    <table id="task-table" class="common-table">
                        <thead>
                            <tr>
                                <th style="width: 200px;">任务名称</th>
                                <th style="width: 80px;">步骤数</th>
                                <th style="min-width: 200px;">执行信息</th>
                                <th style="width: 160px;">开始时间</th>
                                <th style="width: 160px;">结束时间</th>
                                <th style="width: 100px;">执行状态</th>
                                <th style="width: 120px;">操作</th>
                            </tr>
                        </thead>
                        <tbody></tbody>
                    </table>
                </div>
                
                <div id="pipeline-table-container" class="table-container" style="display: none;">
                    <table id="pipeline-table" class="common-table">
                        <thead>
                            <tr>
                                <th style="min-width: 200px;">流水线名称</th>
                                <th style="width: 120px;">模板类型</th>
                                <th style="width: 80px;">是否禁用</th>
                                <th style="width: 120px;">操作</th>
                            </tr>
                        </thead>
                        <tbody></tbody>
                    </table>
                </div>
            </div>
        `;

        document.body.appendChild(mainDiv);
    }

    initTables() {
        this.taskTable = new TaskTable();
        this.pipelineTable = new PipelineTable();
    }

    addEventListeners() {
        document.getElementById("task-list")?.addEventListener("click", () => this.showTaskView());
        document.getElementById("pipeline-list")?.addEventListener("click", () => this.showPipelineView());
        document.getElementById("add-task")?.addEventListener("click", () => new TaskFormModal());
        document.getElementById("add-pipeline")?.addEventListener("click", () => new PipelineFormModal());
        document.getElementById("file-upload-btn")?.addEventListener("click", () => new FileUploadModal());
    }

    showTaskView() {
        this.currentView = 'task';
        this.updateViewVisibility();
        this.updateNavigationButtons();
    }

    showPipelineView() {
        this.currentView = 'pipeline';
        this.updateViewVisibility();
        this.updateNavigationButtons();
    }

    updateViewVisibility() {
        const isTaskView = this.currentView === 'task';

        document.getElementById("title").textContent = isTaskView ? "任务管理" : "流水线管理";

        document.getElementById("task-table-container").style.display = isTaskView ? "block" : "none";
        document.getElementById("add-task").style.display = isTaskView ? "block" : "none";
        document.getElementById("task-table-pagination").style.display = isTaskView ? "flex" : "none";

        document.getElementById("pipeline-table-container").style.display = isTaskView ? "none" : "block";
        document.getElementById("add-pipeline").style.display = isTaskView ? "none" : "block";
        document.getElementById("pipeline-table-pagination").style.display = isTaskView ? "none" : "flex";
    }

    updateNavigationButtons() {
        const taskBtn = document.getElementById("task-list");
        const pipelineBtn = document.getElementById("pipeline-list");

        taskBtn.style.background = this.currentView === 'task' ? 'var(--primary-color)' : 'var(--gray-200)';
        taskBtn.style.color = this.currentView === 'task' ? 'white' : 'var(--gray-700)';

        pipelineBtn.style.background = this.currentView === 'pipeline' ? 'var(--primary-color)' : 'var(--gray-200)';
        pipelineBtn.style.color = this.currentView === 'pipeline' ? 'white' : 'var(--gray-700)';
    }
}

// ==================== 事件监听器 ====================
class EventListener {
    constructor() {
        this.initEventContainer();
        this.listenForEvents();
    }

    initEventContainer() {
        const eventContainer = document.createElement('div');
        eventContainer.id = 'event-container';
        eventContainer.className = 'event-container';
        document.body.appendChild(eventContainer);
    }

    listenForEvents() {
        const eventSource = new EventSource(`${baseUrl}${ConfigManager.API_ENDPOINTS.event}`);

        eventSource.onmessage = (event) => {
            this.displayEventMessage(event.data);
        };

        eventSource.onerror = (event) => {
            console.error('SSE连接错误:', event);
        };

        eventSource.onopen = () => {
            console.log('SSE连接已建立');
        };
    }

    displayEventMessage(message) {
        const eventContainer = document.getElementById('event-container');
        if (!eventContainer) return;

        const messageElement = document.createElement('p');
        messageElement.textContent = message;

        eventContainer.appendChild(messageElement);

        // 添加到容器顶部（最新消息在上方）
        eventContainer.insertBefore(messageElement, eventContainer.firstChild);

        // 保留最近3条消息
        if (eventContainer.children.length > 3) {
            eventContainer.removeChild(eventContainer.firstChild);
        }

        // 自动移除消息
        setTimeout(() => {
            if (messageElement.parentNode) {
                messageElement.classList.add('event-fade-out');
                setTimeout(() => {
                    if (messageElement.parentNode) {
                        messageElement.remove();
                    }
                }, 300);
            }
        }, 1000);
    }
}

// ==================== 文件上传管理器 ====================
class FileUploadManager {
    constructor() {
        this.uploads = new Map();
        this.nextFileId = 1;
        this.activeUploads = new Set();
        this.isGlobalPaused = false;
        this.endpoint = `${baseUrl}/api/v1/files/`;
    }

    async initializeElements() {
        this.fileInput = document.querySelector('#file-upload-input');
        this.selectFileBtn = document.querySelector('#select-file-btn');
        this.dropZone = document.querySelector('#drop-zone');
        this.uploadQueue = document.querySelector('#upload-queue');
        this.fileList = document.querySelector('#file-list');
        this.uploadHistory = document.querySelector('#upload-history');

        // 按钮
        this.startAllBtn = document.querySelector('#start-all-btn');
        this.pauseAllBtn = document.querySelector('#pause-all-btn');
        this.resumeAllBtn = document.querySelector('#resume-all-btn');
        this.clearAllBtn = document.querySelector('#clear-all-btn');

        // 配置
        this.chunkInput = document.querySelector('#chunk-size');
        this.taskidInput = document.querySelector('#task-id');
        this.parallelInput = document.querySelector('#parallel-uploads');
        this.queueConcurrencyInput = document.querySelector('#queue-concurrency');

        // 统计
        this.totalCountEl = document.querySelector('#total-count');
        this.waitingCountEl = document.querySelector('#waiting-count');
        this.uploadingCountEl = document.querySelector('#uploading-count');
        this.pausedCountEl = document.querySelector('#paused-count');
        this.errorCountEl = document.querySelector('#error-count');

        this.initEvents();
    }


    initEvents() {
        if (!this.selectFileBtn) return;

        this.selectFileBtn.addEventListener('click', () => this.fileInput.click());
        this.fileInput.addEventListener('change', (e) => this.handleFileSelect(e.target.files));

        // 拖拽事件
        this.dropZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            this.dropZone.classList.add('dragover');
        });

        this.dropZone.addEventListener('dragleave', (e) => {
            e.preventDefault();
            this.dropZone.classList.remove('dragover');
        });

        this.dropZone.addEventListener('drop', (e) => {
            e.preventDefault();
            this.dropZone.classList.remove('dragover');
            this.handleFileSelect(e.dataTransfer.files);
        });

        // 按钮事件
        this.startAllBtn?.addEventListener('click', () => this.startAllUploads());
        this.pauseAllBtn?.addEventListener('click', () => this.pauseAllUploads());
        this.resumeAllBtn?.addEventListener('click', () => this.resumeAllUploads());
        this.clearAllBtn?.addEventListener('click', () => this.clearAllFiles());
        this.queueConcurrencyInput?.addEventListener('change', () => this.processQueue());
    }

    handleFileSelect(files) {
        if (!files || files.length === 0) return;

        Array.from(files).forEach(file => {
            this.addFileToQueue(file);
        });

        this.uploadQueue.classList.remove('d-none');
        this.fileInput.value = '';
        this.updateStats();
    }

    addFileToQueue(file) {
        const fileId = this.nextFileId++;
        const uploadInfo = {
            id: fileId,
            file: file,
            upload: null,
            status: 'waiting',
            progress: 0,
            element: null,
            startTime: null,
            pausedByUser: false
        };

        this.uploads.set(fileId, uploadInfo);
        this.createFileElement(uploadInfo);
        this.updateStats();
    }

    createFileElement(uploadInfo) {
        const fileElement = document.createElement('div');
        fileElement.className = 'file-item waiting';
        fileElement.innerHTML = `
            <div class="file-header">
                <div class="file-icon">📄</div>
                <div class="file-info">
                    <h4 class="file-name">${uploadInfo.file.name}</h4>
                    <div class="file-meta">
                        <span>📦 ${this.formatFileSize(uploadInfo.file.size)}</span>
                        <span>🔢 ID: ${uploadInfo.id}</span>
                    </div>
                </div>
            </div>

            <div class="file-status">
                <span class="status-badge status-waiting">排队中</span>
            </div>

            <div class="progress-container d-none">
                <div class="progress-header">
                    <span class="progress-text">0%</span>
                    <span class="progress-speed"></span>
                </div>
                <div class="progress-bar-container">
                    <div class="progress-bar" style="width: 0%"></div>
                </div>
            </div>

            <div class="file-actions">
                <button class="btn btn-sm button-sure start-btn" data-id="${uploadInfo.id}">开始</button>
                <button class="btn btn-sm button-cancel pause-btn d-none" data-id="${uploadInfo.id}">暂停</button>
                <button class="btn btn-sm button-sure resume-btn d-none" data-id="${uploadInfo.id}">恢复</button>
                <button class="btn btn-sm button-cancel remove-btn" data-id="${uploadInfo.id}">移除</button>
            </div>
        `;

        uploadInfo.element = fileElement;
        this.fileList.appendChild(fileElement);

        // 绑定事件
        fileElement.querySelector('.start-btn')?.addEventListener('click', () => this.startUpload(uploadInfo.id));
        fileElement.querySelector('.pause-btn')?.addEventListener('click', () => this.pauseUpload(uploadInfo.id));
        fileElement.querySelector('.resume-btn')?.addEventListener('click', () => this.resumeUpload(uploadInfo.id));
        fileElement.querySelector('.remove-btn')?.addEventListener('click', () => this.removeFile(uploadInfo.id));
    }

    getQueueConcurrency() {
        const value = parseInt(this.queueConcurrencyInput?.value, 10);
        return isNaN(value) ? 3 : Math.max(1, Math.min(10, value));
    }

    processQueue() {
        if (this.isGlobalPaused) return;

        const maxConcurrency = this.getQueueConcurrency();
        const currentUploading = this.activeUploads.size;

        if (currentUploading >= maxConcurrency) return;

        const waitingFiles = Array.from(this.uploads.values())
            .filter(uploadInfo => uploadInfo.status === 'waiting')
            .sort((a, b) => a.id - b.id);

        const slotsAvailable = maxConcurrency - currentUploading;
        const filesToStart = waitingFiles.slice(0, slotsAvailable);

        filesToStart.forEach(uploadInfo => {
            this.startUploadInternal(uploadInfo.id);
        });
    }

    startUpload(fileId) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo) return;

        if (uploadInfo.status === 'paused') {
            uploadInfo.pausedByUser = false;
            this.resumeUpload(fileId);
            return;
        }

        if (uploadInfo.status !== 'waiting') return;

        if (this.activeUploads.size < this.getQueueConcurrency() && !this.isGlobalPaused) {
            this.startUploadInternal(fileId);
        }
    }

    async startUploadInternal(fileId) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo || uploadInfo.status === 'uploading') return;

        // 确保tus已加载
        if (!window.tus) {
            Utils.showToast('上传库未加载，请稍后重试', 'error');
            this.updateFileStatus(fileId, 'error');
            return;
        }

        this.activeUploads.add(fileId);
        this.updateFileStatus(fileId, 'uploading');

        let chunkSize = Number.parseInt(this.chunkInput?.value, 10);
        if (Number.isNaN(chunkSize)) chunkSize = 5242880;

        let parallelUploads = Number.parseInt(this.parallelInput?.value, 10);
        if (Number.isNaN(parallelUploads)) parallelUploads = 6;

        const options = {
            endpoint: this.endpoint,
            chunkSize: chunkSize,
            addRequestId: true,
            uploadDataDuringCreation: true,
            removeFingerprintOnSuccess: true,
            retryDelays: [0, 1000, 3000, 5000],
            parallelUploads: parallelUploads,
            metadata: {
                filename: uploadInfo.file.name,
                filetype: uploadInfo.file.type,
                task_id: this.taskidInput?.value.toString() || 'default-task',
            },
            metadataForPartialUploads: {
                task_id: this.taskidInput?.value.toString() || 'default-task',
            },
            onError: (error) => this.handleUploadError(fileId, error),
            onProgress: (bytesUploaded, bytesTotal) => this.handleUploadProgress(fileId, bytesUploaded, bytesTotal),
            onSuccess: () => this.handleUploadSuccess(fileId),
        };

        uploadInfo.upload = new window.tus.Upload(uploadInfo.file, options);
        uploadInfo.startTime = Date.now();

        try {
            const previousUploads = await uploadInfo.upload.findPreviousUploads();
            if (previousUploads.length > 0) {
                uploadInfo.upload.resumeFromPreviousUpload(previousUploads[0]);
            }
            uploadInfo.upload.start();
        } catch (error) {
            this.handleUploadError(fileId, error);
        }
    }

    pauseUpload(fileId, isUserAction = true) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo || !uploadInfo.upload) return;

        uploadInfo.pausedByUser = isUserAction;
        uploadInfo.upload.abort();
        this.activeUploads.delete(fileId);
        this.updateFileStatus(fileId, 'paused');
    }

    resumeUpload(fileId) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo || !uploadInfo.upload || uploadInfo.status !== 'paused') return;

        uploadInfo.pausedByUser = false;

        if (this.activeUploads.size >= this.getQueueConcurrency()) {
            Utils.showToast('当前上传数量已达上限', 'warning');
            return;
        }

        this.resumeUploadInternal(fileId);
    }

    resumeUploadInternal(fileId) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo || !uploadInfo.upload || uploadInfo.status !== 'paused') return;

        if (this.activeUploads.size >= this.getQueueConcurrency()) {
            uploadInfo.pausedByUser = false;
            this.updateFileStatus(fileId, 'waiting');
            return;
        }

        this.activeUploads.add(fileId);
        uploadInfo.upload.start();
        this.updateFileStatus(fileId, 'uploading');
    }

    removeFile(fileId) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo) return;

        if (uploadInfo.upload) {
            uploadInfo.upload.abort();
        }

        this.activeUploads.delete(fileId);
        uploadInfo.element.remove();
        this.uploads.delete(fileId);

        if (this.uploads.size === 0) {
            this.uploadQueue.classList.add('d-none');
        }

        this.updateStats();
        this.processQueue();
    }

    updateFileStatus(fileId, status) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo) return;

        uploadInfo.status = status;

        const element = uploadInfo.element;
        const statusBadge = element.querySelector('.status-badge');
        const startBtn = element.querySelector('.start-btn');
        const pauseBtn = element.querySelector('.pause-btn');
        const resumeBtn = element.querySelector('.resume-btn');
        const progressContainer = element.querySelector('.progress-container');

        element.className = `file-item ${status}`;
        statusBadge.className = `status-badge status-${status}`;

        switch (status) {
            case 'waiting':
                statusBadge.textContent = '排队中';
                startBtn.classList.remove('d-none');
                pauseBtn.classList.add('d-none');
                resumeBtn.classList.add('d-none');
                progressContainer.classList.add('d-none');
                break;

            case 'uploading':
                statusBadge.textContent = '上传中';
                startBtn.classList.add('d-none');
                pauseBtn.classList.remove('d-none');
                resumeBtn.classList.add('d-none');
                progressContainer.classList.remove('d-none');
                break;

            case 'paused':
                statusBadge.textContent = '已暂停';
                startBtn.classList.add('d-none');
                pauseBtn.classList.add('d-none');
                resumeBtn.classList.remove('d-none');
                break;

            case 'error':
                statusBadge.textContent = '失败';
                startBtn.classList.remove('d-none');
                pauseBtn.classList.add('d-none');
                resumeBtn.classList.add('d-none');
                break;
        }

        this.updateStats();
    }

    updateStats() {
        const statusCounts = {
            total: this.uploads.size,
            waiting: 0,
            uploading: 0,
            paused: 0,
            error: 0
        };

        this.uploads.forEach(uploadInfo => {
            statusCounts[uploadInfo.status]++;
        });

        if (this.totalCountEl) this.totalCountEl.textContent = statusCounts.total;
        if (this.waitingCountEl) this.waitingCountEl.textContent = statusCounts.waiting;
        if (this.uploadingCountEl) this.uploadingCountEl.textContent = statusCounts.uploading;
        if (this.pausedCountEl) this.pausedCountEl.textContent = statusCounts.paused;
        if (this.errorCountEl) this.errorCountEl.textContent = statusCounts.error;
    }

    handleUploadProgress(fileId, bytesUploaded, bytesTotal) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo) return;

        const percentage = ((bytesUploaded / bytesTotal) * 100).toFixed(1);
        uploadInfo.progress = percentage;

        const progressBar = uploadInfo.element.querySelector('.progress-bar');
        const progressText = uploadInfo.element.querySelector('.progress-text');
        const speedElement = uploadInfo.element.querySelector('.progress-speed');

        progressBar.style.width = `${percentage}%`;
        progressText.textContent = `${percentage}%`;

        if (uploadInfo.startTime) {
            const elapsed = (Date.now() - uploadInfo.startTime) / 1000;
            const speed = bytesUploaded / elapsed;
            speedElement.textContent = this.formatFileSize(speed) + '/s';
        }
    }

    handleUploadSuccess(fileId) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo) return;

        this.activeUploads.delete(fileId);
        this.addToHistory(uploadInfo);
        this.removeFile(fileId);
        this.processQueue();
        Utils.showToast('文件上传成功', 'success');
    }

    handleUploadError(fileId, error) {
        const uploadInfo = this.uploads.get(fileId);
        if (!uploadInfo) return;

        this.activeUploads.delete(fileId);
        this.updateFileStatus(fileId, 'error');

        const progressBar = uploadInfo.element.querySelector('.progress-bar');
        progressBar.classList.add('error');

        console.error(`文件 ${uploadInfo.file.name} 上传失败:`, error);
        Utils.showToast('文件上传失败', 'error');
        this.processQueue();
    }

    addToHistory(uploadInfo) {
        const historyItem = document.createElement('div');
        historyItem.className = 'history-item';
        historyItem.innerHTML = `
            <div class="history-header">
                <div class="success-icon">✓</div>
                <div class="history-info">
                    <h4 class="history-name">${uploadInfo.file.name}</h4>
                </div>
            </div>

            <div class="history-meta">
                <div class="history-meta-item">
                    <span>📦 ${this.formatFileSize(uploadInfo.file.size)}</span>
                </div>
                <div class="history-meta-item">
                    <span>🕒 ${new Date().toLocaleString()}</span>
                </div>
            </div>

            <div class="history-actions">
                <a href="${uploadInfo.upload.url}" class="btn btn-sm button-sure" target="_blank">下载</a>
                <button class="btn btn-sm button-cancel copy-link-btn">复制链接</button>
            </div>
        `;

        // 清除空状态
        const emptyState = this.uploadHistory.querySelector('.empty-state');
        if (emptyState) {
            emptyState.remove();
        }

        this.uploadHistory.insertBefore(historyItem, this.uploadHistory.firstChild);

        // 绑定复制链接事件
        historyItem.querySelector('.copy-link-btn')?.addEventListener('click', () => {
            this.copyToClipboard(uploadInfo.upload.url);
        });
    }

    startAllUploads() {
        this.isGlobalPaused = false;

        this.uploads.forEach((uploadInfo, fileId) => {
            if (uploadInfo.status === 'paused') {
                uploadInfo.pausedByUser = false;
                this.updateFileStatus(fileId, 'waiting');
            }
        });

        this.processQueue();
    }

    pauseAllUploads() {
        this.isGlobalPaused = true;

        this.uploads.forEach((uploadInfo, fileId) => {
            if (uploadInfo.status === 'uploading') {
                this.pauseUpload(fileId, false);
            }
        });
    }

    resumeAllUploads() {
        this.isGlobalPaused = false;

        const pausedFiles = Array.from(this.uploads.values())
            .filter(uploadInfo => uploadInfo.status === 'paused' && !uploadInfo.pausedByUser);

        if (pausedFiles.length === 0) {
            this.processQueue();
            return;
        }

        pausedFiles.forEach(uploadInfo => {
            this.resumeUploadInternal(uploadInfo.id);
        });

        setTimeout(() => {
            this.processQueue();
        }, 100);
    }

    clearAllFiles() {
        if (!confirm('确定要清空所有文件吗？')) return;

        const uploadIds = Array.from(this.uploads.keys());
        uploadIds.forEach(fileId => {
            this.removeFile(fileId);
        });
        this.isGlobalPaused = false;
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    }

    copyToClipboard(text) {
        if (navigator.clipboard) {
            navigator.clipboard.writeText(text).then(() => {
                Utils.showToast('链接已复制到剪贴板', 'success');
            }).catch(() => {
                this.fallbackCopyTextToClipboard(text);
            });
        } else {
            this.fallbackCopyTextToClipboard(text);
        }
    }

    fallbackCopyTextToClipboard(text) {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-9999px';
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        try {
            document.execCommand('copy');
            Utils.showToast('链接已复制到剪贴板', 'success');
        } catch (err) {
            Utils.showToast('复制失败，请手动复制链接', 'error');
        }
        document.body.removeChild(textArea);
    }
}

// ==================== 文件上传模态框 ====================
class FileUploadModal extends BaseModal {
    constructor() {
        super({
            id: 'file-upload',
            title: '文件上传中心',
            size: 'large'
        });

        this.uploadManager = new FileUploadManager();
        
        // 检查 tus 是否可用
        if (!window.tus) {
            Utils.showToast('上传库未加载，请刷新页面重试', 'error');
            console.error('tus.js 未加载');
            return;
        }
        
        this.create();
    }

    getBodyContent() {
        return `
            <div class="card-body" style="position: fixed; inset: 81px 24px 24px; display: flex; flex-direction: column; gap: 20px;">
                <div style="background: white; border-radius: var(--border-radius); padding: 20px; flex-shrink: 0;">
                    <div class="upload-layout">
                        <div class="config-section">
                            <div class="section-title">上传配置</div>
                            <div class="config-content">
                                <div class="form-group">
                                    <label for="task-id" class="form-label">Task ID</label>
                                    <input type="text" class="form-control" id="task-id" value="task-123456">
                                    <div class="form-help">任务ID: 用于上传到指定流水线下</div>
                                </div>
                                <div class="form-group">
                                    <label for="chunk-size" class="form-label">分块大小</label>
                                    <input type="number" class="form-control" id="chunk-size" value="5242880">
                                    <div class="form-help">默认: 5MB</div>
                                </div>
                                <div class="form-group">
                                    <label for="parallel-uploads" class="form-label">单文件并行数</label>
                                    <input type="number" class="form-control" id="parallel-uploads" value="6">
                                    <div class="form-help">分块并行数</div>
                                </div>
                                <div class="form-group">
                                    <label for="queue-concurrency" class="form-label">队列并发数</label>
                                    <input type="number" class="form-control" id="queue-concurrency" value="3" min="1" max="10">
                                    <div class="form-help">同时上传文件数</div>
                                </div>
                            </div>
                        </div>

                        <div class="upload-section">
                            <div class="upload-zone" id="drop-zone">
                                <div class="upload-icon">📁</div>
                                <h3 class="upload-title">拖放文件到这里或点击选择</h3>
                                <p class="upload-subtitle">支持多文件同时上传，任意格式</p>
                                <button type="button" class="btn button-sure" id="select-file-btn">选择文件</button>
                            </div>
                        </div>
                    </div>

                    <input type="file" class="d-none" id="file-upload-input" multiple>

                    <div id="upload-queue" class="d-none" style="margin-top: 20px;">
                        <div class="stats" id="stats-grid">
                            <div class="stat-item">
                                <div class="stat-label">总计</div>
                                <div class="stat-value" id="total-count">0</div>
                            </div>
                            <div class="stat-item stat-waiting">
                                <div class="stat-label">排队中</div>
                                <div class="stat-value" id="waiting-count">0</div>
                            </div>
                            <div class="stat-item stat-uploading">
                                <div class="stat-label">上传中</div>
                                <div class="stat-value" id="uploading-count">0</div>
                            </div>
                            <div class="stat-item stat-paused">
                                <div class="stat-label">已暂停</div>
                                <div class="stat-value" id="paused-count">0</div>
                            </div>
                            <div class="stat-item stat-error">
                                <div class="stat-label">失败</div>
                                <div class="stat-value" id="error-count">0</div>
                            </div>
                        </div>

                        <div class="action-bar">
                            <button class="btn button-sure" id="start-all-btn">开始全部</button>
                            <button class="btn button-cancel" id="pause-all-btn">暂停全部</button>
                            <button class="btn button-sure" id="resume-all-btn">恢复全部</button>
                            <button class="btn button-cancel" id="clear-all-btn">清空队列</button>
                        </div>

                        <div class="file-list" id="file-list"></div>
                    </div>
                </div>

                <div style="background: white; border-radius: var(--border-radius); padding: 20px; flex: 1; overflow: hidden; display: flex; flex-direction: column;">
                    <h6 style="margin: 0 0 16px 0; color: var(--gray-700);">上传记录</h6>
                    <div class="history-list" id="upload-history" style="flex: 1; overflow-y: auto;">
                        <div class="empty-state">
                            <div class="empty-icon">📦</div>
                            <div class="empty-title">暂无上传记录</div>
                            <p class="empty-text">上传完成的文件会显示在这里</p>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    bindCustomEvents() {
        this.uploadManager.initializeElements();
    }
}

// 初始化高层级z-index
window.highestZIndex = 1000;

