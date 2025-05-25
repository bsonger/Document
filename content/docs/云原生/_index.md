---
title: "Cloud Native"
weight: 2
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: true
---

# 生产环境云原生部署架构图

## 架构图（文字版）

```
+--------------------------------------------------------------------------------------------------------+
|                                             Kubernetes 集群                                             |
|                                                                                                        |
|  +---------------------+    +-------------------+    +---------------------+                          |
|  |     Cilium CNI      |<-->| NetworkPolicy/EBPF |<-->|  Pod 网络通信与安全策略   |                          |
|  +---------------------+    +-------------------+    +---------------------+                          |
|                                                                                                        |
|  +---------------------+                                                                           |
|  |  cert-manager       |----> 自动签发/续期 TLS 证书（整合 ACME / Vault / Issuer）                    |
|  +---------------------+                                                                           |
|                                                                                                        |
|  +---------------------+                                                                           |
|  |     Keycloak        |----> OIDC/SAML 登录认证（整合 ArgoCD、Harbor、Grafana 等）                  |
|  +---------------------+                                                                           |
|                                                                                                        |
|  +---------------------+    +---------------------+                                                  |
|  |       Harbor        |<--> 镜像仓库（私有仓库，结合 Keycloak OIDC）      |                            |
|  +---------------------+    +---------------------+                                                  |
|                                                                                                        |
|  +----------------------+   +----------------------+   +---------------------------+                |
|  |  Argo Workflows      |   |   Argo CD            |   | Argo Events / Rollouts    |                |
|  |  (CI流水线调度器)     |   | (GitOps 部署)        |   | （事件驱动 / 灰度发布）      |                |
|  +----------------------+   +----------------------+   +---------------------------+                |
|     ↑ GitOps触发构建           ↑ 从Git拉取应用部署配置         ↑ CRDs / 事件驱动CI流程                             |
|                                                                                                        |
|  +---------------------+    +---------------------+    +--------------------------+                  |
|  |  Prometheus         |<---| Node Exporter / Kubelet |<-- Kubernetes 指标收集                      |
|  +---------------------+    +---------------------+    +--------------------------+                  |
|         ↓ Alertmanager        ↓ Grafana (集成Keycloak登录)                                          |
|         告警规则与通知        可视化监控面板                                                         |
|                                                                                                        |
|  +---------------------+                                                                           |
|  |         Loki        |<--- 日志收集（整合 promtail / fluentbit）                                   |
|  +---------------------+                                                                           |
|                                                                                                        |
+--------------------------------------------------------------------------------------------------------+
```

# ☁️ Cloud-Native Production Environment Architecture

## 🧩 Central Platform:
- 🧱 **Kubernetes Cluster**: Orchestrates containers and manages workloads.

## 🛠️ Key Components:
- 🧑‍💻 **Backstage**: Developer portal providing software catalog, CI/CD integration, documentation, and monitoring dashboards.
- 📦 **Harbor**: Private container registry to manage container images securely.
- 🕸️ **Cilium**: Provides networking, security policies, observability using eBPF technology.
- 🔐 **Cert-manager**: Automates the issuance and renewal of TLS certificates.
- 🛡️ **Keycloak**: Identity and access management system providing authentication and authorization.
- 🚀 **Argo Suite**:
    - 📥 **Argo CD**: GitOps continuous deployment tool.
    - 🔁 **Argo Workflows**: Workflow orchestration and CI pipelines.
    - 🧲 **Argo Events/Rollouts**: Event-driven workflows and canary deployments.
- 📈 **Prometheus**: Monitoring and alerting toolkit.
- 📄 **Loki**: Log aggregation system for collecting and managing logs.

## 🔗 Integration and Interactions:
- **Backstage** interacts with all components, providing a unified view for developers and administrators.
- **Harbor** is integrated with Kubernetes and Keycloak for secure image management.
- **Cilium** manages cluster networking, providing observability data to Prometheus.
- **Cert-manager** integrates with Kubernetes and Keycloak for automated certificate management.
- **Keycloak** provides centralized authentication services for Kubernetes, Backstage, Harbor, and Argo Suite.
- **Argo Suite** integrates with Kubernetes and Git repositories for streamlined deployments and CI/CD processes.
- **Prometheus** gathers metrics from all components, with dashboards presented via Backstage.
- **Loki** aggregates logs from Kubernetes and applications, accessible via dashboards.

✨ This architecture ensures scalability, security, observability, and ease of use for cloud-native application deployment and management.
