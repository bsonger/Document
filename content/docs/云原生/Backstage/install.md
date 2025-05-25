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
# Backstage Helm 安装指南（生产环境）

本文档介绍如何使用 Helm 在 Kubernetes 中部署生产级别的 Backstage 开发者门户。

---

## 🧰 前置条件

确保以下组件已安装并配置：

- Kubernetes 集群（版本 ≥ 1.21）
- Helm 3.x
- Docker（构建镜像用）
- 可访问的容器镜像仓库
- PostgreSQL 数据库（可选，Helm 可自动安装）

---

## 🚀 步骤一：构建并推送 Backstage 镜像

```bash
# 初始化项目（如未创建）
npx @backstage/create-app

cd my-backstage-app

# 构建生产镜像
docker build -t my-org/backstage:latest .
docker push my-org/backstage:latest
```

---

## 🔧 步骤二：添加 Helm 仓库

```bash
helm repo add backstage https://backstage.github.io/charts
helm repo update
```

---

## 📝 步骤三：创建 `values.yaml` 配置文件

以下是一个简化的生产配置：

```yaml
image:
  repository: my-org/backstage
  tag: latest
  pullPolicy: IfNotPresent

replicaCount: 2

ingress:
  enabled: true
  host: backstage.example.com
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
  tls:
    - hosts:
        - backstage.example.com
      secretName: backstage-tls

postgresql:
  enabled: true
  postgresqlPassword: strong-password

appConfig:
  app:
    baseUrl: https://backstage.example.com
  backend:
    baseUrl: https://backstage.example.com
    cors:
      origin: https://backstage.example.com
    database:
      client: pg
      connection:
        host: postgres
        user: postgres
        password: strong-password
        database: backstage
```

> ✅ 推荐使用自定义域名 + HTTPS + TLS。

---

## 🛠 步骤四：安装 Backstage

```bash
helm install backstage backstage/backstage -n backstage --create-namespace -f values.yaml
```

---

## 🔁 步骤五：更新部署（如有变更）

```bash
helm upgrade backstage backstage/backstage -n backstage -f values.yaml
```

---

## 🧪 步骤六：验证部署

```bash
kubectl get pods -n backstage
kubectl get svc -n backstage
```

访问：https://backstage.example.com

---

## 📌 附加建议

- 配合 cert-manager 自动生成 TLS 证书
- 使用 External PostgreSQL 实现持久化和高可用
- 使用 Fluent Bit 收集日志，配合 Loki 或 ELK

---

## 📚 参考链接

- https://backstage.io/docs
- https://github.com/backstage/charts
