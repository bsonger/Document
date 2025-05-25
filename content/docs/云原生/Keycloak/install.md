# 🏗️ Keycloak 生产环境部署指南

## 📦 部署方式概览

Keycloak 支持多种部署方式，常见包括：

- Docker / Podman 容器部署
- Kubernetes / OpenShift 部署（推荐生产使用）
- 使用独立 JAR 文件运行
- Helm Chart 部署（适合 Kubernetes 环境）

本指南以 Kubernetes 环境 + PostgreSQL 数据库为例进行说明。

## 🔧 环境准备

### 系统要求
- Kubernetes 集群（v1.21+）
- Ingress Controller（如 nginx-ingress）
- PostgreSQL 数据库
- 持久化存储（PVC）
- Cert Manager（可选，用于 TLS）

### 建议资源配置
| 组件      | CPU   | 内存   |
|-----------|-------|--------|
| Keycloak  | 2核   | 4Gi    |
| PostgreSQL| 1核   | 2Gi    |

## 🐳 使用 Helm 部署 Keycloak

### 1️⃣ 添加 Helm 仓库
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### 2️⃣ 安装 PostgreSQL（可选，或使用外部 DB）
```bash
helm install keycloak-postgresql bitnami/postgresql \
  --set auth.postgresPassword=yourpassword \
  --set primary.persistence.enabled=true
```

### 3️⃣ 安装 Keycloak
```bash
helm install keycloak bitnami/keycloak \
  --set auth.adminUser=admin \
  --set auth.adminPassword=adminpassword \
  --set externalDatabase.host=keycloak-postgresql \
  --set externalDatabase.user=postgres \
  --set externalDatabase.password=yourpassword \
  --set proxy=edge \
  --set ingress.enabled=true \
  --set ingress.hostname=keycloak.example.com
```

## 🔒 TLS 配置（可选）
如果集群中部署了 cert-manager，可通过如下方式配置 TLS：
```yaml
ingress:
  enabled: true
  hostname: keycloak.example.com
  tls: true
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  extraTls:
    - hosts:
        - keycloak.example.com
      secretName: keycloak-tls
```

## 📁 数据持久化配置
Keycloak 默认开启持久化，确保 statefulset 挂载了 PVC：
```yaml
persistence:
  enabled: true
  size: 8Gi
  storageClass: standard
```

## 🩺 健康检查与监控

- 默认开启 liveness/readiness probe
- 可通过 Prometheus Operator 集成 Keycloak Exporter 实现监控

## 🛠️ 常见问题排查

| 问题 | 解决方法 |
|------|----------|
| Admin 页面打不开 | 检查 ingress 或 admin 密码是否正确 |
| 数据丢失 | 确保 PVC 正常挂载并配置持久化 |
| 登录失败 | 检查数据库连接配置是否正确 |

## 📚 参考链接

- 官方文档：[https://www.keycloak.org/documentation](https://www.keycloak.org/documentation)
- Helm Chart：[https://artifacthub.io/packages/helm/bitnami/keycloak](https://artifacthub.io/packages/helm/bitnami/keycloak)
- GitHub：[https://github.com/keycloak/keycloak](https://github.com/keycloak/keycloak)

---

✅ 建议将 Keycloak 与 LDAP/AD 集成，并定期备份数据库，启用高可用部署以增强稳定性。
