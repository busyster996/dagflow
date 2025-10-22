# yaegi 

exec yaegi type script

## Usage

```text
step:
  - name: yaegi
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
```