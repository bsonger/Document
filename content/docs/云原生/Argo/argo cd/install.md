---
title: "install"
weight: 2
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---

# 🚀 Argo CD + ApplicationSet 安装指南（生产推荐）

本指南包含 Argo CD 与 ApplicationSet Controller 的完整安装流程，适用于生产环境部署。

---

## ✅ 一、准备工作

- Kubernetes 集群（v1.21+ 推荐）
- Helm 3.x
- 具备集群管理员权限
- 域名（如使用 SSO 或 Ingress）

---

## 🧰 二、添加 Argo Helm 仓库

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
```

---

## 📦 三、安装 Argo CD（含 ApplicationSet）

```bash
helm install argocd argo/argo-cd   --namespace argocd --create-namespace   --set server.service.type=ClusterIP   --set dex.enabled=true   --set applicationset.enabled=true   --set configs.params.server.insecure=false
```

### 验证服务状态

```bash
kubectl get pods -n argocd
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

访问 UI: `https://localhost:8080`

### 获取初始 admin 密码

```bash
kubectl -n argo get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

---

## 🧩 四、ApplicationSet Controller 简介

ApplicationSet 允许你**批量创建和动态生成多个 Application**，适用于：

- 多微服务部署
- 多环境/多集群 GitOps
- 自动化生成 YAML 应用资源

常用 Generators 包括：
- Git Generator（根据目录结构生成）
- List Generator（静态定义多个服务）
- Cluster Generator（为每个集群生成一个应用）

官方文档：https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/

---

## 🔒 五、启用 SSO（可选，生产建议）

修改 `argocd-cm` 和 `argocd-secret` 配置，启用 OIDC/OAuth2 认证（如 GitHub、Dex、Auth0）。

> 完整配置见生产部署文档

```yaml
data:
  oidc.config: |
    name: GitHub
    issuer: https://github.com/login/oauth
    clientID: YOUR_CLIENT_ID
    clientSecret: $oidc.github.clientSecret
```

---

## 📣 六、集成 Argo CD Notifications（可选）

```bash
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj-labs/argocd-notifications/stable/manifests/install.yaml
```

添加通知通道（Slack/Webhook/Email）并配置 `argocd-notifications-cm` 和 Secret。

---

## 🛡 七、生产部署建议

| 类别 | 推荐做法 |
|------|----------|
| 高可用性 | 为 controller/server 设置副本数，使用 LoadBalancer 或 Ingress |
| 安全性 | 启用 SSO，禁用默认 admin，配置 RBAC |
| 多集群支持 | 注册多个集群，使用 ApplicationSet 进行批量管理 |
| 可观测性 | 启用 Prometheus、Grafana 监控 |
| 资源隔离 | 使用 Argo CD Projects 区分业务/环境 |
| 存储与备份 | 所有 manifests 存入 Git，配置定期备份 |

---

## 🔗 推荐文档

- Argo CD Helm Chart：https://github.com/argoproj/argo-helm/tree/main/charts/argo-cd
- ApplicationSet 示例：https://github.com/argoproj-labs/applicationset
- Argo CD Operator Manual：https://argo-cd.readthedocs.io/en/stable/operator-manual/

