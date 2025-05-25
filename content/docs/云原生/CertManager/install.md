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

# cert-manager 生产环境部署指南

![cert-manager 部署](https://cert-manager.io/docs/images/cert-manager-overview.svg)

# cert-manager 生产环境部署指南

本文档介绍如何在 Kubernetes 生产环境中部署 `cert-manager`，并确保其安全、稳定、高可用运行。

## 一、准备工作

- 已部署 Kubernetes 集群（版本建议 >= 1.19）
- 集群中已安装 `kubectl`、具备管理员权限
- 推荐为 cert-manager 创建独立的命名空间
- 网络能访问目标证书颁发机构（如 Let's Encrypt）

## 二、部署步骤

### 1. 创建命名空间

```bash
kubectl create namespace cert-manager
```

### 2. 添加 Jetstack Helm 仓库

```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update
```

### 3. 安装 cert-manager（使用 CRDs）

```bash
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.14.3 \
  --set installCRDs=true
```

> ✅ 说明：请根据实际需要调整版本号。

### 4. 验证部署

```bash
kubectl get pods -n cert-manager
```

应看到如下组件正在运行：

- cert-manager
- cert-manager-webhook
- cert-manager-cainjector

### 5. 创建 ClusterIssuer 示例（以 Let's Encrypt 为例）

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod-account-key
    solvers:
      - http01:
          ingress:
            class: nginx
```

```bash
kubectl apply -f clusterissuer.yaml
```

### 6. 创建 Certificate 示例（自动签发证书）

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-tls-cert
  namespace: default
spec:
  secretName: my-tls-cert-secret
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  commonName: example.com
  dnsNames:
    - example.com
    - www.example.com
```

### 7. 配合 Ingress 使用

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
    - hosts:
        - example.com
      secretName: my-tls-cert-secret
  rules:
    - host: example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: example-service
                port:
                  number: 80
```

## 三、安全与高可用建议

- 为 cert-manager 配置资源限制和 pod disruption budget
- 将证书 Secret 配置备份策略（如 Velero）
- 监控证书状态，结合 Prometheus/Grafana 设置告警
- 使用 DNS-01 验证适应多 ingress controller 和私有网络环境
- 建议结合 ExternalDNS 和自动化 DNS 提供商（如 Cloudflare）

## 四、常见问题排查

| 问题                             | 解决方案                                                             |
|----------------------------------|----------------------------------------------------------------------|
| Challenge 一直 pending           | 检查 DNS 或 Ingress 配置，确保验证路径可达                          |
| secret 未创建                    | 确保 Issuer 正确配置，邮箱正确，ACME 服务可访问                     |
| 证书未续签                       | 检查 cert-manager 日志，确认 webhook、cainjector 正常运行           |

## 五、参考资料

- [cert-manager 官方文档](https://cert-manager.io/docs/)
- [Helm 安装参考](https://cert-manager.io/docs/installation/helm/)
- [Let's Encrypt 使用指南](https://letsencrypt.org/)
