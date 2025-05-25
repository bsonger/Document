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

# 🏗️ 在生产环境安装 Cilium 并集成 Envoy

## 一、生产环境安装 Cilium

### 1️⃣ 环境准备
- Kubernetes v1.23+（建议 v1.25+）
- Linux 内核版本 >= 5.10（支持高级 eBPF 特性）
- 已安装 Helm、kubectl、containerd 或其他 CRI

### 2️⃣ 添加 Cilium Helm 仓库
```bash
helm repo add cilium https://helm.cilium.io/
helm repo update
```

### 3️⃣ 安装 Cilium（推荐开启 kube-proxy 替代）
```bash
helm install cilium cilium/cilium \
  --version 1.14.4 \
  --namespace kube-system \
  --set kubeProxyReplacement=strict \
  --set k8sServiceHost=<API_SERVER_HOST> \
  --set k8sServicePort=6443 \
  --set enableHubble=true \
  --set hubble.relay.enabled=true \
  --set hubble.ui.enabled=true
```

> 注意：替换 `<API_SERVER_HOST>` 为你的 Kubernetes API Server 地址。

### 4️⃣ 安装 CLI 工具（可选）
```bash
cilium status
cilium connectivity test
```

---

## 二、集成 Envoy 实现 L7 策略或流量代理

### 📌 Cilium 中的 Envoy 用法
Cilium 内部集成了 Envoy 作为透明代理，配合 eBPF socket hooks 实现 L7 网络策略与应用协议过滤（如 HTTP、gRPC、Kafka）。你无需手动部署 Envoy Sidecar。

### ✅ 启用 L7 HTTP 策略功能
Cilium 默认启用了内置 Envoy，用于处理 L7 策略。你只需部署 CiliumNetworkPolicy 即可实现应用层控制。

### 示例策略：仅允许 GET 请求 `/api`
```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: l7-http-policy
spec:
  endpointSelector:
    matchLabels:
      app: frontend
  ingress:
    - fromEndpoints:
        - matchLabels:
            app: backend
      toPorts:
        - ports:
            - port: "80"
              protocol: TCP
          rules:
            http:
              - method: GET
                path: "/api"
```

> Cilium 会自动在需要时为目标 Pod 启动 Envoy 实例，作为 L7 数据流的过滤器。

---

## 三、验证与调试

### 查看 Envoy 状态
```bash
cilium envoy config <pod-name>
cilium monitor --type l7
```

### 查看可视化界面（Hubble）
```bash
kubectl port-forward -n kube-system svc/hubble-ui 12000:80
# 打开浏览器 http://localhost:12000
```

---

## 四、总结
- 在生产环境中部署 Cilium 时推荐启用 kube-proxy 替代（strict 模式）以提升性能
- Cilium 原生集成 Envoy，无需单独部署 Sidecar 即可实现 L7 策略控制
- 搭配 Hubble 可实现端到端网络可视化、HTTP/gRPC 请求监控

Cilium 提供了轻量、高性能的 Service Mesh 替代方案，适合对安全、观测、性能有高要求的生产集群环境。
