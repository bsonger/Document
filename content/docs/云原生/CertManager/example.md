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

# cert-manager 复杂生产环境使用示例

![复杂使用示意图](https://cert-manager.io/docs/images/cert-manager-overview.svg)

# cert-manager 复杂生产环境使用示例

本示例展示了如何在生产环境中使用 cert-manager 实现以下目标：

- 使用 Let's Encrypt 作为证书颁发机构
- 使用 DNS-01 验证（支持多 ingress controller 和私有环境）
- 将证书应用于 Ingress（支持多域名）
- 自动续期与 Secret 管理
- 多租户支持（ClusterIssuer + 命名空间 Certificate）

---

## 一、前提条件

- Kubernetes 集群已部署 cert-manager
- DNS 提供商支持 API 控制（如 Cloudflare、AliDNS、Route53）
- 集群可访问 DNS API（配置 Secret 以供认证）

---

## 二、配置 Cloudflare DNS API Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloudflare-api-token-secret
  namespace: cert-manager
type: Opaque
stringData:
  api-token: CLOUDFLARE_API_TOKEN
```

> 请将 `CLOUDFLARE_API_TOKEN` 替换为实际 token。建议使用 read/write 权限控制子域名。

---

## 三、创建 ClusterIssuer 使用 DNS-01 验证

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-dns
spec:
  acme:
    email: your-email@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-dns-account-key
    solvers:
    - dns01:
        cloudflare:
          email: your-email@example.com
          apiTokenSecretRef:
            name: cloudflare-api-token-secret
            key: api-token
```

---

## 四、创建生产 Certificate 对象（多域名 + 自定义 DNS）

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: multi-domain-cert
  namespace: production
spec:
  secretName: tls-multi-domain
  duration: 2160h  # 90天
  renewBefore: 360h  # 15天前续期
  issuerRef:
    name: letsencrypt-dns
    kind: ClusterIssuer
  commonName: app.example.com
  dnsNames:
    - app.example.com
    - api.example.com
    - dashboard.example.com
```

---

## 五、配置支持 HTTPS 的 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  namespace: production
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-dns
spec:
  tls:
    - hosts:
        - app.example.com
        - api.example.com
        - dashboard.example.com
      secretName: tls-multi-domain
  rules:
    - host: app.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app-service
                port:
                  number: 80
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 80
    - host: dashboard.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: dashboard-service
                port:
                  number: 80
```

---

## 六、建议与注意事项

- 推荐为每个租户或业务线使用独立的 Certificate + Secret
- 可使用 ExternalDNS 实现自动化 DNS 管理
- 配合 Prometheus 监控证书剩余有效期（cert-manager 提供 metrics）
- 保持 cert-manager Controller、webhook 和 cainjector 高可用部署

---

## 七、故障排查

```bash
kubectl describe certificate multi-domain-cert -n production
kubectl logs -l app=cert-manager -n cert-manager
```

---

## 八、参考文档

- [cert-manager DNS-01 官方文档](https://cert-manager.io/docs/configuration/acme/dns01/)
- [Cloudflare API Token 权限配置](https://developers.cloudflare.com/api/tokens/create/)
