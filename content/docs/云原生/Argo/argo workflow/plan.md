---
title: "plan"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---
# 📘 Argo Workflows 学习计划 - 第 1 周

## 🎯 学习目标
- 理解 Argo Workflows 的基本概念、架构和组件。
- 能够在本地或 Kubernetes 集群中安装并运行 Workflows。
- 掌握编写基础 Workflow YAML 的能力，了解参数化、模板复用及定时调度任务（CronWorkflow）。

---

## 📅 每日任务安排

### ✅ 周一：Workflows 概述与架构

- 学习内容：
  - Argo Workflows 是什么？解决什么问题？
  - 工作流的基本术语（模板、步骤、参数等）
  - 架构组件：Workflow Controller、UI Server、CRD

- 资料：
  - 官方文档：https://argoproj.github.io/argo-workflows/
  - 架构图参考：![Argo Workflows Architecture](https://argoproj.github.io/argo-workflows/assets/argo-workflows-architecture.png)

---

### ✅ 周二：安装 Argo Workflows

- 操作步骤：
  ```bash
  kubectl create namespace argo
  kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo-workflows/stable/manifests/install.yaml
  kubectl -n argo port-forward deployment/argo-server 2746:2746
  ```

- 访问 UI：http://localhost:2746

- 验证安装：
  ```bash
  kubectl get pods -n argo
  ```

---

### ✅ 周三：创建第一个 Workflow

- 示例 YAML：
  ```yaml
  apiVersion: argoproj.io/v1alpha1
  kind: Workflow
  metadata:
    generateName: hello-world-
  spec:
    entrypoint: whalesay
    templates:
    - name: whalesay
      container:
        image: docker/whalesay
        command: [cowsay]
        args: ["hello world"]
  ```

- 提交命令：
  ```bash
  kubectl create -f hello-world.yaml -n argo
  ```

- 查看执行：Argo UI 或 `kubectl get wf -n argo`

---

### ✅ 周四：参数与模板复用

- 使用参数：
  ```yaml
  args:
    parameters:
    - name: message
      value: "Hello Argo"
  ```

- 模板复用：定义多个模板并通过 `steps` 引用

- 示例项目：https://github.com/argoproj/argo-workflows/tree/master/examples

---

### ✅ 周五：定时任务（CronWorkflow）

- 示例 YAML：
  ```yaml
  apiVersion: argoproj.io/v1alpha1
  kind: CronWorkflow
  metadata:
    name: hello-cron
  spec:
    schedule: "*/5 * * * *"
    workflowSpec:
      entrypoint: whalesay
      templates:
      - name: whalesay
        container:
          image: docker/whalesay
          command: [cowsay]
          args: ["Hello from cron!"]
  ```

- 提交命令：
  ```bash
  kubectl create -f cron-workflow.yaml -n argo
  ```

- 验证命令：`kubectl get cronworkflow -n argo`

---

## 📝 总结建议
- 每晚完成操作记录一份运行截图或 UI 截图。
- 笔记记录遇到的问题和解决方法。
- 可尝试将 Workflow 与其他容器任务组合成多步骤流程。
