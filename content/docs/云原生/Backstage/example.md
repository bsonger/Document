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

# Backstage 生产环境使用示例（复杂用例）

本文展示一个包含多个插件、权限控制、自定义服务目录、CI/CD 集成、文档生成和 Kubernetes 可观测性的复杂 Backstage 使用场景。

---

## 🧩 系统架构图示意

```
开发者 --> Backstage Portal --> 插件系统
                          |
                          +--> Catalog (服务注册)
                          +--> TechDocs (文档平台)
                          +--> ArgoCD Plugin (CI/CD 状态)
                          +--> Kubernetes Plugin (集群状态)
                          +--> Auth (GitHub OAuth)
                          +--> External Monitoring (Prometheus/Grafana)
```

---

## 🗃 软件目录配置（catalog-info.yaml）

```yaml
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: payment-service
  description: 支付服务，负责处理所有支付相关逻辑
  annotations:
    github.com/project-slug: myorg/payment-service
    backstage.io/techdocs-ref: dir:.
    argocd/app-name: payment-service
    prometheus.io/rule: payment-alerts
spec:
  type: service
  lifecycle: production
  owner: team-payments
  system: payments-platform
  dependsOn:
    - component:default/database
```

---

## 🔒 身份认证与权限控制

使用 GitHub OAuth + RBAC 策略：

```yaml
auth:
  providers:
    github:
      production:
        clientId: <client-id>
        clientSecret: <client-secret>
        enterpriseInstanceUrl: https://github.myorg.com

permission:
  enabled: true
  policies:
    - policy: allow
      resourceType: component
      actions: ['read']
      subjects:
        - group: developers
    - policy: deny
      resourceType: component
      actions: ['delete']
```

---

## 🚀 CI/CD 可视化（ArgoCD 插件）

在 `app-config.yaml` 中配置：

```yaml
argocd:
  appLocatorMethods:
    - type: config
      instances:
        - name: main
          url: https://argocd.example.com
          token: ${ARGOCD_TOKEN}
```

前端插件可展示 Argo 应用部署状态。

---

## 📄 技术文档（TechDocs）

项目结构：

```
.
├── docs
│   └── index.md
├── mkdocs.yml
└── catalog-info.yaml
```

示例 `mkdocs.yml`:

```yaml
site_name: Payment Service Docs
nav:
  - Home: index.md
plugins:
  - techdocs-core
```

构建并部署 TechDocs：

```bash
yarn techdocs:build
yarn techdocs:publish
```

---

## 🛰 Kubernetes 插件配置

配置多个集群的访问凭证：

```yaml
kubernetes:
  clusterLocatorMethods:
    - type: config
      clusters:
        - name: prod-cluster
          url: https://k8s-prod-api.example.com
          authProvider: serviceAccount
          serviceAccountToken: ${K8S_TOKEN}
          skipTLSVerify: false
```

可视化展示 Pod、Service、Deployment 状态。

---

## 📊 监控系统集成（Prometheus + Grafana）

Prometheus 配置：

```yaml
prometheus:
  instances:
    - name: main
      url: http://prometheus.monitoring.svc:9090
```

Grafana：

可通过 iframe 嵌入 Grafana 仪表盘至自定义插件或页面。

---

## 🎯 结语

通过 Backstage 的插件化架构，可以将 DevOps 工具链整合进一个统一平台，提升开发者效率。此示例展示了一个集成 ArgoCD、Kubernetes、Prometheus、TechDocs 等复杂插件的生产环境用例。

---

## 📚 参考链接

- https://backstage.io/docs
- https://github.com/backstage/backstage
