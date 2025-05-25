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

# 🚀 Argo Events 生产环境部署指南

本指南介绍如何在 Kubernetes 上以生产级配置部署 Argo Events，结合 NATS 作为 EventBus、Ingress 曝露 Webhook、TLS 安全、以及 Prometheus 监控等增强功能。

---

## ✅ 一、安装前准备

- Kubernetes 1.21+ 集群
- 已安装 Ingress Controller（如 NGINX）
- 可选组件：Prometheus、Cert-Manager、Argo Workflows

---

## 🧱 二、部署 Argo Events 控制器 + CRD

```bash
kubectl create ns argo-events

kubectl apply -n argo-events -f https://github.com/argoproj/argo-events/releases/latest/download/install.yaml
```

这会安装：
- CRDs: EventSource, Sensor, EventBus
- 控制器 Deployment

---

## 📡 三、部署 NATS EventBus（默认推荐）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: EventBus
metadata:
  name: default
  namespace: argo-events
spec:
  nats:
    native:
      replicas: 3
      auth: token
    containerTemplate:
      resources:
        requests:
          cpu: 100m
          memory: 128Mi
```

创建命令：
```bash
kubectl apply -f eventbus-nats.yaml
```

---

## 🌐 四、暴露 Webhook Source（使用 Ingress + TLS）

### EventSource 示例：Webhook

```yaml
apiVersion: argoproj.io/v1alpha1
kind: EventSource
metadata:
  name: webhook-source
  namespace: argo-events
spec:
  service:
    ports:
      - port: 12000
        targetPort: 12000
  webhook:
    push:
      endpoint: /webhook
      method: POST
      port: 12000
```

### 配置 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: events-webhook
  namespace: argo-events
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
    - hosts: [events.example.com]
      secretName: webhook-cert
  rules:
    - host: events.example.com
      http:
        paths:
          - path: /webhook
            pathType: Prefix
            backend:
              service:
                name: webhook-source-eventsource-svc
                port:
                  number: 12000
```

---

## 🔐 五、安全与认证建议

| 类别 | 建议 |
|------|------|
| 接口认证 | 使用网关或 Ingress 层的 Basic Auth 或 JWT 验证 |
| TLS 加密 | 使用 Cert-Manager 自动签发证书 |
| Secret 加密 | 使用 Sealed Secrets 或 External Secret 管理器 |
| RBAC 最小权限 | 限制 Sensor 控制器的角色为最小可运行权限 |

---

## 📈 六、监控与可观测性

### 启用 metrics

所有 Argo Events 控制器默认支持 Prometheus metrics，可通过如下配置暴露端口：

```yaml
spec:
  containers:
    - name: sensor-controller
      ports:
        - name: metrics
          containerPort: 8080
```

### Prometheus 配置抓取

```yaml
scrape_configs:
  - job_name: 'argo-events'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        regex: argo-events
        action: keep
```

---

## 🔄 七、自动触发 Argo Workflows 示例（Sensor）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: webhook-sensor
  namespace: argo-events
spec:
  dependencies:
    - name: push
      eventSourceName: webhook-source
      eventName: push
  triggers:
    - template:
        name: start-workflow
        k8s:
          operation: create
          source:
            resource:
              apiVersion: argoproj.io/v1alpha1
              kind: Workflow
              metadata:
                generateName: triggered-
              spec:
                entrypoint: hello
                templates:
                  - name: hello
                    container:
                      image: alpine
                      command: [echo]
                      args: ["Hello Argo Events"]
```

---

## ✅ 总结部署要点

| 组件 | 是否必须 | 说明 |
|------|----------|------|
| EventSource/Sensor Controller | ✅ 是 | 事件处理核心 |
| EventBus（NATS） | ✅ 是 | 默认消息通道 |
| Ingress + Cert-Manager | ✅ 推荐 | Webhook 安全接入 |
| Prometheus | ⛔ 可选 | 可观测性增强 |
| Argo Workflows | ⛔ 可选 | 如果使用 Events 触发 Workflow |

---

## 📚 推荐资源

- Argo Events 官网：https://argoproj.github.io/argo-events/
- EventBus 类型：https://argoproj.github.io/argo-events/eventbus/
- 示例合集：https://github.com/argoproj/argo-events/tree/master/examples
