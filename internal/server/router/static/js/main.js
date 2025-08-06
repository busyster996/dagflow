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
      ping -c 4 1.1.1.1`,
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

// 初始化高层级z-index
window.highestZIndex = 1000;

