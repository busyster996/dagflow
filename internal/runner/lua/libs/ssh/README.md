# ssh

## Usage

```lua
local ssh = require("ssh")
session, err = ssh.auth{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}
command = session:execute{command = "echo true"}
print(command.output)
```