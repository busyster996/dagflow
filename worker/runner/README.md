# Default runner executor

## bash/sh/cmd/powershell

```text
step:
  - name: sh
    type: sh
    content: |-
      ping -c 4 1.1.1.1
```

## python2/python3

```text
step:
  - name: python
    type: python
    content: |-
      import os
      os.system('ping -c 4 1.1.1.1')
```

## mkdir

```text
step:
  - name: mkdir
    type: mkdir
    content: |-
      path: a/b/c/d/e
```

## touch

```text
step:
  - name: touch
    type: touch
    content: |-
      path: a/b/c/d/e/f.txt
      overwrite: true
      content: |-
        content of f.txt      
```