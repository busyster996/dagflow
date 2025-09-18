# cmd

## Usage

```lua
local cmd = require("cmd")
local res, err = cmd.exec("bash", "-c", "echo hello world")
if err then error(err) end
print(res)

--- or timeout seconds
local res, err = cmd.exec_by_timeout(15, "bash", "-c", "sleep 5; echo hello world")
if err then error(err) end
print(res)
```