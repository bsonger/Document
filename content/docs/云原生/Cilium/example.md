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

# 🧠 Cilium 生产环境复杂示例（L7 策略 + Egress Gateway + ClusterMesh）

本示例模拟一个典型的生产微服务架构场景，包含以下特性：

- 使用 Cilium 作为 CNI
- 基于标签的 L3/L4 策略控制
- 精细的 L7 HTTP 策略（基于路径、方法）
- Egress Gateway 限定 Pod 出口 IP
- 跨命名空间策略
- 多集群通信（ClusterMesh 简要配置）

---

## 📦 应用拓扑说明

```
frontend (namespace: frontend)
   |
   |--> backend (namespace: backend)
               |
               |--> external API (egress: api.external.com)
```

---

## 📑 1. 跨命名空间的 L7 策略

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: allow-frontend-to-backend
  namespace: backend
spec:
  endpointSelector:
    matchLabels:
      app: backend
  ingress:
    - fromEndpoints:
        - matchLabels:
            app: frontend
        - namespaceSelector:
            matchLabels:
              name: frontend
      toPorts:
        - ports:
            - port: "8080"
              protocol: TCP
          rules:
            http:
              - method: "GET"
                path: "/api/v1/status"
              - method: "POST"
                path: "/api/v1/submit"
```

---

## 🌐 2. Egress Gateway 配置（出公网使用固定 IP）

### Egress Gateway Policy

```yaml
apiVersion: cilium.io/v2alpha1
kind: CiliumEgressGatewayPolicy
metadata:
  name: backend-egress
spec:
  egress:
    - podSelector:
        matchLabels:
          app: backend
      namespaceSelector:
        matchLabels:
          name: backend
      destinationCIDRs:
        - "0.0.0.0/0"
      egressGateway: 
        nodeSelector:
          matchLabels:
            role: egress-gw
```

> 对应节点需设置标识：`kubectl label node <node> role=egress-gw`

---

## 🔐 3. DNS-based 策略（只允许访问特定域名）

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: dns-egress-policy
  namespace: backend
spec:
  endpointSelector:
    matchLabels:
      app: backend
  egress:
    - toFQDNs:
        - matchName: "api.external.com"
```

---

## ☁️ 4. ClusterMesh 基础配置步骤（简略）

```bash
# 在两个集群中都安装 Cilium 且开启 ClusterMesh
helm upgrade cilium cilium/cilium   --namespace kube-system   --set cluster.name=cluster-a   --set cluster.id=1   --set clustermesh.useAPIServer=true   --set etcd.enabled=true
```

- 使用共享 DNS 或 CoreDNS + etcd 连接两个集群
- 服务和身份会在两个集群间同步

---

## ✅ 5. 验证策略

```bash
cilium status
cilium monitor --type l7
cilium connectivity test
kubectl exec <frontend-pod> -- curl -X GET backend.backend.svc.cluster.local:8080/api/v1/status
```

---

## 🧾 总结

此示例演示了 Cilium 在实际生产中如何组合多个功能：

- 精确的应用层访问控制（HTTP Path/Method）
- 命名空间隔离
- 外部服务访问控制
- Egress IP 管理
- 多集群服务互通（ClusterMesh）

Cilium 可替代多种传统组件，提供统一的、高性能的网络 + 安全 + 可观测平台。
