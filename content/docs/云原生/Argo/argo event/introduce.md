---
title: "introduce"
weight: 2
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---

# 📡 Argo Events 架构理解

Argo Events 是 Argo 项目的一个事件驱动引擎，专为 Kubernetes 原生环境设计，支持通过外部事件自动触发工作流、部署、任务执行等动作，是构建 GitOps 自动化的关键工具之一。

---

## ✅ 核心概念

### 1. **EventSource**
定义事件来源（例如 Webhook、S3 上传、Kafka 消息、Cron 计划任务等），Argo Events 会监听这些事件源并产生事件。

```yaml
apiVersion: argoproj.io/v1alpha1
kind: EventSource
metadata:
  name: webhook-source
spec:
  webhook:
    example:
      endpoint: /trigger
      method: POST
      port: 12000
```

---

### 2. **Sensor**
Sensor 监听 EventSource 的事件，当满足触发条件时，激活 Trigger。可以设置依赖关系、条件判断。

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: webhook-sensor
spec:
  dependencies:
    - name: example-dep
      eventSourceName: webhook-source
      eventName: example
  triggers:
    - template:
        name: trigger-workflow
        k8s:
          group: argoproj.io
          version: v1alpha1
          resource: workflows
          operation: create
          source:
            resource:
              apiVersion: argoproj.io/v1alpha1
              kind: Workflow
              metadata:
                generateName: triggered-wf-
              spec:
                entrypoint: whalesay
                templates:
                  - name: whalesay
                    container:
                      image: docker/whalesay
                      command: [cowsay]
                      args: ["Hello from event"]
```

---

### 3. **Trigger**
指定事件触发后的行为，例如启动 Argo Workflow、Kubernetes Job、发送通知等。支持模板化、参数注入、条件触发等复杂逻辑。

---

## 🔄 支持的触发器类型

| 类型 | 描述 |
|------|------|
| `webhook` | 接收 HTTP 请求（如 GitLab Push、Jenkins 通知） |
| `s3` / `minio` | 文件上传触发工作流（例如日志、模型） |
| `kafka` / `nats` | 通过消息队列异步触发流程 |
| `calendar` / `cron` | 定时执行任务，代替 CronWorkflow |
| `file` | 本地文件系统变化触发 |
| `resource` | 监听 Kubernetes 对象变化（如 Pod 状态） |
| `amqp`, `mqtt`, `sqslike` | 企业集成常用消息源 |
| `slack` | 消息触发自动响应（如 bot 流程） |

完整支持列表参考官方文档：https://argoproj.github.io/argo-events/sources/

---

## 🧩 架构图

![Argo Events Architecture](https://argoproj.github.io/argo-events/assets/events-architecture.png)

### 架构简要说明

| 组件 | 功能 |
|------|------|
| **EventSource Controller** | 创建并监听事件来源 |
| **Sensor Controller** | 检测事件是否满足触发条件 |
| **EventBus（NATS）** | 内部事件总线（可选）用于异步解耦 |
| **Trigger Handler** | 执行具体的触发操作（如创建 Workflow） |

---

## 📚 推荐资料

- 官方文档：https://argoproj.github.io/argo-events/
- GitHub：https://github.com/argoproj/argo-events
- 示例合集：https://github.com/argoproj/argo-events/tree/master/examples

