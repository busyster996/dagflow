# sftp

```text
step:
  - name: sftp
    type: sftp
    content: |-
      method: password
      host: 192.168.1.1
      port: 22
      username: root
      password: <PASSWORD>
      direction: upload
      source: xxx.txt
      target: /tmp/xxx.txt
```