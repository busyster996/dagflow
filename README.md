# dagflow

[![Go Report Card](https://goreportcard.com/badge/github.com/busyster996/dagflow)](https://goreportcard.com/report/github.com/busyster996/dagflow)
[![Go Reference](https://pkg.go.dev/badge/github.com/busyster996/dagflow.svg)](https://pkg.go.dev/github.com/busyster996/dagflow)

An `API` for cross-platform custom orchestration of execution steps without any third-party dependencies.
Based on `DAG` , it implements the scheduling function of sequential execution of dependent steps and concurrent execution of non-dependent steps.

It provides `API` remote operation mode, batch execution of `Shell` , `Powershell` , `Python` and other commands,
and easily completes common management tasks such as running automated operation and maintenance scripts, polling processes, installing or uninstalling software, updating applications, and installing patches.

## Operating system remote execution interface

![](images/dag.png)

## Feature

- [x] support `Windows` / `Linux` / `Mac`
- [x] Dynamically adjust the amount workers
- [x] Orchestrating execution based on directed acyclic graph ( `DAG` )
- [x] Supports forced termination of tasks or steps
- [x] Supports suspension and resumption of tasks or steps
- [x] Support timeout for tasks or steps
- [x] Task-level Workspace isolation
- [x] Browse, upload and download tasks in Workspace
- [x] Self-update, use parameter `--self_url`
- [x] WebShell
- [x] Step retry policy mechanism
- [x] Support delayed Task
- [x] Send events before/after a task or step is executed
- [ ] Task or step plugin implementation

## Help
```text
Usage:
  dagflow_linux_amd64_v1 [command]

Available Commands:
  api         start api service
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  standalone  start standalone service
  version     print version information and quit
  worker      start a worker service

Flags:
      --db_url string        database type. [sqlite,mysql,postgres,sqlserver] (default "sqlite://localhost")
      --enable_self_update   enable self update (default true)
      --help                 Print usage
      --kind_id int        data kind id (default 1)
      --log_level string     log level [debug,info,warn,error] (default "debug")
      --log_output string    log output [file,stdout] (default "file")
      --mq_url string        message queue url. [inmemory,amqp] (default "inmemory://localhost")
      --node_id int          node id (default 1)
      --root_dir string      root directory (default "/usr/local/dagflow")
      --self_url string      self Update URL (default "https://oss.yfdou.com/tools/dagflow")
  -v, --version              Print version information and quit

Use "dagflow_linux_amd64_v1 [command] --help" for more information about a command
```

## How to use
### Windows
Open PowerShell in management mode to add services
```powershell
New-Service -Name dagflow -BinaryPathName "C:\dagflow\bin\dagflow_windows_amd64_v1.exe standalone" -DisplayName  "dagflow " -StartupType Automatic
sc.exe failure dagflow reset= 0 actions= restart/0/restart/0/restart/0
sc.exe start dagflow
```

### Linux
```shell
echo > /etc/systemd/system/dagflow.service <<EOF
[Unit]
Description=Operating system remote execution interface
Documentation=https://github.com/busyster996/dagflow.git
After=network.target nss-lookup.target

[Service]
NoNewPrivileges=true
ExecStart=/usr/local/dagflow/bin/dagflow_linux_amd64_v1 standalone
Restart=on-failure
RestartSec=10s
LimitNOFILE=infinity

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now dagflow.service
```

## Local compilation (Linux)

+ Depends on the Docker environment

```shell
git clone https://github.com/busyster996/dagflow.git
cd dagflow
make
```

## Request Example

![](images/dag_exec.png)

```text
name: 测试
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
    retryPolicy:
      interval: 1s
      maxInterval: 10s
      maxAttempts: 3
      multiplier: 3.0
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
```

### Create a task

```shell
# By default, the execution is in order.
curl -X POST -H "Content-Type:application/json" -d '"name": "test",
"timeout": "10m",
"env": [
  {
    "name": "TEST_SITE",
    "value" : "www.google.com"
  }
],
"step": [
  {
    "type": "bash", # support[python2,python3,bash,sh,cmd,powershell]
    "content": "env", # Script content
    "env": [ # Environment variable injection
      {
        "name": "TEST_SITE",
        "value" : "www.google.com"
      }
    ]
  },
  {
    "type": "bash", # support[python2,python3,bash,sh,cmd,powershell]
    "content": "curl ${TEST_SITE}", # Script content
    "env": [ # Environment variable injection
      {
        "name": "TEST_SITE",
        "value" : "www.baidu.com"
      }
    ]
  }
]' 'http://localhost:2376/api/v1/task' 

# Concurrent Execution
curl -X POST -H "Content-Type:application/json" -d '"name": "test",
"timeout": "10m",
"env": [
  {
    "name": "TEST_SITE",
    "value" : "www.google.com"
  }
],
"kind": "dag",
"step": [
  {
    "type": "bash", # support[python2,python3,bash,sh,cmd,powershell]
    "content": "env", # Script content
    "env": [ # Environment variable injection
      {
        "name": "TEST_SITE",
        "value" : "www.google.com"
      }
    ]
  },
  {
    "type": "bash", # support[python2,python3,bash,sh,cmd,powershell]
    "content": "curl ${TEST_SITE}", # Script content
    "env": [ # Environment variable injection
      {
        "name": "TEST_SITE",
        "value" : "www.baidu.com"
      }
    ]
  }
]' 'http://localhost:2376/api/v1/task'

# Customized orchestration execution
curl -X POST -H "Content-Type:application/json" -d '"name": "test",
"timeout": "10m",
"env": [
  {
    "name": "TEST_SITE",
    "value" : "www.google.com"
  }
],
"kind": dag,
"step": [
  {
    "name": "step0",
    "type": "bash", # support[python2,python3,bash,sh,cmd,powershell]
    "content": "env", # Script content
    "env": [ # Environment variable injection
      {
        "name": "TEST_SITE",
        "value" : "www.google.com"
      }
    ]
  },
  {
    "name": "step1",
    "type": "bash", # support[python2,python3,bash,sh,cmd,powershell]
    "content": "curl ${TEST_SITE}", # Script content
    "env": [ # Environment variable injection
      {
        "name": "TEST_SITE",
        "value" : "www.baidu.com"
      }
    ],
    "depends": [
      "step1"
    ]
  }
]' 'http://localhost:2376/api/v1/task'
```

### Get the task list

```shell
curl -X GET -H "Content-Type:application/json" 'http://localhost:2376/api/v1/task'
```

### Get task details

```shell
curl -X GET -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}
```

### Get task step list

```shell
curl -X GET -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}/step
```

### Get the task working directory

```shell
curl -X GET -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}/workspace
```

### Task Control

```shell
# Task to force kill
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}?action=kill

# Pause task execution [Only pending tasks can be paused]
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}?action=pause

# Pause task execution (pause for 5 minutes) [Only tasks to be run can be paused]
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}?action=pause&duration=5m

# Continue the task
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}?action=resume
```

### Get step console output

```shell
curl -X GET -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}/step/{step name}
```

### Step Control

```shell
# Steps to force kill
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}/step/{step name}?action=kill

# Pause step execution [Only pending steps can be paused]
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}/step/{step name}?action=pause

# Pause step execution (pause for 5 minutes) [Only steps to be run can be paused]
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}/step/{step name}?action=pause&duration=5m

# Continue to step
curl -X PUT -H "Content-Type:application/json" http://localhost:2376/api/v1/task/{task name}/step/{step name}?action=resume
```

[Notes]  
+ code:  
  - 0: success
  - 1001: running
  - 1002: failed
  - 1003: not found
  - 1004: pending
  - 1005: paused
  - 1006: skipped

## Script language support
+ [bash/sh/ps1/bat/python2/python3/mkdir/touch](worker/runner/README.md)
+ [lua](worker/runner/lua/README.md)
+ [scp](worker/runner/scp/README.md)
+ [ssh](worker/runner/ssh/README.md)
+ [sftp](worker/runner/sftp/README.md)
+ [yaegi](worker/runner/yaegi/README.md)
+ [docker](worker/runner/docker/README.md)
+ [kubectl](worker/runner/kubectl/README.md)

## Swagger API documentation
[Swagger API documentation](docs/swagger.json)
