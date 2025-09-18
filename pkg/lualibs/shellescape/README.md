# shellescape

## Usage

```lua
local shellescape = require("shellescape")

escaped = shellescape.quote("foo bar baz")
-- 'foo bar baz'

escaped_command = shellescape.quote_command({ "echo", "foo bar baz" })
-- echo 'foo bar baz'
```