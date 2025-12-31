import{s as p,m as d,w as m}from"./element-plus-D2VUXABr.js";const g=()=>{if(!window.APP_CONFIG)throw console.error("❌ APP_CONFIG未初始化！请检查index.html中的配置脚本"),new Error("APP_CONFIG未初始化");const{baseUrl:s,wsBaseUrl:t,api:e}=window.APP_CONFIG;return{baseUrl:s,wsBaseUrl:t,apiEndpoints:{task:e.task,pipeline:e.pipeline,event:e.event}}},a=g(),f=a.baseUrl,h=a.wsBaseUrl,T={task:a.apiEndpoints.task,pipeline:a.apiEndpoints.pipeline,event:a.apiEndpoints.event,files:"/api/v1/files"},y=`#name: 测试测试任务
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
        	log.Fatalf("读取响应体失败: %v", err)
        	return
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
        	log.Fatalf("读取响应体失败: %v", err)
        	return
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
        	log.Fatalf("读取响应体失败: %v", err)
        	return
        }
        
        // 打印响应内容
        fmt.Println("HTTP 响应内容:")
        fmt.Println(string(body))
      }
            `,b=`params:
  # 在这里添加流水线运行参数
  imageTag: latest`;function w(s="grid",t=null){const e=t?localStorage.getItem(t):null,n=p(e||s),o=u=>{n.value=u},l=()=>{o("grid")},i=()=>{o("table")},r=()=>{n.value=n.value==="grid"?"table":"grid"},c=d(()=>n.value==="grid"),v=d(()=>n.value==="table");return t&&m(n,u=>{localStorage.setItem(t,u)}),{viewMode:n,setViewMode:o,toGridView:l,toTableView:i,toggleViewMode:r,isGridView:c,isTableView:v}}function x(s=!1,t="sidebar-collapsed"){const e=localStorage.getItem(t),n=p(e?e==="true":s),o=()=>{n.value=!n.value},l=()=>{n.value=!1},i=()=>{n.value=!0},r=d(()=>n.value?"48px":"160px");return m(n,c=>{localStorage.setItem(t,c.toString())}),{isCollapsed:n,toggle:o,expand:l,collapse:i,width:r}}function _(s="default"){const t=p(s),e=p([s]);return{activeTab:t,tabHistory:e,switchTab:r=>{t.value!==r&&(t.value=r,e.value.push(r),e.value.length>10&&e.value.shift())},goBack:()=>{if(e.value.length>1){e.value.pop();const r=e.value[e.value.length-1];t.value=r}},isActive:r=>t.value===r,reset:()=>{t.value=s,e.value=[s]}}}const C=(s,t)=>{const e=s.__vccOpts||s;for(const[n,o]of t)e[n]=o;return e};export{f as A,b as P,y as T,h as W,C as _,T as a,w as b,_ as c,x as u};
