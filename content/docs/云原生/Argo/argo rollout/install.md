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
# 🚀 Argo Rollouts 生产环境部署指南

本指南介绍如何在生产环境中部署 Argo Rollouts，集成监控系统，实现金丝雀发布、蓝绿部署、自动分析等高级发布策略。

---

## ✅ 一、组件功能概览

| 组件 | 说明 |
|------|------|
| Rollouts Controller | 替代原生 Deployment 控制器，实现高级发布策略 |
| AnalysisRun | 进行 Prometheus 等指标分析 |
| Experiment | 支持 AB 测试对比部署效果 |
| Argo Rollouts Kubectl Plugin | 用于 CLI 管理 Rollout 状态 |
| Dashboard UI（可选） | 与 Argo CD 集成的 Rollouts 插件界面 |

---

## 🧱 二、部署 Rollouts Controller

### 方法一：快速安装（官方 YAML）

```bash
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
```

### 方法二：使用 Helm（推荐生产）

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

helm install argo-rollouts argo/argo-rollouts   --namespace argo-rollouts --create-namespace   --set controller.replicaCount=2   --set metrics.enabled=true   --set dashboard.enabled=true
```

---

## 🔧 三、配置 Prometheus 集成（指标分析支持）

Rollouts 默认支持将 AnalysisRun 数据输出到 Prometheus，示例配置：

```yaml
analysis:
  templates:
    - name: success-rate-check
      metrics:
        - name: error-rate
          provider:
            prometheus:
              address: http://prometheus.monitoring.svc.cluster.local:9090
              query: rate(http_requests_total{status=~"5.."}[1m])
```

Prometheus 地址需要与你集群中实际部署的服务匹配。

---

## 🎨 四、与 Argo CD 集成 UI

如果你使用 Argo CD，并启用了 Rollouts Dashboard 插件，可以在应用中看到 Rollout 状态图：

```bash
--set dashboard.enabled=true  # 在 Helm 安装时启用
```

UI 可展示：当前阶段、金丝雀比例、活动服务/预览服务等信息。

---

## 🛡️ 五、生产部署建议

| 分类 | 建议 |
|------|------|
| 高可用性 | controller 设置为多副本，使用 readiness/liveness probe |
| 安全性 | 限制 namespace 访问权限，审计 AnalysisRun 输出 |
| 日志分析 | 配合 Loki/EFK 收集 controller 日志 |
| 指标分析 | 强烈建议配合 Prometheus/Datadog |
| 流量控制 | 建议使用 Istio/Nginx 服务网关进行金丝雀流量控制 |
| 自动化 | 可结合 GitOps（Argo CD）管理 Rollout YAML |
| Rollback 支持 | 启用 `automatic rollback` 条件判断策略 |

---

## 🔄 六、金丝雀/蓝绿发布策略示例（基础）

### 金丝雀发布策略

```yaml
strategy:
  canary:
    steps:
    - setWeight: 10
    - pause: {duration: 2m}
    - setWeight: 50
    - pause: {}
    analysis:
      templates:
        - templateName: success-rate-check
```

### 蓝绿部署策略

```yaml
strategy:
  blueGreen:
    activeService: app-service-active
    previewService: app-service-preview
    autoPromotionEnabled: false
```

---

## 📚 推荐阅读

- 官方文档：https://argo-rollouts.readthedocs.io/
- 指标分析配置：https://argo-rollouts.readthedocs.io/en/stable/features/analysis/
- Prometheus 监控接入：https://argo-rollouts.readthedocs.io/en/stable/features/metrics/
