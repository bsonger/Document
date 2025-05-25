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

# Backstage 详解介绍

## 🎯 什么是 Backstage？

Backstage 是由 Spotify 开源的开发者门户平台框架，用于构建集中化的开发者体验平台（Developer Portal），帮助组织统一管理软件组件、CI/CD 流程、监控工具、文档等资源。

GitHub 地址：[https://github.com/backstage/backstage](https://github.com/backstage/backstage)

---

## 🧠 Backstage 原理与架构

Backstage 基于插件架构，提供统一的 UI 界面来聚合不同的开发工具和资源。其核心原理为：

- **软件目录（Software Catalog）**：注册管理组织内的所有软件组件（如微服务、库、文档等）。
- **插件机制**：通过插件系统实现功能扩展，如 CI/CD、监控、API 管理等。
- **统一的 UI 界面**：基于 React 构建，通过插件动态渲染内容。
- **配置驱动**：基于 YAML 文件配置软件组件元数据。

```
+---------------------+
|     Developer       |
+---------------------+
         |
         v
+---------------------+         +---------------------+
|  Backstage Frontend | <-----> | Backstage Backend   |
+---------------------+         +---------------------+
         |                               |
         v                               v
+----------------+           +--------------------------+
|  Plugins       |           | Software Catalog + Auth  |
+----------------+           +--------------------------+
         |
         v
+---------------------------+
| External Tools (CI/CD,   |
| Monitoring, Docs, etc.)  |
+---------------------------+
```

---

## 🧩 核心功能模块

- **软件目录（Catalog）**：集中注册和管理所有组件。
- **TechDocs**：支持以 Markdown 为基础的文档自动生成和展示。
- **CI/CD 集成**：支持 Jenkins、GitHub Actions、GitLab 等。
- **监控可视化**：集成 Grafana、Prometheus、Sentry 等工具。
- **权限控制**：支持 RBAC、LDAP、OAuth 等。

---

## 🧱 技术栈

| 组件         | 技术                          |
|--------------|-------------------------------|
| 前端         | React + TypeScript            |
| 后端         | Node.js + Express             |
| 插件系统     | 插件基于 React 和 Node 构建   |
| 配置管理     | YAML                           |
| 部署         | 支持 Docker / Kubernetes      |

---

## 🔌 插件机制

Backstage 插件分为两种：

- **前端插件**：提供 UI 功能，独立模块化。
- **后端插件**：与外部服务交互，如数据库、CI 服务等。

可以通过 `npx @backstage/create-app` 快速初始化插件和项目。

---

## 🚀 部署方式

支持以下部署方式：

- Docker Compose
- Kubernetes（Helm Charts 可用）
- 静态文件 + Node 后端服务

---

## 🌐 使用场景

- 多团队协作统一入口
- DevOps 工具集成
- 内部服务发现与治理
- 自动化文档平台
- 微服务治理门户

---

## 🌍 社区与生态

- 提供丰富的官方插件
- 支持第三方插件市场
- 与 CNCF 等社区合作紧密

---

## 📚 参考文档

- 官方文档：https://backstage.io/docs
- 插件市场：https://backstage.io/plugins
