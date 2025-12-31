/**
 * 应用配置管理
 * 统一管理 API 端点、状态映射、模板等配置
 */

// 导入设计令牌
import { STATUS_COLORS, PAGINATION } from '@/constants/design-tokens'

// 导入动态API配置
import { apiConfig, API_ENDPOINTS as API_ENDPOINTS_CONFIG, getApiUrl, getWsUrl } from './api'

// API 配置（兼容旧代码）
export const API_BASE_URL = apiConfig.baseUrl
export const WS_BASE_URL = apiConfig.wsBaseUrl

// API 端点（兼容旧代码，使用动态配置）
export const API_ENDPOINTS = {
  task: apiConfig.apiEndpoints.task,
  pipeline: apiConfig.apiEndpoints.pipeline,
  event: apiConfig.apiEndpoints.event,
  files: '/api/v1/files', // 文件上传端点
}

// 导出新的API配置工具函数
export { getApiUrl, getWsUrl, API_ENDPOINTS_CONFIG }

// 重新导出状态颜色（保持向后兼容）
export { STATUS_COLORS }

// 分页默认配置（保持向后兼容）
export const PAGINATION_DEFAULTS = {
  page: PAGINATION.defaultPage,
  size: PAGINATION.defaultSize,
  sizes: PAGINATION.pageSizes,
}

export const TASK_YAML_TEMPLATE = `#name: 测试测试任务
desc: 这是一段任务描述
kind: dag
timeout: 2m
env:
  - name: GLOBAL_NAME
    value: "全局变量"
step:
  - name: shell0-0
    desc: 执行shell脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    type: sh
    content: |-
      ping -c 4 1.1.1.1
  - name: shell0-1
    desc: 执行shell脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    depends:
      - shell0-0
    type: sh
    content: |-
      ping -c 4 1.1.1.1
  - name: python0-0
    desc: 执行python脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    type: py3
    content: |-
      import subprocess
      command = ["ping", "-c", "4", "1.1.1.1"]
      try:
          result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
          print("Ping 命令的输出：")
          print(result.stdout)
      except subprocess.CalledProcessError as e:
          print("执行 ping 命令时发生错误：")
          print(e.stderr)
  - name: python0-1
    desc: 执行python脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    depends:
      - python0-0
    type: py3
    content: |-
      import subprocess
      command = ["ping", "-c", "4", "1.1.1.1"]
      try:
          result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
          print("Ping 命令的输出：")
          print(result.stdout)
      except subprocess.CalledProcessError as e:
          print("执行 ping 命令时发生错误：")
          print(e.stderr)
  - name: shell
    desc: 执行shell脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    type: sh
    content: |-
      ping -c 4 1.1.1.1
  - name: python
    desc: 执行python脚本
    timeout: 2m
    env:
      - name: Test
        value: "test_env"
    depends:
      - shell
    type: py3
    content: |-
      import subprocess
      command = ["ping", "-c", "4", "1.1.1.1"]
      try:
          result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
          print("Ping 命令的输出：")
          print(result.stdout)
      except subprocess.CalledProcessError as e:
          print("执行 ping 命令时发生错误：")
          print(e.stderr)
  - name: yaegi
    desc: 执行yaegi脚本
    env:
      - name: Test
        value: "test_env"
    depends:
      - python
    type: yaegi
    content: |-
      import (
        "context"
        "fmt"
        "os/exec"
        
        "github.com/tidwall/gjson"
      )
      func EvalCall(ctx context.Context, params gjson.Result) {
        fmt.Println(params)
        cmd := exec.Command("ping", "-c", "4", "1.1.1.1")
        output, err := cmd.CombinedOutput()
        if err != nil {
          fmt.Println("执行 ping 命令时发生错误：", err)
          return
        }
        fmt.Println("Ping 命令的输出：")
        fmt.Println(string(output))
      }
  - name: 聚合测试
    desc: 等待所有脚本执行完成
    env:
      - name: Test
        value: "test_env"
    depends:
      - yaegi
      - 多分支执行2
    type: sh
    content: |-
      echo "done done"
  - name: 多分支执行
    desc: 测试多分支执行
    env:
      - name: Test
        value: "test_env"
    type: yaegi
    content: |-
      import (
        "context"
        "fmt"
        "io"
        "log"
        "net/http"
        
        "github.com/tidwall/gjson"
      )
      func EvalCall(ctx context.Context, params gjson.Result) {
        resp, err := http.Get("https://www.baidu.com")
        if err != nil {
          log.Fatalf("HTTP 请求失败: %v", err)
          return
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
          log.Printf("HTTP 请求失败，状态码: %d", resp.StatusCode)
          return
        }
        // 读取响应体
        body, err := io.ReadAll(resp.Body)
        if err != nil {
        \tlog.Fatalf("读取响应体失败: %v", err)
        \treturn
        }
        
        // 打印响应内容
        fmt.Println("HTTP 响应内容:")
        fmt.Println(string(body))
      }
  - name: 多分支执行1
    desc: 测试多分支执行
    env:
      - name: Test
        value: "test_env"
    depends:
      - 多分支执行
      - shell0-1
      - python0-1
    type: yaegi
    content: |-
      import (
        "context"
        "fmt"
        "io"
        "log"
        "net/http"
        
        "github.com/tidwall/gjson"
      )
      func EvalCall(ctx context.Context, params gjson.Result) {
        resp, err := http.Get("https://www.baidu.com")
        if err != nil {
          log.Fatalf("HTTP 请求失败: %v", err)
          return
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
          log.Printf("HTTP 请求失败，状态码: %d", resp.StatusCode)
          return
        }
        // 读取响应体
        body, err := io.ReadAll(resp.Body)
        if err != nil {
        \tlog.Fatalf("读取响应体失败: %v", err)
        \treturn
        }
        
        // 打印响应内容
        fmt.Println("HTTP 响应内容:")
        fmt.Println(string(body))
      }
  - name: 多分支执行2
    desc: 测试多分支执行
    env:
      - name: Test
        value: "test_env"
    depends:
      - 多分支执行1
    type: yaegi
    content: |-
      import (
        "context"
        "fmt"
        "io"
        "log"
        "net/http"
        
        "github.com/tidwall/gjson"
      )
      func EvalCall(ctx context.Context, params gjson.Result) {
        resp, err := http.Get("https://www.baidu.com")
        if err != nil {
          log.Fatalf("HTTP 请求失败: %v", err)
          return
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
          log.Printf("HTTP 请求失败，状态码: %d", resp.StatusCode)
          return
        }
        // 读取响应体
        body, err := io.ReadAll(resp.Body)
        if err != nil {
        \tlog.Fatalf("读取响应体失败: %v", err)
        \treturn
        }
        
        // 打印响应内容
        fmt.Println("HTTP 响应内容:")
        fmt.Println(string(body))
      }
            `

export const PIPELINE_PARAMS_TEMPLATE = `params:
  # 在这里添加流水线运行参数
  imageTag: latest`