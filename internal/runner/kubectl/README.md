# kubectl

## kubectl@status
```text
step:
  - name: kubectl@status
    type: kubectl@status
    env:
      - name: KUBECONFIG
        value: /root/.kube/admin.conf
    content: |-
        resources:
          - name: nginx
            namespace: default
            kind: Deployment
```

## kubectl@update
```text
step:
  - name: kubectl@update
    type: kubectl@update
    env:
      - name: KUBECONFIG
        value: /root/.kube/admin.conf
    content: |-
      resources:
        - name: nginx
          namespace: default
          imageTag: latest
          kind: Deployment
```