---
title: "install"
weight: 3
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
bookCollapseSection: false
---
# 📦 Argo Workflows 安装与启动参数详解

## ✅ 安装 Argo Workflows

### 方法一：官方快速安装

```bash
kubectl create namespace argo
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo-workflows/stable/manifests/install.yaml
kubectl -n argo port-forward deployment/argo-server 2746:2746
```

### 方法二：使用 Helm 安装（推荐）

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

helm install argo-workflows argo/argo-workflows   --namespace argo --create-namespace   --set server.secure=true   --set controller.workflowNamespaces="{argo}"   --set useDefaultArtifactRepo=true
```

---

## ⚙️ argo-server 常见启动参数说明

| 参数 | 作用 |
|------|------|
| `--secure` | 是否启用 TLS（建议开启） |
| `--auth-mode` | 认证模式，可选：`server`, `client`, `sso`, `hybrid` |
| `--configmap` | 配置 Workflow Controller 的配置来源 |
| `--namespaced` | 是否启用命名空间隔离（每个用户只能看自己 namespace 的 workflow） |
| `--basehref` | 设置 UI 路由前缀（如部署在子路径下） |
| `--namespace` | 设置默认的 namespace |
| `--loglevel` | 设置日志级别（例如 `info`, `debug`） |

---

## 🔐 启用 SSO（Single Sign-On）

Argo Workflows 支持基于 OIDC 的 SSO 认证。以下为基本配置步骤：

### 步骤一：配置 `argo-server` 使用 SSO 模式

```bash
--auth-mode sso
--configmap workflow-controller-configmap
```

### 步骤二：创建 SSO 配置（ConfigMap）

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: workflow-controller-configmap
  namespace: argo
data:
  sso: |
    issuer: https://accounts.google.com
    clientId: my-client-id
    clientSecret: $SSO_CLIENT_SECRET
    redirectUrl: https://argo.example.com/oauth2/callback
    scopes:
      - email
      - profile
    rbac:
      enabled: true
      defaultPolicy: 'role:readonly'
      scopes: ['groups']
```

> 注意：clientId 与 clientSecret 来自你的 OIDC 提供商（如 Google、Auth0、Dex）

---

## 💾 数据持久化配置（连接数据库）

Argo Workflows 默认使用 Kubernetes 的状态进行调度，但也可集成 PostgreSQL 以提升高可用和性能（通常用于 Argo Server + Workflow Archive）。

### 步骤一：配置 PostgreSQL 信息

```yaml
controller:
  persistence:
    archive: true
    postgresql:
      host: <postgres-host>
      port: 5432
      database: argo
      tableName: argo_workflows
      userName: argo
      passwordSecret:
        name: argo-postgres-secret
        key: password
```

### 步骤二：创建密钥

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: argo-postgres-secret
  namespace: argo
type: Opaque
data:
  password: <base64-encoded-password>
```

---

## 🔍 验证部署

```bash
kubectl get pods -n argo
kubectl logs deployment/argo-server -n argo
kubectl get workflow -n argo
```

访问 UI: http://localhost:2746

---

## 📚 更多参考

- 官方 Helm Chart 文档: https://github.com/argoproj/argo-helm/tree/main/charts/argo-workflows
- SSO 文档: https://argoproj.github.io/argo-workflows/argo-server-sso/
- Artifact Persistence: https://argoproj.github.io/argo-workflows/workflow-archive/

