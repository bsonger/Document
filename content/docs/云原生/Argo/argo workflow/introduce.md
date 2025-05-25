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
# 📘 Argo Workflows 基础概念

## ✅ Argo Workflows 是什么？

Argo Workflows 是一个用于在 Kubernetes 上运行 **容器化工作流** 的开源工具。它允许你通过自定义的 YAML 文件，定义一系列任务的执行流程，这些任务以容器的形式运行，并可以通过步骤或依赖顺序进行编排。

它主要用于：
- **CI/CD 流程编排**
- **数据处理流水线**
- **自动化测试任务**
- **批处理作业运行**
- **机器学习任务执行（如 Kubeflow Pipelines）**

Argo Workflows 通过 CRD（自定义资源定义）扩展了 Kubernetes，原生集成在集群中运行，无需引入外部服务。

---

## 🔧 基本术语解释

| 术语 | 含义 |
|------|------|
| **Workflow** | 工作流的实例，是任务的运行单元 |
| **Template** | 可复用的任务模板，每个 template 定义一个操作 |
| **Step** | 指定任务的顺序执行逻辑 |
| **DAG** | 有向无环图，用于定义任务之间的依赖 |
| **Parameter** | 参数变量，允许模板传入数据 |
| **Artifact** | 文件或数据对象，在任务之间传递或保存结果 |
| **Entrypoint** | 定义 Workflow 的起始模板 |
| **WorkflowTemplate** | 可以全局复用的模板定义（跨 Workflow 共享） |

示意结构：
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: hello-world-
spec:
  entrypoint: hello
  templates:
  - name: hello
    container:
      image: alpine:3.7
      command: ["echo"]
      args: ["hello world"]
```

---

## 🧩 Argo Workflows 架构组件

![Argo Workflows Architecture](https://argoproj.github.io/argo-workflows/assets/argo-workflows-architecture.png)

| 组件 | 作用 |
|------|------|
| **Workflow Controller** | 核心控制器，监听并调度 Workflow 的每一步任务 |
| **argo-server (UI Server)** | 提供 Web UI 管理界面，可查看执行状态、日志 |
| **CRD** | 包括 Workflow、WorkflowTemplate、CronWorkflow 等资源，扩展 Kubernetes |
| **Executor** | Sidecar 容器，负责执行每个任务的命令并报告状态 |
| **Artifact Repository** | 任务之间的数据传递（如 S3、OSS、PVC） |

---

## 📚 推荐学习资源

- 官网文档：[https://argoproj.github.io/argo-workflows](https://argoproj.github.io/argo-workflows)
- GitHub 项目：[https://github.com/argoproj/argo-workflows](https://github.com/argoproj/argo-workflows)
- 样例库：[https://github.com/argoproj/argo-workflows/tree/master/examples](https://github.com/argoproj/argo-workflows/tree/master/examples)

