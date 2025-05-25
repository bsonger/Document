---
title: "example"
weight: 4
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---
# 🚀 最复杂的 Argo Workflow 示例

本示例涵盖 Argo Workflow 的核心高级功能，包括：

- DAG + 步骤式混用
- 参数传递与输出
- Artifact 使用
- Script 模板
- Sidecar 容器
- 条件执行（when）
- 重试策略
- 并发控制
- TTL 策略
- 持久化存储（Artifact S3 示例）
- 手动暂停与继续
- 工作流模板复用（WorkflowTemplate）

---

## 🧾 Workflow YAML 示例

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: complex-workflow-
spec:
  entrypoint: main-dag
  arguments:
    parameters:
      - name: env
        value: "prod"

  ttlStrategy:
    secondsAfterSuccess: 3600

  parallelism: 2

  templates:
    - name: main-dag
      dag:
        tasks:
          - name: generate
            template: generate-data
          - name: process
            dependencies: [generate]
            template: process-data
            arguments:
              artifacts:
                - name: input-data
                  from: "{{tasks.generate.outputs.artifacts.out-data}}"
          - name: conditional-step
            dependencies: [generate]
            template: conditional-script
            when: "{{inputs.parameters.env}} == prod"
          - name: deploy
            dependencies: [process, conditional-step]
            template: deploy-app

    - name: generate-data
      outputs:
        artifacts:
          - name: out-data
            path: /tmp/output.txt
      container:
        image: alpine
        command: [sh, -c]
        args: ["echo 'data-content' > /tmp/output.txt"]

    - name: process-data
      inputs:
        artifacts:
          - name: input-data
            path: /tmp/input.txt
      container:
        image: python:3.8
        command: [python]
        args:
          - -c
          - |
            with open('/tmp/input.txt') as f:
                data = f.read()
            print(f"Processed: {data}")

    - name: conditional-script
      script:
        image: bash:5.1
        command: [bash]
        source: |
          echo "Running only in production"

    - name: deploy-app
      container:
        image: curlimages/curl
        command: ["curl"]
        args: ["-X", "POST", "http://example.com/deploy"]

    - name: sidecar-demo
      container:
        image: busybox
        command: [sh, -c]
        args: ["echo main container && sleep 10"]
      sidecars:
        - name: logger
          image: busybox
          command: ["sh", "-c"]
          args: ["while true; do echo sidecar running; sleep 2; done"]
          mirrorVolumeMounts: true

    - name: retry-step
      retryStrategy:
        limit: 3
        retryPolicy: "Always"
      container:
        image: alpine
        command: [sh, -c]
        args: ["exit 1"]

    - name: approval-step
      suspend: {}
```

---

## 🧠 功能说明对照表

| 功能模块 | 说明 |
|----------|------|
| DAG 模式 | 使用 `main-dag` 模板构建任务依赖 |
| Artifact | 任务间通过 `/tmp/output.txt` 文件传递数据 |
| 参数传递 | 通过 `arguments.parameters` 控制部署环境 |
| 条件执行 | 只有在 `env == prod` 时才运行 `conditional-script` |
| Script 模板 | 内联 bash 脚本 |
| Sidecar 容器 | `sidecar-demo` 展示日志采集 Sidecar |
| 重试策略 | `retry-step` 会失败并最多重试 3 次 |
| Suspend 手动批准 | `approval-step` 模拟人工审核点 |
| TTL 策略 | 完成 3600 秒后自动清理 |
| 并发限制 | `parallelism: 2` 限制并发任务数为 2 |

---

## 📦 要求的环境依赖

- 已安装 Argo Workflows Controller & Server
- Artifact 存储（如 MinIO、S3）已配置（示例略去具体 config）
- 如果使用 WorkflowTemplate，可拆分 `generate`, `process`, `deploy` 为可复用模板

---

## 📌 小提示

- 可以通过 Argo UI 或 CLI 提交此 Workflow
- 若需手动 Resume suspend 任务：
  ```bash
  argo resume <workflow-name>
  ```

---

## 🔗 相关参考

- Argo 官方文档：https://argoproj.github.io/argo-workflows/
- 示例仓库：https://github.com/argoproj/argo-workflows/tree/master/examples
